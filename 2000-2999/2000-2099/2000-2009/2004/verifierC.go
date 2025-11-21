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
	refSourceC = "2000-2999/2000-2099/2000-2009/2004/2004C.go"
	maxValue   = int64(1e9)
	maxSumN    = 200000
)

type caseData struct {
	n int
	k int64
	a []int64
}

type testCase struct {
	cases []caseData
}

func (tc testCase) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tc.cases))
	for _, cs := range tc.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.k)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.FormatInt(v, 10))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
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
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2004C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceC))
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

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
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

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{cases: []caseData{
		{n: 2, k: 5, a: []int64{1, 10}},
		{n: 3, k: 0, a: []int64{15, 12, 10}},
		{n: 5, k: 12, a: []int64{3, 1, 2, 4, 6}},
		{n: 2, k: 46, a: []int64{6, 9}},
	}})

	tests = append(tests, testCase{cases: []caseData{
		{n: 2, k: 0, a: []int64{1, 1}},
		{n: 2, k: 10, a: []int64{1, 1}},
		{n: 3, k: 1, a: []int64{1, 100, 2}},
	}})

	rng := rand.New(rand.NewSource(2004))
	for len(tests) < 40 {
		tests = append(tests, randomBundle(rng, 50000))
	}

	tests = append(tests, stressCase(maxSumN))
	return tests
}

func randomBundle(rng *rand.Rand, budget int) testCase {
	var cases []caseData
	sumN := 0
	for sumN < budget {
		n := rng.Intn(10) + 2
		if sumN+n > budget {
			break
		}
		sumN += n
		k := rng.Int63n(1_000_000_000)
		a := make([]int64, n)
		for i := range a {
			a[i] = rng.Int63n(maxValue) + 1
		}
		cases = append(cases, caseData{n: n, k: k, a: a})
	}
	if len(cases) == 0 {
		cases = append(cases, caseData{n: 2, k: 0, a: []int64{1, 1}})
	}
	return testCase{cases: cases}
}

func stressCase(total int) testCase {
	a := make([]int64, total)
	for i := 0; i < total; i++ {
		a[i] = int64((i % 1000) + 1)
	}
	return testCase{cases: []caseData{{n: total, k: 1_000_000_000, a: a}}}
}
