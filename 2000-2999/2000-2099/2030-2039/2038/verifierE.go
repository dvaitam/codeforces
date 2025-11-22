package main

import (
	"bytes"
	"fmt"
	"math"
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
	v []int
	h []int
}

func callerDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier path")
	}
	return filepath.Dir(file), nil
}

func buildOracle() (string, func(), error) {
	dir, err := callerDir()
	if err != nil {
		return "", nil, err
	}
	tmpDir, err := os.MkdirTemp("", "oracle-2038E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2038E.go")
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

func parseAnswer(out string) (float64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", fields[0], err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite output %v", val)
	}
	return val, nil
}

func compareAnswers(expected, actual float64) error {
	diff := math.Abs(expected - actual)
	den := math.Max(1, math.Abs(expected))
	if diff/den > 1e-6+1e-12 {
		return fmt.Errorf("expected %.12f, got %.12f (diff %.3e)", expected, actual, diff)
	}
	return nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.v {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, x := range tc.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(x))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		// Samples from the statement.
		{n: 2, v: []int{1, 2}, h: []int{2}},
		{n: 3, v: []int{3, 0, 0}, h: []int{6, 9}},
		{n: 5, v: []int{10, 0, 0, 0, 5}, h: []int{11, 1, 2, 5}},
		// Small custom cases.
		{n: 2, v: []int{0, 0}, h: []int{5}},
		{n: 3, v: []int{3, 3, 0}, h: []int{1, 5}},
		{n: 4, v: []int{7, 7, 7, 7}, h: []int{4, 4, 4}},
		{n: 3, v: []int{1_000_000, 1_000_000, 1_000_000}, h: []int{1_000_000, 1_000_000}},
	}
}

func makeEquilibriumCase(n int, rng *rand.Rand) testCase {
	heights := make([]int, n)
	idx := 0
	for idx < n {
		remaining := n - idx
		maxSize := 1 + rng.Intn(5)
		if maxSize > remaining {
			maxSize = remaining
		}
		size := 1 + rng.Intn(maxSize)
		level := rng.Intn(900000) + 1
		if rng.Intn(6) == 0 {
			level = 0
			size = 1
		}
		for c := 0; c < size && idx < n; c++ {
			heights[idx] = level
			idx++
		}
	}

	h := make([]int, n-1)
	for i := 0; i+1 < n; i++ {
		a, b := heights[i], heights[i+1]
		if a == b && a > 0 && rng.Intn(2) == 0 {
			h[i] = rng.Intn(a) + 1 // connect barrels at or below the waterline
		} else {
			minHi := max(a, b) + 1
			if minHi > 1_000_000 {
				minHi = 1_000_000
			}
			h[i] = minHi
			if minHi < 1_000_000 {
				h[i] += rng.Intn(1_000_000 - minHi + 1)
			}
		}
		if h[i] < 1 {
			h[i] = 1
		}
	}

	v := make([]int, n)
	copy(v, heights)
	return testCase{n: n, v: v, h: h}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	for i := 0; i < 15; i++ {
		n := rng.Intn(8) + 2
		tests = append(tests, makeEquilibriumCase(n, rng))
	}
	for i := 0; i < 10; i++ {
		n := rng.Intn(150) + 50
		tests = append(tests, makeEquilibriumCase(n, rng))
	}
	for i := 0; i < 5; i++ {
		n := rng.Intn(5000) + 500
		tests = append(tests, makeEquilibriumCase(n, rng))
	}
	tests = append(tests, makeEquilibriumCase(200000, rng))
	return tests
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (n=%d): %v\ninput:\n%s", idx+1, tc.n, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d (n=%d): %v\ninput:\n%s", idx+1, tc.n, err, input)
			os.Exit(1)
		}
		expAns, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d (n=%d): %v\n%s", idx+1, tc.n, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d (n=%d): %v\n%s", idx+1, tc.n, err, gotOut)
			os.Exit(1)
		}
		if err := compareAnswers(expAns, gotAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch (n=%d): %v\ninput:\n%s", idx+1, tc.n, err, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
