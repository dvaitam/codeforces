package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Problem struct {
	ID        int
	ContestID int
	IndexName string
	Statement string
}

var requestTimeout time.Duration

func main() {
	model := flag.String("model", "", "The AI model to use (e.g., anthropic/claude-3.5-sonnet)")
	provider := flag.String("provider", "openrouter", "Model provider: Gemini, OpenAI, xai, Claude, openrouter")
	dbDSN := flag.String("db", "user:pass@tcp(127.0.0.1:3306)/dbname", "Database DSN")
	maxAttempts := flag.Int("max-attempts", 1, "Maximum attempts to fix syntax errors (1-5)")
	httpTimeout := flag.Duration("timeout", 30*time.Second, "HTTP request timeout")
	flag.Parse()
	requestTimeout = *httpTimeout

	if *maxAttempts < 1 || *maxAttempts > 5 {
		fmt.Println("max-attempts must be between 1 and 5")
		os.Exit(1)
	}

	var apiKeyEnv string
	switch strings.ToLower(*provider) {
	case "openai":
		apiKeyEnv = "OPENAI_API_KEY"
	case "gemini":
		apiKeyEnv = "GEMINI_API_KEY"
	case "xai":
		apiKeyEnv = "XAI_API_KEY"
	case "claude":
		apiKeyEnv = "CLAUDE_API_KEY"
	default:
		apiKeyEnv = "OPENROUTER_API_KEY"
	}

	apiKey := os.Getenv(apiKeyEnv)
	if *model == "" || apiKey == "" {
		fmt.Printf("Usage: go run script.go -model=<model> -db=<dsn> [-max-attempts=1-5] -provider=%s\n", *provider)
		fmt.Printf("Set %s environment variable\n", apiKeyEnv)
		os.Exit(1)
	}

	db, err := sql.Open("mysql", *dbDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS evaluations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			run_id VARCHAR(255),
			model VARCHAR(255),
			problem_id INT,
			prompt TEXT,
			response TEXT,
			success BOOL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS leaderboard (
			id INT AUTO_INCREMENT PRIMARY KEY,
			run_id VARCHAR(255),
			model VARCHAR(255),
			rating INT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		panic(err)
	}

	runID := time.Now().Format("20060102-150405")

	availableRatings := getAvailableRatings(db)
	if len(availableRatings) == 0 {
		panic("No valid problems found in the database")
	}

	rand.Seed(time.Now().UnixNano())

	estimatedRating := 800
	for i := 0; i < 25; i++ {
		actualRating := clampToNearest(estimatedRating, availableRatings)
		fmt.Printf("Attempt %d: Targeting estimated %d (using actual rating %d)\n", i+1, estimatedRating, actualRating)
		problem, verifierFile := getRandomProblem(db, actualRating)
		plainStatement := latexToPlain(problem.Statement)
		prompt := "write a go solution for " + plainStatement + ". Output only the code, with no explanation or additional text."
		fmt.Printf("Sending prompt for Problem ID: %d, Contest ID: %d, Index: %s\n", problem.ID, problem.ContestID, problem.IndexName)
		fmt.Println("Sending prompt...")
		response := sendPrompt(*provider, *model, apiKey, prompt)
		fmt.Println("Response received.")

		code := extractCode(response)
		fmt.Printf("Solution code:\n%s\n", code)

		success := false
		finalResponse := response
		attempt := 1
		var tempBinAbs string
		for attempt <= *maxAttempts {
			fmt.Printf("Verification attempt %d of %d\n", attempt, *maxAttempts)
			buildSuccess, buildErrMsg, builtBinAbs := buildSolution(code)
			tempBinAbs = builtBinAbs
			if !buildSuccess {
				if attempt == *maxAttempts {
					break
				}
				fixPrompt := fmt.Sprintf("The following Go code has compilation errors: %s\n\nFix the errors and output only the corrected code, no explanation.", buildErrMsg)
				fixPrompt += "\n\nOriginal code:\n" + code
				fmt.Println("Sending fix prompt...")
				fixResponse := sendPrompt(*provider, *model, apiKey, fixPrompt)
				code = extractCode(fixResponse)
				finalResponse = fixResponse // Update final response to the corrected one
				fmt.Printf("Corrected code:\n%s\n", code)
				attempt++
				continue
			}

			// Build succeeded, now verify
			verifySuccess := runVerifier(verifierFile, tempBinAbs)
			if verifySuccess {
				success = true
			}
			break // No need for more attempts if build succeeds, verification result is final
		}

		// Clean up temp files if any
		if tempBinAbs != "" {
			os.Remove(tempBinAbs)
			os.Remove(filepath.Dir(tempBinAbs) + "/solution.go") // Clean source too
			os.RemoveAll(filepath.Dir(tempBinAbs))
		}

		_, err = db.Exec(
			"INSERT INTO evaluations (run_id, model, problem_id, prompt, response, success) VALUES (?, ?, ?, ?, ?, ?)",
			runID, *model, problem.ID, prompt, finalResponse, success,
		)
		if err != nil {
			panic(err)
		}

		if success {
			estimatedRating += 100
		} else {
			estimatedRating -= 100
		}
	}

	// Insert into leaderboard
	_, err = db.Exec(
		"INSERT INTO leaderboard (run_id, model, rating) VALUES (?, ?, ?)",
		runID, *model, estimatedRating,
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Evaluation complete. Estimated Codeforces rating for model %s: %d\n", *model, estimatedRating)
}

func buildSolution(code string) (bool, string, string) {
	tempDir, err := os.MkdirTemp("", "gobuild-*")
	if err != nil {
		return false, err.Error(), ""
	}

	tempSrc := filepath.Join(tempDir, "solution.go")
	err = os.WriteFile(tempSrc, []byte(code), 0644)
	if err != nil {
		return false, err.Error(), ""
	}

	// Format the code
	cmd := exec.Command("gofmt", "-w", "solution.go")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return false, err.Error(), ""
	}

	tempBin := filepath.Join(tempDir, "solution")
	cmd = exec.Command("go", "build", "-o", "solution", "solution.go")
	cmd.Dir = tempDir
	var buildStderr bytes.Buffer
	cmd.Stderr = &buildStderr
	if err := cmd.Run(); err != nil {
		return false, buildStderr.String(), tempDir // Return tempDir to clean later
	}

	tempBinAbs, _ := filepath.Abs(tempBin)
	return true, "", tempBinAbs
}

func runVerifier(verifierFile, tempBinAbs string) bool {
	verifierAbs, err := filepath.Abs(verifierFile)
	if err != nil {
		fmt.Printf("Error getting absolute path for verifier: %v\n", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", verifierAbs, tempBinAbs)
	cmd.Dir = filepath.Dir(verifierAbs)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Verification timed out after 10 seconds")
			return false
		}
		fmt.Printf("Verifier failed: %v\n", err)
		return false
	}

	fmt.Printf("Verifier stdout: %s\n", out.String())
	fmt.Printf("Verifier stderr: %s\n", stderr.String())

	return strings.Contains(out.String(), "tests passed")
}

func getAvailableRatings(db *sql.DB) []int {
	rows, err := db.Query("SELECT DISTINCT rating FROM problems WHERE rating IS NOT NULL AND statement IS NOT NULL ORDER BY rating")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var ratings []int
	for rows.Next() {
		var r int
		err := rows.Scan(&r)
		if err != nil {
			panic(err)
		}
		ratings = append(ratings, r)
	}

	var available []int
	for _, r := range ratings {
		if hasValidProblem(db, r) {
			available = append(available, r)
		}
	}

	sort.Ints(available) // Ensure sorted, though query orders it
	return available
}

func hasValidProblem(db *sql.DB, rating int) bool {
	rows, err := db.Query("SELECT contest_id, index_name FROM problems WHERE rating = ? AND statement IS NOT NULL", rating)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var contestID int
		var indexName string
		err := rows.Scan(&contestID, &indexName)
		if err != nil {
			panic(err)
		}
		dir := getContestDir(contestID)
		solFile := filepath.Join(dir, fmt.Sprintf("%d%s.go", contestID, indexName))
		verFile := filepath.Join(dir, fmt.Sprintf("verifier%s.go", indexName))

		if _, err := os.Stat(solFile); err == nil {
			if _, err := os.Stat(verFile); err == nil {
				return true
			}
		}
	}
	return false
}

