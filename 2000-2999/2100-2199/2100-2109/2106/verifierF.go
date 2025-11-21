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
	s string
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2106F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2106F.go")
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
	sb.Grow(totalN + len(tests)*32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
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
		{n: 3, s: "100"}, // sample
		{n: 4, s: "1010"},
		{n: 7, s: "1011001"},
		{n: 1, s: "0"},
		{n: 1, s: "1"},
		{n: 2, s: "01"},
		{n: 2, s: "10"},
		{n: 5, s: "11111"},
		{n: 5, s: "00000"},
	}
}

func randomString(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func randomTests(totalN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 256)
	used := 0
	for used < totalN {
		maxN := totalN - used
		n := rng.Intn(min(2000, maxN)) + 1
		s := randomString(n, rng)
		tests = append(tests, testCase{n: n, s: s})
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	const nLimit = 200_000
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

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	expectedAns, err := parseOutput(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle output invalid: %v\n%s", err, expOut)
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
