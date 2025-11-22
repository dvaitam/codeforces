package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource2045H = "2000-2999/2000-2099/2040-2049/2045/2045H.go"
	maxTotalLen    = 8000
)

type testCase struct {
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2045H)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := tc.s + "\n"

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		refCount, _, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, runErr := runProgram(candidate, input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, runErr, input)
			os.Exit(1)
		}
		candCount, words, parseErr := parseOutput(candOut)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, parseErr, input, candOut)
			os.Exit(1)
		}

		if err := validate(tc.s, refCount, candCount, words); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2045H-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string) (int, []string, error) {
	reader := strings.NewReader(out)
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return 0, nil, fmt.Errorf("missing word count: %v", err)
	}
	k, err := strconv.Atoi(first)
	if err != nil {
		return 0, nil, fmt.Errorf("first token not integer (%s)", first)
	}
	if k < 0 {
		return 0, nil, fmt.Errorf("negative word count %d", k)
	}
	words := make([]string, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &words[i]); err != nil {
			return 0, nil, fmt.Errorf("expected %d words, got %d (%v)", k, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return 0, nil, fmt.Errorf("extra output detected starting with %s", extra)
	}
	return k, words, nil
}

func validate(s string, bestCount, k int, words []string) error {
	if k != bestCount {
		return fmt.Errorf("word count mismatch: expected %d, got %d", bestCount, k)
	}
	if len(words) != k {
		return fmt.Errorf("word list length mismatch: expected %d, got %d", k, len(words))
	}
	var sb strings.Builder
	for i, w := range words {
		if len(w) == 0 {
			return fmt.Errorf("word %d is empty", i+1)
		}
		if i > 0 && !(words[i-1] < w) {
			return fmt.Errorf("words not strictly increasing at position %d (%s, %s)", i, words[i-1], w)
		}
		sb.WriteString(w)
	}
	if sb.String() != s {
		return fmt.Errorf("concatenation mismatch: got %s, expected %s", sb.String(), s)
	}
	return nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	total := totalLen(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for total < maxTotalLen {
		remain := maxTotalLen - total
		n := rng.Intn(min(400, remain)) + 1
		tests = append(tests, testCase{s: randomString(rng, n)})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "A"},
		{s: "ABACUS"},
		{s: "AAAAAA"},
		{s: "EDCBA"},
		{s: "ABABAB"},
		{s: "BAAAAB"},
		{s: "AZBYCXD"},
		{s: "ZZZZ"},
		{s: "ABCDEFGHIJKLMNOP"},
	}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func totalLen(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += len(tc.s)
	}
	return sum
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func resolveSourcePath(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if abs, err := filepath.Abs(cleaned); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	base := filepath.Base(path)
	if abs, err := filepath.Abs(base); err == nil {
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	return "", fmt.Errorf("reference source not found at %s", path)
}
