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

const (
	refSourceB = "2000-2999/2000-2099/2000-2009/2003/2003B.go"
	maxValue   = 100000
)

type testCase struct {
	cases [][]int
}

func (tc testCase) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tc.cases))
	for _, arr := range tc.cases {
		fmt.Fprintln(&b, len(arr))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
		input := tc.input()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tc.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
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
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2003B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceB))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		cases: [][]int{
			{2, 1},
			{1, 1, 2},
			{3, 1, 2},
			{3, 1, 2, 2, 3},
			{10, 2, 5, 2, 7, 9, 2, 5, 10, 7},
		},
	})

	tests = append(tests, testCase{cases: [][]int{
		{1000, 1000},
	}})

	tests = append(tests, testCase{cases: [][]int{
		makeRange(1, 100),
		repeatValue(50, 500),
	}})

	rng := rand.New(rand.NewSource(2003))
	for len(tests) < 40 {
		tests = append(tests, randomTest(rng))
	}

	tests = append(tests, hugeTest(100000))
	return tests
}

func makeRange(start, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = (start + i) % maxValue
		if arr[i] == 0 {
			arr[i] = maxValue
		}
	}
	return arr
}

func repeatValue(val, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = val
	}
	return arr
}

func randomTest(rng *rand.Rand) testCase {
	var cases [][]int
	remaining := 100000
	numCases := rng.Intn(5) + 1
	for i := 0; i < numCases; i++ {
		if remaining < 2 {
			break
		}
		n := rng.Intn(remaining-1) + 2
		if n > remaining {
			n = remaining
		}
		remaining -= n
		cases = append(cases, randomArray(rng, n))
	}
	if len(cases) == 0 {
		cases = append(cases, randomArray(rng, 2))
	}
	return testCase{cases: cases}
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(maxValue) + 1
	}
	return arr
}

func hugeTest(n int) testCase {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = (i % maxValue) + 1
	}
	return testCase{cases: [][]int{arr}}
}
