package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	name  string
	input string
}

var (
	topicWords = map[int][]string{
		1: {"culture", "art", "museum", "poet", "music", "novel", "literature", "theatre", "theater", "movie", "film", "festival", "painting", "artist", "opera", "gallery", "sculpture", "drama", "ballet", "author", "story", "heritage", "classic"},
		2: {"government", "minister", "president", "parliament", "policy", "election", "politics", "state", "party", "law", "constitution", "senate", "congress", "cabinet", "kremlin", "duma", "administration", "security", "diplomat", "military", "war", "conflict", "reform", "prime", "governor", "opposition", "referendum", "authority", "campaign", "official", "regulation", "senator"},
		3: {"economy", "market", "trade", "company", "business", "industry", "investment", "finance", "bank", "profit", "price", "export", "import", "currency", "capital", "product", "factory", "income", "revenue", "earning", "budget", "inflation", "tax", "loan", "credit", "retail", "sales", "stock", "share", "dollar", "euro", "ruble", "yen", "gdp", "oil", "gas", "metal", "energy", "analyst", "contract", "investor", "portfolio", "dividend", "logistics", "freight"},
	}
	fillerWords = []string{"smart", "beaver", "analysis", "global", "urgent", "reported", "today", "yesterday", "during", "while", "because", "after", "before", "around", "beyond", "rapidly", "unexpected", "insightful", "context", "review", "statement", "noted", "commented", "briefing", "scenario"}
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD9.go /path/to/solution")
		os.Exit(1)
	}
	candidate := os.Args[1]
	refSrc := referencePath()
	refBin, cleanup, err := buildReferenceBinary(refSrc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	total := len(tests)
	for i, tc := range tests {
		expectStr, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("reference failed on case %d (%s): %v\n", i+1, tc.name, err)
			os.Exit(1)
		}
		want, err := normalizeLabel(expectStr)
		if err != nil {
			fmt.Printf("reference produced invalid output on case %d (%s): %v\noutput:\n%s\n", i+1, tc.name, err, expectStr)
			os.Exit(1)
		}

		gotStr, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Printf("case %d (%s): runtime error: %v\ninput:\n%s\n", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := normalizeLabel(gotStr)
		if err != nil {
			fmt.Printf("case %d (%s): invalid output: %v\nfull output:\n%s\n", i+1, tc.name, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d (%s) failed: expected %d, got %d\ninput:\n%s\n", i+1, tc.name, want, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", total)
}

func referencePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "0-999/200-299/200-209/207/207D9.go"
	}
	return filepath.Join(filepath.Dir(file), "207D9.go")
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifierD9")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(tmpDir, "ref207D9")
	cmd := exec.Command("go", "build", "-o", binPath, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference solution: %v\n%s", err, out.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, out.String())
	}
	return out.String(), nil
}

func normalizeLabel(out string) (int, error) {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return 0, fmt.Errorf("empty output")
	}
	fields := strings.Fields(trimmed)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %q", trimmed)
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("not an integer: %v", err)
	}
	if val < 1 || val > 3 {
		return 0, fmt.Errorf("label %d out of range", val)
	}
	return val, nil
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(0x5a7b9c3d))
	tests := []testCase{}

	tests = append(tests, testCase{
		name: "art exhibition report",
		input: buildDocument(101, "Gallery Retrospective", []string{
			"Poet Laureate described the museum ballet program and a new opera stage.",
			"This festival review praises classic literature and the national theatre.",
			"Critics from the heritage council wrote stories about the sculptural drama."}),
	})

	tests = append(tests, testCase{
		name: "economic bulletin with digits",
		input: buildDocument(202, "Quarterly Market Sheet", []string{
			"The company finance board cited market trade figures worth $4500 and rising profit margins.",
			"Export contracts reference oil logistics, energy analysts, and doubled revenue.",
			"Investors noted the currency spread at 72.5 percent while the stock desk tracked loan risk."}),
	})

	tests = append(tests, testCase{
		name: "cabinet briefing",
		input: buildDocument(303, "Cabinet Evening Update", []string{
			"The president met the parliament speaker as ministers discussed security reform.",
			"A special senate panel reviewed opposition policy and law revisions with diplomatic envoys.",
			"Military governors warned about the border conflict and campaign schedule."}),
	})

	tests = append(tests, testCase{
		name: "digit heavy default",
		input: buildDocument(404, "Spreadsheet Dump", []string{
			"9911 invoices 2210 7700 5500 3300 1100", "42 24 84 48 12 60 72 36 18 6"}),
	})

	tests = append(tests, testCase{
		name:  "short artistic blurb",
		input: buildDocument(505, "Poem", []string{"music ballet", "opera"}),
	})

	tests = append(tests, testCase{
		name:  "political slogan",
		input: buildDocument(606, "Speech", []string{"law campaign", "minister"}),
	})

	tests = append(tests, testCase{
		name:  "trade note",
		input: buildDocument(707, "Note", []string{"trade profit", "market"}),
	})

	for topic := 1; topic <= 3; topic++ {
		for i := 0; i < 20; i++ {
			tests = append(tests, randomDoc(rng, topic, i))
		}
	}
	return tests
}

