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
	n  int
	a  []int
	ks []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-212D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "212D.go")
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
	sb.Grow(64 + len(tc.a)*4 + len(tc.ks)*4)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.ks)))
	sb.WriteByte('\n')
	for i, v := range tc.ks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(output string, expected int) ([]float64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d values, got %d", expected, len(fields))
	}
	res := make([]float64, expected)
	for i, f := range fields {
		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q: %v", f, err)
		}
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return nil, fmt.Errorf("non-finite value %v", val)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:  3,
			a:  []int{3, 2, 1},
			ks: []int{1, 2, 3},
		},
		{
			n:  1,
			a:  []int{5},
			ks: []int{1},
		},
		{
			n:  5,
			a:  []int{5, 4, 5, 4, 5},
			ks: []int{1, 2, 4, 5, 3, 1},
		},
		{
			n:  6,
			a:  []int{1, 3, 2, 4, 2, 3},
			ks: []int{1, 2, 3, 4, 5, 6},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 60)
	for len(tests) < 60 {
		n := rng.Intn(40) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(400) + 1
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(1_000_000_000) + 1
		}
		m := rng.Intn(40) + 1
		if rng.Intn(5) == 0 {
			m = rng.Intn(200) + 1
		}
		ks := make([]int, m)
		for i := 0; i < m; i++ {
			ks[i] = rng.Intn(n) + 1
		}
		tests = append(tests, testCase{n: n, a: a, ks: ks})
	}
	return tests
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expectedOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		actualOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expectedVals, err := parseOutput(expectedOut, len(tc.ks))
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d: %v\n%s", idx+1, err, expectedOut)
			os.Exit(1)
		}
		actualVals, err := parseOutput(actualOut, len(tc.ks))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target produced invalid output on test %d: %v\noutput:\n%s", idx+1, err, actualOut)
			os.Exit(1)
		}
		for q := range expectedVals {
			diff := math.Abs(expectedVals[q] - actualVals[q])
			tol := 1e-9 * math.Max(1.0, math.Abs(expectedVals[q]))
			if diff > tol {
				fmt.Fprintf(os.Stderr, "test %d query %d mismatch: expected %.15f, got %.15f (diff %.3e > tol %.3e)\ninput:\n%s", idx+1, q+1, expectedVals[q], actualVals[q], diff, tol, input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed.")
}