func clampToNearest(target int, available []int) int {
	if len(available) == 0 {
		panic("No available ratings")
	}

	if target <= available[0] {
		return available[0]
	}
	if target >= available[len(available)-1] {
		return available[len(available)-1]
	}

	// Find the position where target would be inserted
	idx := sort.SearchInts(available, target)
	if available[idx] == target {
		return target
	}

	// Compare distance to available[idx-1] and available[idx]
	prev := available[idx-1]
	next := available[idx]
	distPrev := target - prev
	distNext := next - target

	if distPrev < distNext {
		return prev
	} else if distNext < distPrev {
		return next
	} else {
		// Equal distance, prefer the higher one
		return next
	}
}

func getRandomProblem(db *sql.DB, rating int) (Problem, string) {
	rows, err := db.Query("SELECT id, contest_id, index_name, statement FROM problems WHERE rating = ? AND statement IS NOT NULL", rating)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var candidates []Problem
	for rows.Next() {
		var p Problem
		err := rows.Scan(&p.ID, &p.ContestID, &p.IndexName, &p.Statement)
		if err != nil {
			panic(err)
		}
		candidates = append(candidates, p)
	}

	var valid []Problem
	var validVerifiers []string
	for _, p := range candidates {
		dir := getContestDir(p.ContestID)
		solFile := filepath.Join(dir, fmt.Sprintf("%d%s.go", p.ContestID, p.IndexName))
		verFile := filepath.Join(dir, fmt.Sprintf("verifier%s.go", p.IndexName))

		if _, err := os.Stat(solFile); err == nil {
			if _, err := os.Stat(verFile); err == nil {
				valid = append(valid, p)
				validVerifiers = append(validVerifiers, verFile)
			}
		}
	}

	if len(valid) == 0 {
		panic(fmt.Sprintf("No valid problems found for rating %d", rating))
	}

	idx := rand.Intn(len(valid))
	return valid[idx], validVerifiers[idx]
}

