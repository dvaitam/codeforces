package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

type inputData struct {
	n, l, k int
	letters string
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildReference() (string, error) {
	outPath := filepath.Join(verifierDir, "ref1267L.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1267L.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseInput(data string) (inputData, error) {
	reader := strings.NewReader(data)
	var n, l, k int
	var letters string
	if _, err := fmt.Fscan(reader, &n, &l, &k); err != nil {
		return inputData{}, fmt.Errorf("failed to read n,l,k: %v", err)
	}
	if _, err := fmt.Fscan(reader, &letters); err != nil {
		return inputData{}, fmt.Errorf("failed to read letters: %v", err)
	}
	return inputData{n: n, l: l, k: k, letters: letters}, nil
}

func parseWords(out string, n int) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	buf := make([]byte, 0, 1024)
	scanner.Buffer(buf, 1<<20)
	scanner.Split(bufio.ScanWords)
	var words []string
	for len(words) < n {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("failed to scan word: %v", err)
			}
			return nil, fmt.Errorf("expected %d words, got %d", n, len(words))
		}
		words = append(words, scanner.Text())
	}
	if scanner.Scan() {
		return nil, fmt.Errorf("too many words, unexpected token %q", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}

func countLetters(s string) [26]int {
	var cnt [26]int
	for _, ch := range s {
		if ch < 'a' || ch > 'z' {
			continue
		}
		cnt[ch-'a']++
	}
	return cnt
}

func validateWords(data inputData, words []string, expectedKth string) error {
	if len(words) != data.n {
		return fmt.Errorf("expected %d words, got %d", data.n, len(words))
	}
	prev := ""
	for i, w := range words {
		if len(w) != data.l {
			return fmt.Errorf("word %d has length %d, expected %d", i+1, len(w), data.l)
		}
		for _, ch := range w {
			if ch < 'a' || ch > 'z' {
				return fmt.Errorf("word %d contains invalid character %q", i+1, ch)
			}
		}
		if i > 0 && prev > w {
			return fmt.Errorf("words are not sorted at positions %d and %d", i, i+1)
		}
		prev = w
	}
	inputCnt := countLetters(data.letters)
	var outputCnt [26]int
	for _, w := range words {
		for _, ch := range w {
			outputCnt[ch-'a']++
		}
	}
	if inputCnt != outputCnt {
		return fmt.Errorf("output letters multiset does not match input")
	}
	kIdx := data.k - 1
	if kIdx < 0 || kIdx >= len(words) {
		return fmt.Errorf("invalid k index %d", data.k)
	}
	if words[kIdx] != expectedKth {
		return fmt.Errorf("k-th word mismatch: expected %s, got %s", expectedKth, words[kIdx])
	}
	return nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	data, err := parseInput(tc.input)
	if err != nil {
		return fmt.Errorf("invalid test input: %v", err)
	}
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut)
	}
	refWords, err := parseWords(refOut, data.n)
	if err != nil {
		return fmt.Errorf("failed to parse reference output: %v\noutput:\n%s", err, refOut)
	}
	if err := validateWords(data, refWords, refWords[data.k-1]); err != nil {
		return fmt.Errorf("reference produced invalid output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate runtime error: %v\n%s", err, candOut)
	}
	candWords, err := parseWords(candOut, data.n)
	if err != nil {
		return fmt.Errorf("failed to parse candidate output: %v\noutput:\n%s", err, candOut)
	}
	if err := validateWords(data, candWords, refWords[data.k-1]); err != nil {
		return fmt.Errorf("invalid output: %v\noutput:\n%s", err, candOut)
	}
	return nil
}

func formatInput(n, l, k int, letters string) string {
	return fmt.Sprintf("%d %d %d\n%s\n", n, l, k, letters)
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_letter", input: formatInput(1, 1, 1, "a")},
		{name: "two_words", input: formatInput(2, 2, 1, "bbaa")},
		{name: "duplicate_letters", input: formatInput(3, 1, 2, "abc")},
		{name: "all_same", input: formatInput(4, 3, 3, strings.Repeat("z", 12))},
	}
}

func randomLetters(n int, rng *rand.Rand) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	return sb.String()
}

func randomCase(name string, rng *rand.Rand, maxN, maxL int) testCase {
	n := rng.Intn(maxN) + 1
	l := rng.Intn(maxL) + 1
	k := rng.Intn(n) + 1
	letters := randomLetters(n*l, rng)
	return testCase{
		name:  name,
		input: formatInput(n, l, k, letters),
	}
}

func largeCase(name string, n, l, k int, seed int64) testCase {
	rng := rand.New(rand.NewSource(seed))
	letters := randomLetters(n*l, rng)
	return testCase{name: name, input: formatInput(n, l, k, letters)}
}

func generateTests() []testCase {
	tests := manualTests()
	deterministicSeeds := []int64{1, 2, 3, 4, 5, 6}
	for idx, seed := range deterministicSeeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomCase(fmt.Sprintf("deterministic_%d", idx+1), rng, 10, 10))
	}
	tests = append(tests,
		largeCase("larger_1", 30, 40, 15, 42),
		largeCase("larger_2", 50, 50, 25, 99),
		largeCase("larger_3", 80, 20, 60, 123),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		maxN := 20 + len(tests)%10
		maxL := 20 + (len(tests)*3)%30
		tests = append(tests, randomCase(fmt.Sprintf("random_%d", len(tests)+1), rng, maxN, maxL))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
