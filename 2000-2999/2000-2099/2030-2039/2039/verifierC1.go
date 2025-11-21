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
	x int
	m int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2039C1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC1")
	cmd := exec.Command("go", "build", "-o", outPath, "2039C1.go")
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.x, tc.m))
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{x: 1, m: 1},
		{x: 1, m: 10},
		{x: 2, m: 2},
		{x: 6, m: 9},
		{x: 5, m: 7},
		{x: 2, m: 36},
		{x: 4, m: 1},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	sumX := 0
	for sumX < 8_000_000 && len(tests) < 50 {
		x := rng.Intn(1_000_000) + 1
		m := rng.Int63n(1_000_000_000) + 1
		tests = append(tests, testCase{x: x, m: m})
		sumX += x
	}
	for len(tests) < 60 {
		x := rng.Intn(1_000_000) + 1
		m := rng.Int63n(1_000_000_000_000) + 1
		tests = append(tests, testCase{x: x, m: m})
	}
	return tests
}

func stressTests() []testCase {
	tests := make([]testCase, 0, 3)
	tests = append(tests,
		testCase{x: 1_000_000, m: 1_000_000_000_000_000_000},
		testCase{x: 999_983, m: 999_983},
		testCase{x: 999_983, m: 1_000_000_000_000_000_000},
	)
	return tests
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at test %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
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
	tests = append(tests, stressTests()...)
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
