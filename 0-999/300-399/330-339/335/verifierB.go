package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "0-999/300-399/330-339/335/335B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refPal := parseOutput(refOut)

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		ans := parseOutput(candOut)

		if err := validatePalindrome(ans); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): candidate output invalid palindrome: %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		if !isSubsequence(ans, strings.TrimSpace(tc.input)) {
			fmt.Fprintf(os.Stderr, "test %d (%s): output is not subsequence of input\ninput:\n%soutput:\n%s", idx+1, tc.name, tc.input, candOut)
			os.Exit(1)
		}
		if len(ans) < len(refPal) {
			fmt.Fprintf(os.Stderr, "test %d (%s): palindrome too short (expected >=%d, got %d)\ninput:\n%soutput:\n%s", idx+1, tc.name, len(refPal), len(ans), tc.input, candOut)
			os.Exit(1)
		}
		if len(ans) == 100 && len(refPal) < 100 {
			fmt.Fprintf(os.Stderr, "test %d (%s): palindrome length 100 not possible, but candidate returned 100\ninput:\n%soutput:\n%s", idx+1, tc.name, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier335B-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string) string {
	return strings.TrimSpace(out)
}

func validatePalindrome(s string) error {
	for len(s) > 1 {
		if s[0] != s[len(s)-1] {
			return fmt.Errorf("not a palindrome")
		}
		s = s[1 : len(s)-1]
	}
	return nil
}

func isSubsequence(sub, input string) bool {
	s := []rune(strings.TrimSpace(sub))
	t := []rune(strings.TrimSpace(input))
	i := 0
	for _, ch := range t {
		if i < len(s) && s[i] == ch {
			i++
		}
	}
	return i == len(s)
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample1", input: "bbbabcbbb\n"},
		{name: "sample2", input: "rquwmzexectvnbanemsmdufrg\n"},
		{name: "single", input: "a\n"},
		longTest(),
	}
	tests = append(tests, randomTests(40)...)
	return tests
}

func longTest() testCase {
	var sb strings.Builder
	n := 50000
	for i := 0; i < n; i++ {
		sb.WriteByte('a' + byte(i%26))
	}
	sb.WriteByte('\n')
	return testCase{name: "long", input: sb.String()}
}

func randomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(200) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		tests = append(tests, testCase{name: fmt.Sprintf("random_%d", i+1), input: sb.String() + "\n"})
	}
	return tests
}
