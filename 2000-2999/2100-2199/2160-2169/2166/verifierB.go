package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	a int64
	b int64
	n int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTestSuite(rng)
	input := buildInput(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2166B.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2166B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, expectedOutputs int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	values, err := parseAnswers(out, expectedOutputs)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", expected, len(fields), out)
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func buildTestSuite(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	for len(tests) < 500 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{a: 28, b: 1, n: 6},
		{a: 9, b: 6, n: 2},
		{a: 10, b: 3, n: 1},
		{a: 10, b: 1, n: 10},
		{a: 9, b: 2, n: 1},
		{a: 5, b: 5, n: 6},
		{a: 6, b: 2, n: 7},
		{a: 9, b: 1, n: 9},
		{a: 3, b: 2, n: 6},
		{a: 8, b: 1, n: 7},
		{a: 8, b: 2, n: 4},
		{a: 1, b: 1, n: 1},
		{a: 5, b: 5, n: 1},
		{a: 1000000000, b: 1, n: 1000000000},
		{a: 1000000000, b: 1000000000, n: 1000000000},
		{a: 1000000000, b: 500000000, n: 1},
		{a: 2, b: 1, n: 200000000},
	}
}

func randomTest(rng *rand.Rand) testCase {
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(a) + 1
	n := rng.Int63n(1_000_000_000) + 1
	return testCase{a: a, b: b, n: n}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for idx, tc := range tests {
		if idx > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d %d", tc.a, tc.b, tc.n))
	}
	sb.WriteByte('\n')
	return sb.String()
}
