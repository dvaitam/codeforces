package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./1532C.go"

type testCase struct {
	input string
	data  []query
}

type query struct {
	n int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := evaluate(tc, refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1532C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func evaluate(tc testCase, refOut, candOut string) error {
	refStrings := splitLines(refOut, len(tc.data))
	if len(refStrings) != len(tc.data) {
		return fmt.Errorf("reference produced %d lines, expected %d", len(refStrings), len(tc.data))
	}
	candStrings := splitLines(candOut, len(tc.data))
	if len(candStrings) != len(tc.data) {
		return fmt.Errorf("candidate produced %d lines, expected %d", len(candStrings), len(tc.data))
	}

	for i, q := range tc.data {
		target := minFreq(refStrings[i], q.k)
		if target == -1 {
			return fmt.Errorf("reference output invalid on query %d; got %q", i+1, refStrings[i])
		}
		candStr := candStrings[i]
		if err := validateString(candStr, q, target); err != nil {
			return fmt.Errorf("query %d invalid: %v", i+1, err)
		}
	}
	return nil
}

func splitLines(out string, limit int) []string {
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" && len(filtered) == limit {
			continue
		}
		filtered = append(filtered, strings.TrimSpace(line))
	}
	return filtered
}

func validateString(s string, q query, target int) error {
	if len(s) != q.n {
		return fmt.Errorf("expected length %d, got %d", q.n, len(s))
	}
	counts := make([]int, q.k)
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < 'a' || ch >= byte('a'+q.k) {
			return fmt.Errorf("invalid character %q", ch)
		}
		idx := int(ch - 'a')
		counts[idx]++
	}
	for i := 0; i < q.k; i++ {
		if counts[i] == 0 {
			return fmt.Errorf("missing letter %c", 'a'+i)
		}
	}
	for i := 0; i < q.k; i++ {
		if counts[i] < target {
			return fmt.Errorf("letter %c appears %d times, less than required %d", 'a'+i, counts[i], target)
		}
	}
	return nil
}

func minFreq(s string, k int) int {
	if len(s) == 0 {
		return -1
	}
	counts := make([]int, k)
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < 'a' || ch >= byte('a'+k) {
			return -1
		}
		counts[int(ch-'a')]++
	}
	min := len(s)
	for i := 0; i < k; i++ {
		if counts[i] == 0 {
			return -1
		}
		if counts[i] < min {
			min = counts[i]
		}
	}
	return min
}

func buildTests() []testCase {
	return []testCase{
		makeTest("3\n7 3\n4 4\n6 2\n", []query{{7, 3}, {4, 4}, {6, 2}}),
		makeTest("2\n1 1\n5 1\n", []query{{1, 1}, {5, 1}}),
		makeTest("4\n5 3\n6 3\n8 4\n10 5\n", []query{{5, 3}, {6, 3}, {8, 4}, {10, 5}}),
		makeTest("2\n26 26\n26 5\n", []query{{26, 26}, {26, 5}}),
		makeTest("1\n100 4\n", []query{{100, 4}}),
	}
}

func makeTest(input string, qs []query) testCase {
	return testCase{input: input, data: qs}
}
