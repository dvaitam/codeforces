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

const refSourceB = "2000-2999/2000-2099/2020-2029/2024/2024B.go"

type caseData struct {
	n int
	k int64
	a []int64
}

type testInput struct {
	cases []caseData
}

func (ti testInput) buildInput() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.cases))
	for _, cs := range ti.cases {
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
			if candVals[i] != refVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2024B-ref-*")
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
	tests = append(tests, smallEdgeTests())

	rng := rand.New(rand.NewSource(2024))
	for len(tests) < 40 {
		tests = append(tests, randomBundle(rng, 200000))
	}

	tests = append(tests, stressTest())
	return tests
}

func sampleTests() testInput {
	return testInput{cases: []caseData{
		{n: 2, k: 1, a: []int64{1, 1}},
		{n: 2, k: 2, a: []int64{1, 2}},
		{n: 3, k: 4, a: []int64{2, 1, 3}},
		{n: 2, k: 1, a: []int64{3, 10}},
		{n: 2, k: 1_000_000_000, a: []int64{1_000_000_000, 1_000_000_000}},
	}}
}

func smallEdgeTests() testInput {
	return testInput{cases: []caseData{
		{n: 1, k: 1, a: []int64{1}},
		{n: 1, k: 1_000_000_000, a: []int64{1_000_000_000}},
		{n: 5, k: 5, a: []int64{1, 1, 1, 1, 1}},
		{n: 5, k: 3, a: []int64{1, 2, 3, 4, 5}},
	}}
}

func randomBundle(rng *rand.Rand, limit int) testInput {
	var (
		cases []caseData
		used  int
	)
	for used < limit {
		n := rng.Intn(5000) + 1
		if used+n > limit {
			n = limit - used
		}
		if n == 0 {
			break
		}
		a := make([]int64, n)
		var sum int64
		for i := range a {
			a[i] = rng.Int63n(1_000_000_000) + 1
			sum += a[i]
		}
		var k int64
		if sum > 0 {
			k = randInt63(rng, 1, sum)
		} else {
			k = 0
		}
		cases = append(cases, caseData{n: n, k: k, a: a})
		used += n
		if rng.Intn(4) == 0 {
			break
		}
	}
	if len(cases) == 0 {
		cases = append(cases, caseData{n: 1, k: 1, a: []int64{1}})
	}
	return testInput{cases: cases}
}

func stressTest() testInput {
	n := 200000
	a := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		a[i] = 1_000_000_000
		sum += a[i]
	}
	return testInput{cases: []caseData{
		{n: n, k: sum / 2, a: a},
	}}
}

func randInt63(rng *rand.Rand, low, high int64) int64 {
	if low > high {
		low, high = high, low
	}
	if low == high {
		return low
	}
	return low + rng.Int63n(high-low+1)
}
