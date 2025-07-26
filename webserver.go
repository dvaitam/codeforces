package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type contestInfo struct {
	ID       string
	Path     string
	Problems []string
}

var contests map[string]*contestInfo
var db *sql.DB

var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Contests</h1>
<ul>
{{range .}}
<li><a href="/contest/{{.ID}}">{{.ID}}</a></li>
{{end}}
</ul>
</body></html>`))

var contestTmpl = template.Must(template.New("contest").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Contest {{.ID}}</h1>
<ul>
{{range .Problems}}
<li><a href="/contest/{{$.ID}}/problem/{{.}}">Problem {{.}}</a></li>
{{end}}
</ul>
<h2>Add Problem</h2>
<form action="/addproblem" method="post">
<input type="hidden" name="contest" value="{{.ID}}">
Letter: <input name="letter"><br>
Admin Key: <input type="password" name="adminkey"><br>
<textarea name="statement" rows="10" cols="80"></textarea><br>
<input type="submit" value="Add Problem">
</form>
</body></html>`))

var addProblemTmpl = template.Must(template.New("addproblem").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Add Problem</h1>
<form action="/addproblem" method="post">
Contest ID: <input name="contest" value="{{.Contest}}"><br>
Letter: <input name="letter" value="{{.Letter}}"><br>
Admin Key: <input type="password" name="adminkey"><br>
<textarea name="statement" rows="10" cols="80">{{.Statement}}</textarea><br>
<input type="submit" value="Add Problem">
</form>
</body></html>`))

var problemTmpl = template.Must(template.New("problem").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
body { max-width: 800px; margin: auto; font-family: sans-serif; }
pre { white-space: pre-wrap; word-wrap: break-word; }
textarea { width: 100%; }
</style>
</head>
<body>
<h1>Contest {{.Contest}} Problem {{.Letter}}</h1>
<pre>{{.Statement}}</pre>
<form action="/contest/{{.Contest}}/problem/{{.Letter}}/submit" method="post" enctype="multipart/form-data">
<select name="lang">
<option value="c">C</option>
<option value="cpp">C++</option>
<option value="java">Java</option>
<option value="python">Python 3</option>
<option value="go">Go</option>
<option value="rust">Rust</option>
</select><br>
<textarea name="code" rows="20" cols="80"></textarea><br>
<input type="file" name="file"><br>
<input type="submit" value="Submit">
</form>
</body></html>`))

var resultTmpl = template.Must(template.New("result").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
body { max-width: 800px; margin: auto; font-family: sans-serif; }
pre { white-space: pre-wrap; word-wrap: break-word; }
</style>
</head>
<body>
<h1>Result for Contest {{.Contest}} Problem {{.Letter}}</h1>
<pre>{{.Output}}</pre>
<a href="/contest/{{.Contest}}/problem/{{.Letter}}">Back</a>
</body></html>`))

var textTmpl = template.Must(template.New("text").Parse(`<!DOCTYPE html><html><body><pre>{{.}}</pre></body></html>`))

var leaderboardTmpl = template.Must(template.New("leaderboard").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Leaderboard</h1>
<table border="1">
<tr><th>Run ID</th><th>Model</th><th>Rating</th><th>Timestamp</th></tr>
{{range .Leaders}}
<tr><td><a href="/leaderboard?run={{.RunID}}">{{.RunID}}</a></td><td>{{.Model}}</td><td>{{.Rating}}</td><td>{{.Timestamp}}</td></tr>
{{end}}
</table>
{{if .Evals}}
<h2>Evaluation History for {{.RunID}}</h2>
<table border="1">
<tr><th>Run ID</th><th>Model</th><th>Problem ID</th><th>Rating</th><th>Success</th><th>Timestamp</th><th>Prompt</th><th>Response</th></tr>
{{range .Evals}}
<tr><td>{{.RunID}}</td><td>{{.Model}}</td><td><a href="/contest/{{.ContestID}}/problem/{{.IndexName}}">{{.ProblemID}}</a></td><td>{{.Rating}}</td><td>{{.Success}}</td><td>{{.Timestamp}}</td><td><a href="/evaluation/prompt/{{.ID}}">View</a></td><td><a href="/evaluation/response/{{.ID}}">View</a></td></tr>
{{end}}
</table>
{{end}}
</body></html>`))

