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
	p []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2138C1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC1")
	cmd := exec.Command("go", "build", "-o", outPath, "2138C1.go")
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
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}
	sb.Grow(totalN*4 + len(tests)*32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		res[i] = v
	}
	return res, nil
}

func compareAnswers(expected, actual []int64) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return fmt.Errorf("answer %d mismatch: expected %d, got %d", i+1, expected[i], actual[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, k: 0, p: []int{1}},                // all ones
		{n: 2, k: 2, p: []int{1}},                // all zeros
		{n: 2, k: 1, p: []int{1}},                // balanced
		{n: 5, k: 0, p: []int{1, 2, 3, 4}},       // chain all ones
		{n: 5, k: 5, p: []int{1, 2, 3, 4}},       // chain all zeros
		{n: 7, k: 3, p: []int{1, 1, 2, 2, 3, 3}}, // small tree
	}
}

func randomTree(n int, rng *rand.Rand) []int {
	p := make([]int, n-1)
	for i := 2; i <= n; i++ {
		p[i-2] = rng.Intn(i-1) + 1
	}
	return p
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 128)
	used := 0
	for used < totalN {
		remain := totalN - used
		n := rng.Intn(min(1000, remain)) + 2
		k := rng.Intn(n + 1)
		p := randomTree(n, rng)
		tests = append(tests, testCase{n: n, k: k, p: p})
		used += n
	}
	return tests
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
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

	tests := deterministicTests()
	const nLimit = 50_000 // within t<=50, n<=1000 gives margin
	used := totalN(tests)
	if used < nLimit {
		tests = append(tests, randomTests(nLimit-used)...)
	}

	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	expectedAns, err := parseOutput(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}
	actualAns, err := parseOutput(actOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\n%s", err, actOut)
		os.Exit(1)
	}

	if err := compareAnswers(expectedAns, actualAns); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
