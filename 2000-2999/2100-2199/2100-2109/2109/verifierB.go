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

const refSource = "2109B.go"

type testCase struct {
	n, m int64
	a, b int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	want, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, want[i], got[i])
			fmt.Fprintf(os.Stderr, "n=%d m=%d a=%d b=%d\n", tc.n, tc.m, tc.a, tc.b)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2109B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Join(dir, refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutputs(out string, expected int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output starting at token %q", tokens[expected])
	}
	res := make([]int64, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		if val < 0 {
			return nil, fmt.Errorf("negative answer %d", val)
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.a, tc.b))
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	add := func(tc testCase) {
		tests = append(tests, tc)
	}

	for _, tc := range sampleTests() {
		add(tc)
	}

	// Small edge positions.
	add(testCase{n: 2, m: 2, a: 1, b: 1})
	add(testCase{n: 2, m: 2, a: 2, b: 2})
	add(testCase{n: 2, m: 3, a: 1, b: 3})
	add(testCase{n: 3, m: 2, a: 3, b: 1})
	add(testCase{n: 4, m: 4, a: 2, b: 3})

	// Large extremes near constraints.
	add(testCase{n: 1_000_000_000, m: 1_000_000_000, a: 1, b: 1})
	add(testCase{n: 1_000_000_000, m: 1_000_000_000, a: 1_000_000_000, b: 1_000_000_000})
	add(testCase{n: 1_000_000_000, m: 17, a: 999_999_999, b: 9})
	add(testCase{n: 17, m: 1_000_000_000, a: 9, b: 123_456_789})
	add(testCase{n: 500_000_000, m: 750_000_000, a: 250_000_000, b: 375_000_000})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 500 {
		n := randRange(rng, 2, 1_000_000_000)
		m := randRange(rng, 2, 1_000_000_000)
		a := randRange(rng, 1, n)
		b := randRange(rng, 1, m)
		add(testCase{n: n, m: m, a: a, b: b})
	}

	return tests
}

func sampleTests() []testCase {
	return []testCase{
		{n: 2, m: 2, a: 1, b: 1},
		{n: 3, m: 3, a: 2, b: 2},
		{n: 2, m: 7, a: 1, b: 4},
		{n: 2, m: 7, a: 2, b: 2},
		{n: 8, m: 9, a: 4, b: 6},
		{n: 9, m: 9, a: 5, b: 5},
		{n: 2, m: 20, a: 2, b: 11},
		{n: 22, m: 99, a: 20, b: 70},
	}
}

func randRange(rng *rand.Rand, l, r int64) int64 {
	return l + rng.Int63n(r-l+1)
}
