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
	refSourceH = "2000-2999/2000-2099/2000-2009/2008/2008H.go"
	maxValue   = 100000
)

type caseData struct {
	n       int
	q       int
	a       []int
	queries []int
}

type testCase struct {
	cases []caseData
}

func (tc testCase) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tc.cases))
	for _, cs := range tc.cases {
		fmt.Fprintf(&b, "%d %d\n", cs.n, cs.q)
		for i, v := range cs.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
		for _, x := range cs.queries {
			fmt.Fprintln(&b, x)
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate")
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
		refVals, err := parseOutputs(refOut, totalAnswers(tc))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(refVals))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d at query %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func totalAnswers(tc testCase) int {
	sum := 0
	for _, cs := range tc.cases {
		sum += cs.q
	}
	return sum
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2008H-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceH))
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
	tests = append(tests, sampleTest())
	tests = append(tests, edgeTest())

	rng := rand.New(rand.NewSource(2008))
	for len(tests) < 30 {
		tests = append(tests, randomBundle(rng, 50000, 50000))
	}

	tests = append(tests, stressTest())
	return tests
}

func sampleTest() testCase {
	return testCase{cases: []caseData{
		{
			n: 5, q: 5,
			a:       []int{1, 2, 3, 4, 5},
			queries: []int{1, 2, 3, 4, 5},
		},
		{
			n: 6, q: 3,
			a:       []int{1, 2, 6, 4, 1, 3},
			queries: []int{2, 1, 5},
		},
	}}
}

func edgeTest() testCase {
	return testCase{cases: []caseData{
		{n: 1, q: 4, a: []int{1}, queries: []int{1, 1, 1, 1}},
		{n: 2, q: 2, a: []int{2, 2}, queries: []int{1, 2}},
	}}
}

func randomBundle(rng *rand.Rand, budgetN, budgetQ int) testCase {
	var cases []caseData
	remainingN := budgetN
	remainingQ := budgetQ
	for remainingN > 0 && remainingQ > 0 {
		maxN := min(remainingN, 1000)
		if maxN == 0 {
			break
		}
		n := rng.Intn(maxN) + 1
		maxQ := min(remainingQ, 1000)
		if maxQ == 0 {
			break
		}
		q := rng.Intn(maxQ) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(n) + 1
		}
		queries := make([]int, q)
		for i := range queries {
			queries[i] = rng.Intn(n) + 1
		}
		cases = append(cases, caseData{n: n, q: q, a: a, queries: queries})
		remainingN -= n
		remainingQ -= q
		if rng.Intn(4) == 0 {
			break
		}
	}
	if len(cases) == 0 {
		cases = append(cases, caseData{n: 1, q: 1, a: []int{1}, queries: []int{1}})
	}
	return testCase{cases: cases}
}

func stressTest() testCase {
	n := maxValue
	q := maxValue
	a := make([]int, n)
	for i := range a {
		a[i] = (i % n) + 1
	}
	queries := make([]int, q)
	for i := range queries {
		queries[i] = (i % n) + 1
	}
	return testCase{cases: []caseData{{n: n, q: q, a: a, queries: queries}}}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
