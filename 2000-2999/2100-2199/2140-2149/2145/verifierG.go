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
	refSource   = "2000-2999/2100-2199/2140-2149/2145/2145G.go"
	targetTests = 120
	maxN        = 2000
	maxM        = 2000
	maxK        = 4000
	mod         = 998244353
)

type testCase struct {
	n int
	m int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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

	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Fprintf(os.Stderr, "test %d: mismatched answer length: expected %d values, got %d\n", idx+1, len(refAns), len(candAns))
			os.Exit(1)
		}
		for i := range refAns {
			if refAns[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "test %d: value %d mismatch: expected %d, got %d\n", idx+1, i+1, refAns[i], candAns[i])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2145G-ref-*")
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

func parseOutput(out string, tc testCase) ([]int64, error) {
	tokens := strings.Fields(out)
	expected := tc.n + tc.m - 1 - min(tc.n, tc.m) + 1
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	ans := make([]int64, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse integer %q at position %d: %v", tokens[i], i+1, err)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		ans[i] = val
	}
	return ans, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0

	add := func(tc testCase) {
		if tc.n < 2 || tc.m < 2 {
			return
		}
		if totalN+tc.n > maxTotalN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Sample inputs from statement.
	add(testCase{n: 2, m: 3, k: 2})
	add(testCase{n: 2, m: 3, k: 3})
	add(testCase{n: 3, m: 2, k: 4})
	add(testCase{n: 2, m: 2, k: 3})
	add(testCase{n: 2, m: 2, k: 2})
	add(testCase{n: 2, m: 2, k: 1})
	add(testCase{n: 5, m: 3, k: 4})
	add(testCase{n: 3, m: 5, k: 3})
	add(testCase{n: 5, m: 2, k: 2})
	add(testCase{n: 2, m: 5, k: 2})

	// Random cases.
	for len(tests) < targetTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		if remain < 2 {
			break
		}
		n := rng.Intn(min(maxN, remain-1)) + 2
		m := rng.Intn(maxM-1) + 2
		if m > maxM {
			m = maxM
		}
		k := rng.Intn(n+m-1) + 1
		add(testCase{n: n, m: m, k: k})
	}

	if len(tests) == 0 {
		add(testCase{n: 2, m: 2, k: 1})
	}
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
