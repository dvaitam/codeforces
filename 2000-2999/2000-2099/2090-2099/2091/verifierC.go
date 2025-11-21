package main

import (
	"bufio"
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

const (
	refSource        = "2091C.go"
	tempOraclePrefix = "oracle-2091C-"
	randomTests      = 80
	maxN             = 2000
)

type testCase struct {
	n int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(rng, randomTests)...)

	for idx, tc := range tests {
		input := formatInput(tc)
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		expVals, err := parseAnswers(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotVals, err := parseAnswers(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if expVals == nil {
			if gotVals != nil {
				fmt.Fprintf(os.Stderr, "test %d failed: expected -1, got %v\n", idx+1, gotVals)
				fmt.Println("Input:")
				fmt.Print(input)
				fmt.Println("Output:")
				fmt.Print(candOut)
				os.Exit(1)
			}
			continue
		}

		if err := validatePermutation(tc.n, gotVals); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed validation: %v\n", idx+1, err)
			fmt.Println("Input:")
			fmt.Print(input)
			fmt.Println("Output:")
			fmt.Print(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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

func parseAnswers(out string, n int) ([]int, error) {
	out = strings.TrimSpace(out)
	if out == "-1" {
		return nil, nil
	}
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	perm := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		perm[i] = v
	}
	return perm, nil
}

func validatePermutation(n int, perm []int) error {
	if len(perm) != n {
		return fmt.Errorf("expected permutation of length %d, got %d", n, len(perm))
	}
	seen := make([]bool, n+1)
	for _, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range [1,%d]", v, n)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}
	arr := make([]int, n)
	copy(arr, perm)
	for shift := 0; shift < n; shift++ {
		count := 0
		for i := 0; i < n; i++ {
			if arr[i] == i+1 {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("shift %d has %d fixed points", shift, count)
		}
		rotate(arr)
	}
	return nil
}

func rotate(arr []int) {
	if len(arr) == 0 {
		return
	}
	last := arr[len(arr)-1]
	copy(arr[1:], arr[:len(arr)-1])
	arr[0] = last
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("1\n%d\n", tc.n)
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1},
		{n: 2},
		{n: 3},
		{n: 5},
		{n: 7},
	}
}
