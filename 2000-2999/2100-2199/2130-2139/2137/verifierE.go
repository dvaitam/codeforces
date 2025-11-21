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
	refSource   = "2000-2999/2100-2199/2130-2139/2137/2137E.go"
	targetTests = 140
	maxTotalN   = 200000
	maxValue    = 200000
	mod         = 998244353
)

type testCase struct {
	n int
	k int64
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

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

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2137E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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
		res[i] = val % mod
		if res[i] < 0 {
			res[i] += mod
		}
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

	// Samples from statement.
	add(testCase{n: 3, k: 3, a: []int64{0, 2, 1}})
	add(testCase{n: 2, k: 4, a: []int64{0, 2}})
	add(testCase{n: 4, k: 10, a: []int64{0, 0, 1, 1}})
	add(testCase{n: 8, k: 7, a: []int64{6, 6, 2, 4, 3, 0, 1, 8}})
	add(testCase{n: 2, k: 20, a: []int64{0, 0}})

	// Random tests.
	for len(tests) < targetTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain < 2 {
			break
		}
		n := rng.Intn(min(2000, remain-1)) + 2
		k := rng.Int63n(1_000_000_000) + 1
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(min(maxValue, n+5)))
		}
		add(testCase{n: n, k: k, a: a})
	}

	if len(tests) == 0 {
		add(testCase{n: 2, k: 1, a: []int64{0, 1}})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
