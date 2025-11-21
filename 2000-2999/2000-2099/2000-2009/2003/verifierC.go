package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2000-2009/2003/2003C.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2003C.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	name string
	s    string
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "sample1", s: "abc"},
		{name: "sample2", s: "edddf"},
		{name: "sample3", s: "turtle"},
		{name: "sample4", s: "pppppppp"},
		{name: "sample5", s: "codeforces"},
		{name: "all_same", s: "aaaaa"},
		{name: "two_letters", s: "abababab"},
		{name: "increasing", s: "abcdefghijklmnopqrstuvwxyz"},
		{name: "palindrome", s: "racecar"},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalLen := 0
	for len(tests) < 80 && totalLen < 200000 {
		n := rng.Intn(200) + 2
		if totalLen+n > 200000 {
			n = 200000 - totalLen
		}
		if n <= 0 {
			break
		}
		builder := strings.Builder{}
		for i := 0; i < n; i++ {
			ch := rune('a' + rng.Intn(26))
			builder.WriteRune(ch)
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", len(tests)+1),
			s:    builder.String(),
		})
		totalLen += n
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n%s\n", len(tc.s), tc.s))
	}
	return sb.String()
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
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2003C-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2003C")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, count int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) < count {
		return nil, fmt.Errorf("expected %d strings, got %d tokens", count, len(tokens))
	}
	if len(tokens) > count {
		return nil, fmt.Errorf("expected %d strings, but got %d (extra tokens)", count, len(tokens))
	}
	return tokens, nil
}

func freqSignature(s string) [26]int {
	var freq [26]int
	for _, ch := range s {
		if ch < 'a' || ch > 'z' {
			continue
		}
		freq[ch-'a']++
	}
	return freq
}

func goodPairs(s string) int64 {
	n := len(s)
	total := int64(n) * int64(n-1) / 2
	if n <= 1 {
		return 0
	}
	var bad int64
	prevLen := -1
	curLen := 1
	for i := 1; i < n; i++ {
		if s[i] == s[i-1] {
			curLen++
		} else {
			if prevLen != -1 {
				bad += int64(prevLen) * int64(curLen)
			}
			prevLen = curLen
			curLen = 1
		}
	}
	if prevLen != -1 {
		bad += int64(prevLen) * int64(curLen)
	}
	return total - bad
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refStrings, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	targetGood := make([]int64, len(tests))
	origFreq := make([][26]int, len(tests))
	for i, tc := range tests {
		origFreq[i] = freqSignature(tc.s)
		refStr := refStrings[i]
		if len(refStr) != len(tc.s) {
			fmt.Fprintf(os.Stderr, "reference output for test %s (%d) has wrong length: got %d expected %d\n", tc.name, i+1, len(refStr), len(tc.s))
			os.Exit(1)
		}
		if freqSignature(refStr) != origFreq[i] {
			fmt.Fprintf(os.Stderr, "reference output for test %s (%d) is not a permutation of input\n", tc.name, i+1)
			os.Exit(1)
		}
		targetGood[i] = goodPairs(refStr)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userStrings, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		ans := userStrings[i]
		if len(ans) != len(tc.s) {
			fmt.Fprintf(os.Stderr, "test %s (%d): expected length %d, got %d\n", tc.name, i+1, len(tc.s), len(ans))
			os.Exit(1)
		}
		if freqSignature(ans) != origFreq[i] {
			fmt.Fprintf(os.Stderr, "test %s (%d): string is not a permutation of input\n", tc.name, i+1)
			os.Exit(1)
		}
		if good := goodPairs(ans); good != targetGood[i] {
			fmt.Fprintf(os.Stderr, "test %s (%d): expected %d good pairs, got %d\n", tc.name, i+1, targetGood[i], good)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