var modelsTmpl = template.Must(template.New("models").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Models</h1>
<ul>
{{range .}}
<li><a href="/model?name={{.}}">{{.}}</a></li>
{{end}}
</ul>
</body></html>`))

var modelTmpl = template.Must(template.New("model").Parse(`
<!DOCTYPE html>
<html><body>
<h1>Evaluations for {{.Model}}</h1>
<table border="1">
<tr><th>Run ID</th><th>Problem ID</th><th>Rating</th><th>Success</th><th>Timestamp</th><th>Prompt</th><th>Response</th></tr>
{{range .Evals}}
<tr><td>{{.RunID}}</td><td><a href="/contest/{{.ContestID}}/problem/{{.IndexName}}">{{.ProblemID}}</a></td><td>{{.Rating}}</td><td>{{.Success}}</td><td>{{.Timestamp}}</td><td><a href="/evaluation/prompt/{{.ID}}">View</a></td><td><a href="/evaluation/response/{{.ID}}">View</a></td></tr>
{{end}}
</table>
</body></html>`))

func scanContests(root string) (map[string]*contestInfo, error) {
	result := make(map[string]*contestInfo)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if _, err := strconv.Atoi(base); err == nil {
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil
			}
			var probs []string
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, "problem") && strings.HasSuffix(name, ".txt") {
					letter := strings.TrimSuffix(strings.TrimPrefix(name, "problem"), ".txt")
					probs = append(probs, letter)
				}
			}
			if len(probs) > 0 {
				sort.Strings(probs)
				result[base] = &contestInfo{ID: base, Path: path, Problems: probs}
			}
		}
		return nil
	})
	return result, err
}

func findVerifier(dir, letter string) string {
	cand := filepath.Join(dir, "verifier"+letter+".go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	cand = filepath.Join(dir, "verifier.go")
	if _, err := os.Stat(cand); err == nil {
		return cand
	}
	return ""
}

func detectJavaClassName(src []byte) string {
	re := regexp.MustCompile(`(?m)^\s*public\s+(?:class|interface|enum|record)\s+([A-Za-z_][A-Za-z0-9_]*)`)
	if m := re.FindSubmatch(src); m != nil {
		return string(m[1])
	}
	return "Main"
}

func compileSource(srcPath, lang string) (string, string, error) {
	tmpDir, err := os.MkdirTemp("", "submit")
	if err != nil {
		return "", "", err
	}
	exe := filepath.Join(tmpDir, "main")
	var cmd *exec.Cmd
	switch lang {
	case "c":
		cmd = exec.Command("gcc", srcPath, "-O2", "-std=c11", "-o", exe)
	case "cpp":
		cmd = exec.Command("g++", srcPath, "-O2", "-std=c++17", "-o", exe)
	case "go":
		cmd = exec.Command("go", "build", "-o", exe, srcPath)
	case "rust":
		cmd = exec.Command("rustc", "-O", srcPath, "-o", exe)
	case "java":
		javaDir := filepath.Join(tmpDir, "java")
		if err := os.Mkdir(javaDir, 0755); err != nil {
			return "", "", err
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", "", err
		}
		className := detectJavaClassName(data)
		javaSrc := srcPath
		if filepath.Base(srcPath) != className+".java" {
			javaSrc = filepath.Join(tmpDir, className+".java")
			if err := os.WriteFile(javaSrc, data, 0644); err != nil {
				return "", "", err
			}
		}
		cmd = exec.Command("javac", "-d", javaDir, javaSrc)
		exe = filepath.Join(tmpDir, "run-java.sh")
		script := fmt.Sprintf("#!/bin/sh\njava -cp %s %s \"$@\"\n", javaDir, className)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
	case "python":
		absSrc, err := filepath.Abs(srcPath)
		if err != nil {
			return "", "", err
		}
		exe = filepath.Join(tmpDir, "run-python.sh")
		script := fmt.Sprintf("#!/bin/sh\npython3 %s \"$@\"\n", absSrc)
		if err := os.WriteFile(exe, []byte(script), 0755); err != nil {
			return "", "", err
		}
		return exe, "", nil
	default:
		return "", "", fmt.Errorf("unknown language")
	}

	if cmd != nil {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", string(out), err
		}
		if lang == "java" {
			// ensure wrapper is executable
			return exe, string(out), nil
		}
	}
	return exe, "", nil
}

func submitSolution(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lang := r.FormValue("lang")
	var data []byte
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		data, _ = io.ReadAll(file)
	} else {
		data = []byte(r.FormValue("code"))
	}
	extMap := map[string]string{"c": ".c", "cpp": ".cpp", "java": ".java", "python": ".py", "go": ".go", "rust": ".rs"}
	ext := extMap[lang]
	if ext == "" {
		http.Error(w, "unknown language", http.StatusBadRequest)
		return
	}
	srcPath := filepath.Join(c.Path, "user"+strings.ToUpper(letter)+ext)
	if err := os.WriteFile(srcPath, data, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exe, compileOut, err := compileSource(srcPath, lang)
	output := bytes.Buffer{}
	if err != nil {
		output.WriteString("Compilation failed:\n")
		output.WriteString(compileOut)
		output.WriteString(err.Error())
	} else {
		verifier := findVerifier(c.Path, letter)
		if verifier != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "go", "run", filepath.Base(verifier), exe)
			cmd.Dir = c.Path
			res, err := cmd.CombinedOutput()
			output.Write(res)
			if ctx.Err() == context.DeadlineExceeded {
				output.WriteString("\nVerifier timed out after 10 seconds")
			} else if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					output.WriteString(fmt.Sprintf("\nVerifier exited with status %d", ee.ExitCode()))
				} else {
					output.WriteString("\nVerifier error: " + err.Error())
				}
			}
		} else {
			output.WriteString("Compiled successfully. No verifier available.")
		}
	}
	resultTmpl.Execute(w, map[string]string{
		"Contest": c.ID,
		"Letter":  letter,
		"Output":  output.String(),
	})
}

func problemPage(w http.ResponseWriter, r *http.Request, c *contestInfo, letter string) {
	stmtPath := filepath.Join(c.Path, "problem"+letter+".txt")
	data, err := os.ReadFile(stmtPath)
	if err != nil {
		http.Error(w, "problem not found", http.StatusNotFound)
		return
	}
	problemTmpl.Execute(w, map[string]string{
		"Contest":   c.ID,
		"Letter":    letter,
		"Statement": string(data),
	})
}

func contestPage(w http.ResponseWriter, r *http.Request, cid string) {
	c := contests[cid]
	if c == nil {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"+cid), "/")
	if len(parts) >= 3 && parts[1] == "problem" && parts[2] != "" {
		if len(parts) == 4 && parts[3] == "submit" {
			submitSolution(w, r, c, parts[2])
			return
		}
		if r.Method == http.MethodGet {
			problemPage(w, r, c, parts[2])
			return
		}
	}
	if r.URL.Path != "/contest/"+cid {
		http.NotFound(w, r)
		return
	}
	contestTmpl.Execute(w, c)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ids := make([]string, 0, len(contests))
	for id := range contests {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	var list []*contestInfo
	for _, id := range ids {
		list = append(list, contests[id])
	}
	indexTmpl.Execute(w, list)
}

func contestDir(id string) (string, error) {
	n, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	thousands := (n / 1000) * 1000
	tDir := fmt.Sprintf("%d-%d", thousands, thousands+999)
	n %= 1000
	hundreds := (n / 100) * 100
	hDir := fmt.Sprintf("%d-%d", thousands+hundreds, thousands+hundreds+99)
	n %= 100
	tens := (n / 10) * 10
	teDir := fmt.Sprintf("%d-%d", thousands+hundreds+tens, thousands+hundreds+tens+9)
	return filepath.Join(tDir, hDir, teDir, id), nil
}

func ensureContest(id string) (*contestInfo, error) {
	if c := contests[id]; c != nil {
		return c, nil
	}
	dir, err := contestDir(id)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	c := &contestInfo{ID: id, Path: dir}
	contests[id] = c
	return c, nil
}

func addProblemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		addProblemTmpl.Execute(w, map[string]string{
			"Contest":   r.URL.Query().Get("contest"),
			"Letter":    r.URL.Query().Get("letter"),
			"Statement": "",
		})
		return
	case http.MethodPost:
		contestID := r.FormValue("contest")
		letter := strings.ToUpper(r.FormValue("letter"))
		statement := r.FormValue("statement")
		adminkey := r.FormValue("adminkey")
		if adminkey != os.Getenv("ADMIN_KEY") {
			http.Error(w, "admin key mismatch", http.StatusForbidden)
			return
		}
		if contestID == "" || letter == "" {
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}
		c, err := ensureContest(contestID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stmtPath := filepath.Join(c.Path, "problem"+letter+".txt")
		if _, err := os.Stat(stmtPath); err == nil {
			http.Error(w, "problem already exists", http.StatusBadRequest)
			return
		}
		if err := os.WriteFile(stmtPath, []byte(statement), 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		found := false
		for _, l := range c.Problems {
			if l == letter {
				found = true
				break
			}
		}
		if !found {
			c.Problems = append(c.Problems, letter)
			sort.Strings(c.Problems)
		}
		http.Redirect(w, r, "/contest/"+contestID+"/problem/"+letter, http.StatusSeeOther)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/leaderboard" {
		http.NotFound(w, r)
		return
	}

	type Leader struct {
		RunID     string
		Model     string
		Rating    int
		Timestamp string
	}
	var leaders []Leader
	rows, err := db.Query("SELECT run_id, model, rating, timestamp FROM leaderboard ORDER BY rating DESC")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var l Leader
			if err = rows.Scan(&l.RunID, &l.Model, &l.Rating, &l.Timestamp); err == nil {
				leaders = append(leaders, l)
			}
		}
	}

	type Eval struct {
		ID        int
		RunID     string
		Model     string
		ProblemID int
		ContestID int
		IndexName string
		Rating    int
		Success   bool
		Timestamp string
	}
	var evals []Eval
	runIDFilter := r.URL.Query().Get("run")
	if runIDFilter != "" {
		rows, err = db.Query(`SELECT e.id, e.run_id, e.model, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), e.success, e.timestamp
                       FROM evaluations e
                       JOIN problems p ON e.problem_id = p.id
                       WHERE e.run_id = ? ORDER BY e.timestamp DESC`, runIDFilter)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var e Eval
				if err = rows.Scan(&e.ID, &e.RunID, &e.Model, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Success, &e.Timestamp); err == nil {
					evals = append(evals, e)
				}
			}
		}
	}

	leaderboardTmpl.Execute(w, map[string]interface{}{
		"Leaders": leaders,
		"Evals":   evals,
		"RunID":   runIDFilter,
	})
}

func modelHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/model" {
		http.NotFound(w, r)
		return
	}
	modelName := r.URL.Query().Get("name")
	if modelName == "" {
		var models []string
		rows, err := db.Query("SELECT DISTINCT model FROM evaluations ORDER BY model")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var m string
				if err = rows.Scan(&m); err == nil {
					models = append(models, m)
				}
			}
		}
		modelsTmpl.Execute(w, models)
		return
	}

	type Eval struct {
		ID        int
		RunID     string
		ProblemID int
		ContestID int
		IndexName string
		Rating    int
		Success   bool
		Timestamp string
	}
	var evals []Eval
	rows, err := db.Query(`SELECT e.id, e.run_id, e.problem_id, p.contest_id, p.index_name, COALESCE(p.rating, 0), e.success, e.timestamp
                               FROM evaluations e
                               JOIN problems p ON e.problem_id = p.id
                               WHERE e.model = ? ORDER BY e.timestamp DESC`, modelName)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Eval
			if err = rows.Scan(&e.ID, &e.RunID, &e.ProblemID, &e.ContestID, &e.IndexName, &e.Rating, &e.Success, &e.Timestamp); err == nil {
				evals = append(evals, e)
			}
		}
	}

	modelTmpl.Execute(w, map[string]interface{}{
		"Model": modelName,
		"Evals": evals,
	})
}

func evaluationContentHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/evaluation/"), "/")
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	field := parts[0]
	id, err := strconv.Atoi(parts[1])
	if err != nil || (field != "prompt" && field != "response") {
		http.NotFound(w, r)
		return
	}
	var content string
	err = db.QueryRow("SELECT "+field+" FROM evaluations WHERE id = ?", id).Scan(&content)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	textTmpl.Execute(w, content)
}

func main() {
	var err error
	contests, err = scanContests(".")
	if err != nil {
		panic(err)
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "user:pass@tcp(127.0.0.1:3306)/dbname"
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/addproblem", addProblemHandler)
	http.HandleFunc("/contest/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/contest/"), "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}
		contestPage(w, r, parts[0])
	})
	http.HandleFunc("/leaderboard", leaderboardHandler)
	http.HandleFunc("/model", modelHandler)
	http.HandleFunc("/evaluation/", evaluationContentHandler)
	http.ListenAndServe(":8081", nil)
}
