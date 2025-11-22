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

const refSourceD = "2000-2999/2000-2099/2040-2049/2042/2042D.go"

type interval struct {
	l int
	r int
}

type testCase struct {
	ints []interval
}

type testBundle struct {
	cases []testCase
}

func (tb testBundle) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(tb.cases))
	for _, tc := range tb.cases {
		fmt.Fprintln(&b, len(tc.ints))
		for _, it := range tc.ints {
			fmt.Fprintf(&b, "%d %d\n", it.l, it.r)
		}
	}
	return b.String()
}

func main() {
	var candidate string
	if len(os.Args) == 2 {
		candidate = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		candidate = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, bundle := range tests {
		input := bundle.input()

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, bundle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, bundle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if len(refVals[i]) != len(candVals[i]) {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d: expected %d numbers got %d\ninput:\n%sreference: %v\ncandidate: %v\n", idx+1, i+1, len(refVals[i]), len(candVals[i]), input, refVals[i], candVals[i])
				os.Exit(1)
			}
			for j := range refVals[i] {
				if refVals[i][j] != candVals[i][j] {
					fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d position %d\ninput:\n%sreference: %v\ncandidate: %v\n", idx+1, i+1, j+1, input, refVals[i], candVals[i])
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2042D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceD))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
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

func parseOutputs(out string, bundle testBundle) ([][]int64, error) {
	fields := strings.Fields(out)
	pos := 0
	res := make([][]int64, len(bundle.cases))
	for i, tc := range bundle.cases {
		if pos >= len(fields) {
			return nil, fmt.Errorf("missing output for case %d", i+1)
		}
		if len(fields)-pos < len(tc.ints) {
			return nil, fmt.Errorf("case %d: expected %d numbers, got %d", i+1, len(tc.ints), len(fields)-pos)
		}
		cur := make([]int64, len(tc.ints))
		for j := range cur {
			val, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d value %d: %v", i+1, j+1, err)
			}
			cur[j] = val
		}
		res[i] = cur
		pos += len(tc.ints)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing tokens: expected %d consumed, got %d", pos, len(fields))
	}
	return res, nil
}

func buildTests() []testBundle {
	var tests []testBundle

	// Simple and edge cases.
	tests = append(tests, testBundle{cases: []testCase{{ints: []interval{{l: 1, r: 1}}}}})
	tests = append(tests, testBundle{cases: []testCase{{ints: []interval{{l: 3, r: 8}, {l: 2, r: 5}, {l: 4, r: 5}}}}})
	tests = append(tests, testBundle{cases: []testCase{{ints: []interval{{l: 42, r: 42}, {l: 1, r: 1000000000}}}}})
	tests = append(tests, testBundle{cases: []testCase{{ints: []interval{{l: 2, r: 4}, {l: 2, r: 4}, {l: 2, r: 4}}}}})
	tests = append(tests, testBundle{cases: []testCase{{ints: []interval{{l: 6, r: 10}, {l: 3, r: 10}, {l: 3, r: 7}, {l: 5, r: 7}, {l: 4, r: 4}}}}})

	rng := rand.New(rand.NewSource(2042))

	// Random mid-size cases keeping sum n moderate.
	for i := 0; i < 6; i++ {
		var ints []interval
		n := rng.Intn(8) + 5
		for j := 0; j < n; j++ {
			l, r := randomInterval(rng, 50)
			ints = append(ints, interval{l: l, r: r})
		}
		tests = append(tests, testBundle{cases: []testCase{{ints: ints}}})
	}

	// Structured nested/disjoint mix.
	tests = append(tests, testBundle{cases: []testCase{
		{ints: []interval{
			{l: 1, r: 100},
			{l: 1, r: 50},
			{l: 25, r: 75},
			{l: 60, r: 90},
			{l: 10, r: 10},
			{l: 101, r: 200},
			{l: 150, r: 200},
		}},
	}})

	// Stress near constraints: total n close to 2e5.
	tests = append(tests, testBundle{cases: []testCase{stressCase(rng, 200000)}})

	return tests
}

func randomInterval(rng *rand.Rand, maxCoord int) (int, int) {
	a := rng.Intn(maxCoord) + 1
	b := rng.Intn(maxCoord) + 1
	if a > b {
		a, b = b, a
	}
	return a, b
}

func stressCase(rng *rand.Rand, n int) testCase {
	ints := make([]interval, n)
	for i := range ints {
		// Use wide coordinate range including extremes.
		l, r := randomInterval(rng, 1_000_000_000)
		ints[i] = interval{l: l, r: r}
	}
	return testCase{ints: ints}
}
