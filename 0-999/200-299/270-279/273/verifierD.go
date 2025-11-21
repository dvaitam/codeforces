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

const (
	refSource        = "273D.go"
	tempOraclePrefix = "oracle-273D-"
	randomTests      = 120
	maxN             = 150
	maxM             = 150
)

type testCase struct {
	n int
	m int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestCases(randomTests, rng)...)

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
		expStr, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (n=%d, m=%d): %v\n", idx+1, tc.n, tc.m, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		gotStr, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (n=%d, m=%d): %v\n", idx+1, tc.n, tc.m, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		exp, err := parseSingleInt(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, expStr)
			os.Exit(1)
		}
		got, err := parseSingleInt(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, gotStr)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d (n=%d, m=%d)\n", idx+1, exp, got, tc.n, tc.m)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
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
	return strings.TrimSpace(stdout.String()), nil
}

func parseSingleInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value %q is not an integer", fields[0])
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("unexpected extra tokens after result")
	}
	return val, nil
}

func deterministicTests() []testCase {
	var tests []testCase
	for n := 1; n <= 6; n++ {
		for m := 1; m <= 6; m++ {
			tests = append(tests, testCase{n: n, m: m})
		}
	}
	edges := []testCase{
		{n: 1, m: maxM},
		{n: maxN, m: 1},
		{n: maxN, m: maxM},
		{n: maxN, m: maxM - 1},
		{n: maxN - 1, m: maxM},
		{n: 75, m: 150},
		{n: 150, m: 75},
		{n: 149, m: 149},
	}
	tests = append(tests, edges...)
	return tests
}

func randomTestCases(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for len(tests) < count {
		n := rng.Intn(maxN) + 1
		m := rng.Intn(maxM) + 1
		tests = append(tests, testCase{n: n, m: m})
	}
	return tests
}
