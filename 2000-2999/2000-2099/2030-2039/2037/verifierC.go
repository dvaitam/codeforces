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

const maxSum = 400000

type testCase struct {
	n int
}

type result struct {
	ok   bool
	perm []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2037C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "2037C.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 16)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, tests []testCase) ([]result, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([]result, len(tests))
	for ti, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("not enough output for test %d", ti+1)
		}
		token := tokens[idx]
		idx++
		if token == "-1" {
			res[ti] = result{ok: false}
			continue
		}
		val, err := strconv.Atoi(token)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q in test %d", token, ti+1)
		}
		perm := make([]int, tc.n)
		perm[0] = val
		for i := 1; i < tc.n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("not enough numbers for test %d", ti+1)
			}
			v, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q in test %d", tokens[idx], ti+1)
			}
			idx++
			perm[i] = v
		}
		res[ti] = result{ok: true, perm: perm}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output values")
	}
	return res, nil
}

func sieve(limit int) []bool {
	comp := make([]bool, limit+1)
	comp[0], comp[1] = true, true
	for i := 2; i*i <= limit; i++ {
		if !comp[i] {
			for j := i * i; j <= limit; j += i {
				comp[j] = true
			}
		}
	}
	return comp
}

func validatePermutation(res result, n int, comp []bool) error {
	if !res.ok {
		return fmt.Errorf("expected permutation but got -1")
	}
	if len(res.perm) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(res.perm))
	}
	seen := make([]bool, n+1)
	for i, v := range res.perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d at position %d out of range", v, i+1)
		}
		if seen[v] {
			return fmt.Errorf("value %d appears multiple times", v)
		}
		seen[v] = true
	}
	for i := 0; i+1 < n; i++ {
		sum := res.perm[i] + res.perm[i+1]
		if sum >= len(comp) {
			return fmt.Errorf("sum %d exceeds sieve limit", sum)
		}
		if !comp[sum] {
			return fmt.Errorf("sum %d at positions %d,%d is prime", sum, i+1, i+2)
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 6},
		{n: 8},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	totalN := 0
	for len(tests) < cap(tests) && totalN < 190000 {
		n := rng.Intn(2000) + 2
		if totalN+n > 200000 {
			break
		}
		tests = append(tests, testCase{n: n})
		totalN += n
	}
	return tests
}

func compareResults(oracleRes, targetRes []result, tests []testCase, comp []bool) error {
	if len(oracleRes) != len(targetRes) {
		return fmt.Errorf("result count mismatch")
	}
	for i := range oracleRes {
		if !oracleRes[i].ok {
			if targetRes[i].ok {
				return fmt.Errorf("test %d: oracle says -1 but target produced permutation", i+1)
			}
			continue
		}
		if !targetRes[i].ok {
			return fmt.Errorf("test %d: valid permutation exists but target printed -1", i+1)
		}
		if err := validatePermutation(targetRes[i], tests[i].n, comp); err != nil {
			return fmt.Errorf("test %d invalid permutation: %v", i+1, err)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expectedOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	oracleRes, err := parseOutput(expectedOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expectedOut)
		os.Exit(1)
	}
	targetRes, err := parseOutput(actualOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actualOut)
		os.Exit(1)
	}

	comp := sieve(maxSum)
	if err := compareResults(oracleRes, targetRes, tests, comp); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}
