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

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		wantRes, err := parseOutput(want)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, want)
			os.Exit(1)
		}

		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if err := validateOutput(tc.input, got, wantRes.k, wantRes.len); err != nil {
			fmt.Fprintf(os.Stderr, "candidate invalid on test %d: %v\nInput:\n%s\nCandidate output:\n%s\n", idx+1, err, tc.input, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

type result struct {
	k   int
	len int
}

func parseOutput(out string) (*result, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	var k int
	if _, err := fmt.Sscanf(fields[0], "%d", &k); err != nil {
		return nil, fmt.Errorf("failed to parse k: %w", err)
	}
	pals := strings.Split(strings.TrimSpace(out[strings.Index(out, "\n")+1:]), " ")
	if len(pals) != k {
		return nil, fmt.Errorf("expected %d palindromes, got %d", k, len(pals))
	}
	if k == 0 {
		return nil, fmt.Errorf("k must be >= 1")
	}
	length := len(pals[0])
	for i := 1; i < len(pals); i++ {
		if len(pals[i]) != length {
			return nil, fmt.Errorf("palindrome lengths differ")
		}
	}
	return &result{k: k, len: length}, nil
}

func validateOutput(input, output string, bestK, bestLen int) error {
	n, s, err := parseInput(input)
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	var k int
	if _, err := fmt.Sscanf(fields[0], "%d", &k); err != nil {
		return fmt.Errorf("failed to parse k: %w", err)
	}
	if k <= 0 || n%k != 0 {
		return fmt.Errorf("invalid k %d", k)
	}
	palsLine := strings.TrimSpace(output[strings.Index(output, "\n")+1:])
	pals := strings.Fields(palsLine)
	if len(pals) != k {
		return fmt.Errorf("expected %d palindromes, got %d", k, len(pals))
	}
	length := len(pals[0])
	if length == 0 || n != k*length {
		return fmt.Errorf("invalid palindrome length %d", length)
	}
	for _, p := range pals {
		if len(p) != length {
			return fmt.Errorf("palindrome lengths differ")
		}
		if !isPalindrome(p) {
			return fmt.Errorf("string %s is not a palindrome", p)
		}
	}
	if k != bestK {
		return fmt.Errorf("non-optimal k: expected %d, got %d", bestK, k)
	}
	if length != bestLen {
		return fmt.Errorf("non-optimal palindrome length: expected %d, got %d", bestLen, length)
	}
	if !sameMultiset(strings.Join(pals, ""), s) {
		return fmt.Errorf("palindrome multiset does not match original string")
	}
	return nil
}

func parseInput(input string) (int, string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return 0, "", fmt.Errorf("not enough lines")
	}
	var n int
	fmt.Sscanf(lines[0], "%d", &n)
	return n, strings.TrimSpace(lines[1]), nil
}

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func sameMultiset(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	freq := make(map[rune]int)
	for _, ch := range a {
		freq[ch]++
	}
	for _, ch := range b {
		freq[ch]--
		if freq[ch] < 0 {
			return false
		}
	}
	for _, v := range freq {
		if v != 0 {
			return false
		}
	}
	return true
}

func locateReference() (string, error) {
	candidates := []string{
		"883H.go",
		filepath.Join("0-999", "800-899", "880-889", "883", "883H.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 883H.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref883H_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		newTest("1\na\n"),
		newTest("2\naa\n"),
		newTest("3\naaa\n"),
		newTest("6\naabbbb\n"),
		newTest("8\naabbccdd\n"),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTest(rng, rng.Intn(100)+1))
	}
	tests = append(tests, randomTest(rng, 400000))
	return tests
}

func randomTest(rng *rand.Rand, n int) testCase {
	if n < 1 {
		n = 1
	}
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		b.WriteByte(chars[rng.Intn(len(chars))])
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func newTest(data string) testCase {
	return testCase{input: data}
}
