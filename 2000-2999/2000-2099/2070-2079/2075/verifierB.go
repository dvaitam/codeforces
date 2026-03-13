package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	maxTotalN         = 5000
	targetTests       = 80
	maxSingleN        = 500
	maxValue    int64 = 1_000_000_000
)

type testCase struct {
	n int
	k int
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	input := buildInput(tests)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	// Basic validation: answers should be non-negative
	for i := range candAns {
		if candAns[i] < 0 {
			fmt.Fprintf(os.Stderr, "test %d: negative answer %d\n", i+1, candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d test cases).\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d: %v", tok, i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Manual small cases.
	add(testCase{n: 2, k: 1, a: []int64{1, 1}})
	add(testCase{n: 3, k: 1, a: []int64{1, 2, 3}})
	add(testCase{n: 4, k: 2, a: []int64{4, 2, 3, 1}})
	add(testCase{n: 5, k: 3, a: []int64{2, 2, 2, 2, 2}})

	for len(tests) < targetTests && totalN < maxTotalN {
		n := rng.Intn(minVal(100, maxTotalN-totalN)) + 2
		if n > maxSingleN {
			n = maxSingleN
		}
		k := rng.Intn(n-1) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Int63n(maxValue) + 1
		}
		add(testCase{n: n, k: k, a: arr})
	}

	if len(tests) == 0 {
		add(testCase{n: 2, k: 1, a: []int64{1, 1}})
	}
	return tests
}

func minVal(a, b int) int {
	if a < b {
		return a
	}
	return b
}
