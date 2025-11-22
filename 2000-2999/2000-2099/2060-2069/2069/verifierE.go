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
	refSource     = "2069E.go"
	totalLenLimit = 480000
)

type testCase struct {
	s      string
	a, b   int
	ab, ba int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != got[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, expect[i], got[i])
			fmt.Fprintf(os.Stderr, "s(len=%d)=%s\n", len(tc.s), describeString(tc.s))
			fmt.Fprintf(os.Stderr, "a=%d b=%d ab=%d ba=%d\n", tc.a, tc.b, tc.ab, tc.ba)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2069E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Join(dir, refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
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

func parseOutputs(out string, expected int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output starting at token %q", tokens[expected])
	}
	ans := make([]string, expected)
	for i := 0; i < expected; i++ {
		tok := strings.ToUpper(tokens[i])
		if tok != "YES" && tok != "NO" {
			return nil, fmt.Errorf("token %q is not YES/NO", tokens[i])
		}
		ans[i] = tok
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(tc.a))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.b))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.ab))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.ba))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func describeString(s string) string {
	if len(s) <= 50 {
		return s
	}
	return s[:50] + "... (truncated)"
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	totalLen := 0
	add := func(tc testCase) {
		if totalLen+len(tc.s) > totalLenLimit {
			return
		}
		tests = append(tests, tc)
		totalLen += len(tc.s)
	}

	for _, tc := range deterministicTests() {
		add(tc)
	}

	// Large structured cases.
	add(testCase{s: strings.Repeat("A", 80000), a: 80000, b: 0, ab: 0, ba: 0})
	add(testCase{s: strings.Repeat("A", 80000), a: 60000, b: 0, ab: 0, ba: 0})
	add(testCase{s: strings.Repeat("B", 60000), a: 0, b: 60000, ab: 0, ba: 0})
	add(testCase{s: strings.Repeat("AB", 50000), a: 0, b: 0, ab: 50000, ba: 0})
	add(testCase{s: strings.Repeat("BA", 35000), a: 0, b: 0, ab: 0, ba: 35000})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalLen < totalLenLimit-1000 && len(tests) < 400 {
		maxLen := 2000
		if remaining := totalLenLimit - totalLen; remaining < maxLen {
			maxLen = remaining
		}
		tc := randomCase(rng, maxLen)
		add(tc)
	}

	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		// Samples from the statement.
		{s: "A", a: 0, b: 0, ab: 10, ba: 10},
		{s: "B", a: 0, b: 1, ab: 0, ba: 0},
		{s: "ABA", a: 0, b: 0, ab: 1, ba: 1},
		{s: "ABBABAAB", a: 5, b: 5, ab: 0, ba: 0},
		{s: "ABABBAABBAAB", a: 1, b: 1, ab: 2, ba: 3},
		{s: "ABBBBAB", a: 0, b: 3, ab: 2, ba: 0},
		{s: "BAABBA", a: 1, b: 3, ab: 2, ba: 0},

		// Tiny edges.
		{s: "A", a: 1, b: 0, ab: 0, ba: 0},
		{s: "B", a: 0, b: 1, ab: 0, ba: 0},
		{s: "AB", a: 1, b: 1, ab: 0, ba: 0},
		{s: "BA", a: 0, b: 0, ab: 0, ba: 1},
		{s: "AA", a: 2, b: 0, ab: 0, ba: 0},
		{s: "BB", a: 0, b: 2, ab: 0, ba: 0},
		{s: "AABBAABB", a: 4, b: 4, ab: 0, ba: 0},
		{s: "ABAB", a: 0, b: 0, ab: 2, ba: 0},
		{s: "BABA", a: 0, b: 0, ab: 0, ba: 2},
	}
}

func randomCase(rng *rand.Rand, maxLen int) testCase {
	n := 1 + rng.Intn(maxLen)
	s := randomString(rng, n)
	cntA, cntB := countLetters(s)
	maxPairs := n / 2

	var aCap, bCap, abCap, baCap int
	switch rng.Intn(4) {
	case 0: // Generous caps: should be feasible.
		aCap, bCap = cntA, cntB
		abCap, baCap = maxPairs, maxPairs
	case 1: // Tight singles, more pairs.
		aCap = rng.Intn(min(cntA, 3) + 1)
		bCap = rng.Intn(min(cntB, 3) + 1)
		abCap = rng.Intn(maxPairs+5) + maxPairs/2
		baCap = rng.Intn(maxPairs+5) + maxPairs/2
	case 2: // Restrict one pair type heavily.
		aCap = rng.Intn(cntA + 1)
		bCap = rng.Intn(cntB + 1)
		if rng.Intn(2) == 0 {
			abCap = rng.Intn(maxPairs/2 + 1)
			baCap = rng.Intn(maxPairs + 2)
		} else {
			baCap = rng.Intn(maxPairs/2 + 1)
			abCap = rng.Intn(maxPairs + 2)
		}
	default: // Random caps near counts.
		aCap = max(0, cntA-rng.Intn(cntA+3))
		bCap = max(0, cntB-rng.Intn(cntB+3))
		abCap = rng.Intn(maxPairs + 3)
		baCap = rng.Intn(maxPairs + 3)
	}

	return testCase{s: s, a: aCap, b: bCap, ab: abCap, ba: baCap}
}

func randomString(rng *rand.Rand, n int) string {
	var sb strings.Builder
	sb.Grow(n)
	cur := byte('A')
	if rng.Intn(2) == 1 {
		cur = 'B'
	}
	for sb.Len() < n {
		run := 1 + rng.Intn(5)
		for i := 0; i < run && sb.Len() < n; i++ {
			sb.WriteByte(cur)
		}
		if cur == 'A' {
			cur = 'B'
		} else {
			cur = 'A'
		}
	}
	return sb.String()
}

func countLetters(s string) (int, int) {
	a, b := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == 'A' {
			a++
		} else {
			b++
		}
	}
	return a, b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