func getContestDir(contestID int) string {
	top := (contestID / 1000) * 1000
	topStr := fmt.Sprintf("%d-%d", top, top+999)

	second := (contestID / 100) * 100
	secondStr := fmt.Sprintf("%d-%d", second, second+99)

	third := (contestID / 10) * 10
	thirdStr := fmt.Sprintf("%d-%d", third, third+9)

	fourthStr := fmt.Sprintf("%d", contestID)

	return filepath.Join(topStr, secondStr, thirdStr, fourthStr)
}

func latexToPlain(text string) string {
	re := regexp.MustCompile(`\$\$\$(.*?)\$\$\$`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		sub := re.FindStringSubmatch(m)[1]

		// handle simple tables like \begin{array}{|c|c|} ... \end{array}
		if strings.Contains(sub, `\\begin{array}`) {
			arrRe := regexp.MustCompile(`(?s)\\begin{array}{[^}]*}(.*?)\\end{array}`)
			sub = arrRe.ReplaceAllStringFunc(sub, func(t string) string {
				inner := arrRe.FindStringSubmatch(t)[1]
				inner = strings.ReplaceAll(inner, `\\hline`, "")
				// line breaks and columns
				inner = strings.ReplaceAll(inner, `\\\\`, "\n")
				inner = strings.ReplaceAll(inner, `&`, " ")
				textRe := regexp.MustCompile(`\\text{([^{}]*)}`)
				inner = textRe.ReplaceAllString(inner, "$1")
				inner = strings.ReplaceAll(inner, `\\`, "")
				inner = strings.ReplaceAll(inner, "{", "")
				inner = strings.ReplaceAll(inner, "}", "")
				return inner
			})
			return sub
		}

		// remove sizing helpers early so other replacements don't break
		sub = strings.ReplaceAll(sub, `\left`, "")
		sub = strings.ReplaceAll(sub, `\right`, "")

		replacements := map[string]string{
			`\leq`:   "<=",
			`\le`:    "<=",
			`\geq`:   ">=",
			`\ge`:    ">=",
			`\cdot`:  "*",
			`\times`: "x",
			`\dots`:  "...",
		}
		for old, val := range replacements {
			sub = strings.ReplaceAll(sub, old, val)
		}

		// common LaTeX constructs
		fracRe := regexp.MustCompile(`\\frac{([^{}]+)}{([^{}]+)}`)
		sub = fracRe.ReplaceAllString(sub, "$1/$2")
		sub = strings.ReplaceAll(sub, `\lceil`, "ceil(")
		sub = strings.ReplaceAll(sub, `\rceil`, ")")

		textRe := regexp.MustCompile(`\\text{([^{}]*)}`)
		sub = textRe.ReplaceAllString(sub, "$1")

		sub = strings.ReplaceAll(sub, "\\", "")
		sub = strings.ReplaceAll(sub, "left", "")
		sub = strings.ReplaceAll(sub, "right", "")
		sub = strings.ReplaceAll(sub, "{", "")
		sub = strings.ReplaceAll(sub, "}", "")
		sub = strings.ReplaceAll(sub, " ", "")
		return sub
	})
}

