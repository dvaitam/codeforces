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
	"time"
)

const refSourceI = "2000-2999/2000-2099/2040-2049/2045/2045I.go"

type testCase struct {
	n int
	m int
	a []int
}

type testInput struct {
	cases []testCase
}

func (ti testInput) buildInput() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.m)
		for i, v := range cs.a {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
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
		input := tc.buildInput()
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
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput: n=%d m=%d a=%v\nreference: %d\ncandidate: %d\n",
					idx+1, i+1, tc.cases[i].n, tc.cases[i].m, tc.cases[i].a, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2045I-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceI))
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

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomBundle(rng, 10, 1000, 1000, 300000))
	tests = append(tests, randomBundle(rng, 5, 10000, 10000, 300000))
	tests = append(tests, randomBundle(rng, 2, 300000, 300000, 300000))

	return tests
}

func sampleTests() testInput {
	return testInput{cases: []testCase{
		{n: 5, m: 4, a: []int{3, 2, 1, 3, 2}},
		{n: 3, m: 3, a: []int{1, 1, 1}},
	}}
}

func edgeTests() testInput {
	return testInput{cases: []testCase{
		{n: 1, m: 1, a: []int{1}},
		{n: 2, m: 2, a: []int{1, 2}},
		{n: 4, m: 2, a: []int{2, 1, 2, 1}},
	}}
}

func randomBundle(rng *rand.Rand, maxCases, maxN, maxM, limit int) testInput {
	var cases []testCase
	sumN := 0
	for len(cases) < maxCases && sumN < limit {
		n := rng.Intn(maxN) + 1
		if sumN+n > limit {
			n = limit - sumN
			if n <= 0 {
				break
			}
		}
		m := rng.Intn(maxM) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(m) + 1
		}
		cases = append(cases, testCase{n: n, m: m, a: a})
		sumN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, m: 1, a: []int{1}})
	}
	return testInput{cases: cases}
}
