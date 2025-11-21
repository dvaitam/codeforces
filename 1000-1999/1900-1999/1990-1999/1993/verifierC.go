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

type testCase struct {
	n int
	k int
	a []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1993C-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, "1993C.go")
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
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, k: 1, a: []int{5}},
		{n: 2, k: 1, a: []int{1, 100}},
		{n: 4, k: 4, a: []int{2, 3, 4, 5}},
		{n: 4, k: 3, a: []int{2, 3, 4, 5}},
		{n: 5, k: 3, a: []int{1, 10, 100, 1000, 10000}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	totalN := 0
	for len(tests) < 60 && totalN < 200000 {
		n := rng.Intn(5000) + 1
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		k := rng.Intn(n) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1_000_000_000) + 1
		}
		tests = append(tests, testCase{n: n, k: k, a: arr})
		totalN += n
	}
	return tests
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(strings.TrimSpace(expected))
	act := strings.Fields(strings.TrimSpace(actual))
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := range exp {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at test %d: expected %s got %s", i+1, exp[i], act[i])
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

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := compareOutputs(expected, actual, len(tests)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