func sendPrompt(provider, model, apiKey, prompt string) string {
	prompt = latexToPlain(prompt)
	fmt.Printf("Prompt length: %d characters\n", len(prompt))

	var body []byte
	var err error

	lowerProvider := strings.ToLower(provider)
	if lowerProvider == "gemini" {
		gemReq := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"parts": []map[string]string{{"text": prompt}},
				},
			},
		}
		body, err = json.Marshal(gemReq)
	} else {
		messages := []Message{{Role: "user", Content: prompt}}
		reqBody := Request{Model: model, Messages: messages}
		body, err = json.Marshal(reqBody)
	}
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return ""
	}

	client := &http.Client{Timeout: requestTimeout}
	url := ""
	headers := map[string]string{"Content-Type": "application/json"}

	switch strings.ToLower(provider) {
	case "openai":
		url = "https://api.openai.com/v1/chat/completions"
		headers["Authorization"] = "Bearer " + apiKey
	case "gemini":
		url = "https://generativelanguage.googleapis.com/v1beta/models/" + model + ":generateContent?key=" + apiKey
	case "xai":
		url = "https://api.x.ai/v1/chat/completions"
		headers["Authorization"] = "Bearer " + apiKey
	case "claude":
		url = "https://api.anthropic.com/v1/messages"
		headers["x-api-key"] = apiKey
		headers["anthropic-version"] = "2023-06-01"
	default:
		url = "https://openrouter.ai/api/v1/chat/completions"
		headers["Authorization"] = "Bearer " + apiKey
	}

	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			fmt.Printf("Error creating request (attempt %d): %v\n", attempt, err)
			if attempt == maxRetries {
				return ""
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}

		resp, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("Error sending request (attempt %d): %v\n", attempt, err)
			if attempt == maxRetries {
				return ""
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("API error (attempt %d): %s\n", attempt, string(bodyBytes))
			if attempt == maxRetries {
				return ""
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		if strings.ToLower(provider) == "gemini" {
			var gResp struct {
				Candidates []struct {
					Content struct {
						Parts []struct {
							Text string `json:"text"`
						} `json:"parts"`
					} `json:"content"`
				} `json:"candidates"`
			}
			if err = json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
				fmt.Printf("Error decoding response (attempt %d): %v\n", attempt, err)
				if attempt == maxRetries {
					return ""
				}
				time.Sleep(time.Second * time.Duration(attempt))
				continue
			}
			if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
				fmt.Printf("No response from API (attempt %d)\n", attempt)
				if attempt == maxRetries {
					return ""
				}
				time.Sleep(time.Second * time.Duration(attempt))
				continue
			}
			return gResp.Candidates[0].Content.Parts[0].Text
		}

		var apiResp Response
		if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding response (attempt %d): %v\n", attempt, err)
			if attempt == maxRetries {
				return ""
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		if len(apiResp.Choices) == 0 {
			fmt.Printf("No response from API (attempt %d)\n", attempt)
			if attempt == maxRetries {
				return ""
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		return apiResp.Choices[0].Message.Content
	}
	return ""
}

func extractCode(response string) string {
	re := regexp.MustCompile(`(?s)\x60\x60\x60go\s*(.*?)\x60\x60\x60`)
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return strings.TrimSpace(response)
}
