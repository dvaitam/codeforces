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
	n   int
	arr []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2038D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "2038D.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, arr: []int{0}},
		{n: 2, arr: []int{1, 2}},
		{n: 3, arr: []int{3, 4, 6}},
		{n: 3, arr: []int{1, 1, 1}},
		{n: 4, arr: []int{0, 0, 0, 0}},
		{n: 5, arr: []int{1000, 1000, 1000, 1000, 1000}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 30)
	for i := 0; i < 20; i++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(16)
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	for i := 0; i < 10; i++ {
		n := rng.Intn(5000) + 1000
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(1_000_000_000)
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	return tests
}

func stressTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	n := 200000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1_000_000_000)
	}
	return []testCase{
		{n: n, arr: arr},
	}
}

func compareOutputs(expected, actual string) error {
	exp := strings.TrimSpace(expected)
	act := strings.TrimSpace(actual)
	if exp != act {
		return fmt.Errorf("expected %s got %s", exp, act)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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

	for idx, tc := range tests {
		input := buildInput(tc)
		expected, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		actual, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(expected, actual); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