func buildDocument(id int, title string, paragraphs []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", id))
	sb.WriteString(title)
	if !strings.HasSuffix(title, "\n") {
		sb.WriteByte('\n')
	}
	for _, p := range paragraphs {
		sb.WriteString(p)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomDoc(rng *rand.Rand, topic int, idx int) testCase {
	lineCount := rng.Intn(6) + 3
	paragraphs := make([]string, 0, lineCount+1)
	for i := 0; i < lineCount; i++ {
		if rng.Intn(9) == 0 {
			paragraphs = append(paragraphs, "")
			continue
		}
		paragraphs = append(paragraphs, randomSentence(topic, rng))
	}
	if topic == 3 && rng.Intn(3) == 0 {
		paragraphs = append(paragraphs, fmt.Sprintf("budget tables show %d entries and $%d bonuses", rng.Intn(900)+100, rng.Intn(9000)+1000))
	}
	if topic == 2 && rng.Intn(3) == 0 {
		paragraphs = append(paragraphs, fmt.Sprintf("lawmakers counted %d ballots during the referendum", rng.Intn(4000)+600))
	}
	if topic == 1 && rng.Intn(3) == 0 {
		paragraphs = append(paragraphs, fmt.Sprintf("curators added %d new sculptures to the travelling art show", rng.Intn(50)+5))
	}
	title := fmt.Sprintf("Synthetic subject %d report %d", topic, idx)
	return testCase{
		name:  fmt.Sprintf("synthetic subject %d #%d", topic, idx),
		input: buildDocument(rng.Intn(1_000_000), title, paragraphs),
	}
}

func randomSentence(topic int, rng *rand.Rand) string {
	length := rng.Intn(10) + 8
	words := make([]string, length)
	for i := 0; i < length; i++ {
		word := fillerWords[rng.Intn(len(fillerWords))]
		if rng.Intn(4) != 0 {
			word = mutateWord(topicWords[topic][rng.Intn(len(topicWords[topic]))], rng)
		}
		if rng.Intn(7) == 0 {
			word = strings.ToUpper(word)
		} else if rng.Intn(6) == 0 {
			word = capitalize(word)
		}
		if topic == 3 && rng.Intn(9) == 0 {
			word = fmt.Sprintf("$%d", rng.Intn(9000)+100)
		}
		words[i] = word
	}
	sentence := strings.Join(words, " ")
	punctuation := []string{".", "!", "?", "..."}
	sentence += punctuation[rng.Intn(len(punctuation))]
	if topic == 3 && rng.Intn(3) == 0 {
		sentence += fmt.Sprintf(" (%d percent)", rng.Intn(70)+10)
	}
	if topic == 2 && rng.Intn(4) == 0 {
		sentence += fmt.Sprintf(" %d%% turnout", rng.Intn(60)+40)
	}
	if topic == 1 && rng.Intn(4) == 0 {
		sentence += " featuring poets"
	}
	return sentence
}

func mutateWord(base string, rng *rand.Rand) string {
	suffixes := []string{"", "s", "ing", "ed", "al", "ism", "ist", "ian", "ary"}
	prefixes := []string{"", "neo", "post", "pre", "anti", "micro", "ultra"}
	word := base
	if rng.Intn(3) == 0 {
		word = prefixes[rng.Intn(len(prefixes))] + word
	}
	word += suffixes[rng.Intn(len(suffixes))]
	if rng.Intn(5) == 0 {
		word = strings.ReplaceAll(word, "oo", "u")
	}
	if rng.Intn(7) == 0 {
		word = strings.ReplaceAll(word, "ph", "f")
	}
	return word
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
