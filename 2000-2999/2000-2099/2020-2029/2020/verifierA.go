package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource2020A = "2000-2999/2000-2099/2020-2029/2020/2020A.go"

type testCase struct {
	n int64
	k int64
}

type testInput struct {
	cases []testCase
}

func (ti testInput) toString() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.k)
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.toString()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2020A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource2020A))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())

	rng := rand.New(rand.NewSource(2020))
	for len(tests) < 40 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTests() testInput {
	return testInput{cases: []testCase{
		{n: 5, k: 2},
		{n: 3, k: 5},
		{n: 16, k: 4},
		{n: 100, k: 3},
		{n: 6492, k: 10},
		{n: 10, k: 1},
	}}
}

func edgeTests() testInput {
	return testInput{cases: []testCase{
		{n: 1, k: 1},
		{n: 1_000_000_000, k: 1},
		{n: 1_000_000_000, k: 1_000_000_000},
		{n: 999_999_937, k: 2},
	}}
}

func randomTest(rng *rand.Rand) testInput {
	numCases := rng.Intn(10) + 1
	cases := make([]testCase, numCases)
	for i := 0; i < numCases; i++ {
		n := rng.Int63n(1_000_000_000) + 1
		k := rng.Int63n(1_000_000_000) + 1
		cases[i] = testCase{n: n, k: k}
	}
	return testInput{cases: cases}
}
