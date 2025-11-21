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
	tmpDir, err := os.MkdirTemp("", "oracle-2046B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2046B.go")
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
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
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
		{n: 1, arr: []int{5}},
		{n: 2, arr: []int{2, 1}},
		{n: 3, arr: []int{5, 5, 5}},
		{n: 4, arr: []int{3, 5, 1, 9}},
		{n: 5, arr: []int{1, 2, 2, 1, 4}},
		{n: 6, arr: []int{1, 2, 3, 6, 5, 4}},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	totalN := 0
	for len(tests) < 30 && totalN < 80000 {
		n := rng.Intn(8) + 2
		if totalN+n > 80000 {
			n = 80000 - totalN
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(20) + 1
		}
		tests = append(tests, testCase{n: n, arr: arr})
		totalN += n
	}
	for len(tests) < 35 {
		n := rng.Intn(1000) + 500
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1_000_000_000) + 1
		}
		tests = append(tests, testCase{n: n, arr: arr})
	}
	// One large stress case
	n := 100000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1_000_000_000) + 1
	}
	tests = append(tests, testCase{n: n, arr: arr})
	return tests
}

func normalizeOutput(out string) []string {
	return strings.Fields(out)
}

func compareOutputs(expected, actual string, count int) error {
	exp := normalizeOutput(expected)
	act := normalizeOutput(actual)
	if len(exp) != len(act) {
		return fmt.Errorf("token count mismatch: expected %d got %d", len(exp), len(act))
	}
	for i := 0; i < len(exp); i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at token %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
	if err := compareOutputs(expected, actual, len(normalizeOutput(expected))); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
