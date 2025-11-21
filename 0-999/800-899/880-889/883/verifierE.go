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
	"time"
)

const (
	refSource        = "883E.go"
	tempOraclePrefix = "oracle-883E-"
	randomTestsCount = 120
	maxN             = 50
	maxM             = 1000
)

type testCase struct {
	name    string
	n       int
	pattern string
	words   []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)
	tests = append(tests, largeTests()...)

	for idx, tc := range tests {
		input := formatInput(tc)

		expOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expected, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, gotOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if expected != got {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\n", idx+1, tc.name, expected, got)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Candidate output:")
			fmt.Print(gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseAnswer(out string) (int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(out)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", out)
	}
	return val, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	sb.WriteString(tc.pattern)
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(tc.words))
	for _, w := range tc.words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name:    "simple_partial",
			n:       4,
			pattern: "a*b*",
			words:   []string{"aabb", "acbd", "aebf"},
		},
		{
			name:    "no_letter_possible",
			n:       3,
			pattern: "**a",
			words:   []string{"bca", "cba", "dsa"},
		},
		{
			name:    "exact_known",
			n:       3,
			pattern: "aba",
			words:   []string{"aba"},
		},
		{
			name:    "alternating",
			n:       5,
			pattern: "*a*a*",
			words:   []string{"baaaa", "caaab", "daaaa"},
		},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxN-1) + 1
		m := rng.Intn(30) + 1
		if m > maxM {
			m = maxM
		}
		base := randomWord(n, rng)
		pattern := buildPattern(base, rng)
		words := buildWordList(pattern, base, m, rng)
		tests = append(tests, testCase{
			name:    fmt.Sprintf("random_%d", i+1),
			n:       n,
			pattern: pattern,
			words:   words,
		})
	}
	return tests
}

func largeTests() []testCase {
	n := maxN
	base := strings.Repeat("abcde", n/5+1)[:n]
	pattern := buildPattern(base, rand.New(rand.NewSource(1)))
	words := buildWordList(pattern, base, maxM, rand.New(rand.NewSource(2)))
	return []testCase{
		{name: "large_case", n: n, pattern: pattern, words: words},
	}
}

func randomWord(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func buildPattern(base string, rng *rand.Rand) string {
	n := len(base)
	b := []byte(base)
	hasStar := false
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '*'
			hasStar = true
		}
	}
	if !hasStar {
		pos := rng.Intn(n)
		b[pos] = '*'
	}
	return string(b)
}

func buildWordList(pattern, base string, m int, rng *rand.Rand) []string {
	words := make([]string, 0, m)
	seen := make(map[string]bool)
	addWord := func(w string) {
		if !seen[w] {
			seen[w] = true
			words = append(words, w)
		}
	}
	addWord(base)
	for len(words) < m {
		w := make([]byte, len(base))
		for i := 0; i < len(base); i++ {
			if pattern[i] == '*' {
				w[i] = byte('a' + rng.Intn(26))
			} else {
				w[i] = pattern[i]
			}
		}
		addWord(string(w))
	}
	return words
}
