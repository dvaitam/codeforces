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

type testInput struct {
	ms []int
}

func (ti testInput) build() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(ti.ms))
	for _, m := range ti.ms {
		fmt.Fprintln(&b, m)
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.build()

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(tc.ms))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		// Basic validation: answers should be non-negative integers
		for i := range candVals {
			if candVals[i] < 0 {
				fmt.Fprintf(os.Stderr, "test %d case %d: negative answer %d\n", idx+1, i+1, candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func commandFor(path string) *exec.Cmd {
	switch {
	case strings.HasSuffix(path, ".go"):
		return exec.Command("go", "run", path)
	case strings.HasSuffix(path, ".py"):
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
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
	tests = append(tests, testInput{ms: []int{2, 5, 9}})
	tests = append(tests, testInput{ms: []int{1, 2, 3, 4, 5}})
	tests = append(tests, testInput{ms: []int{1000}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const limit = 10000
	sum := 0
	var cur []int
	for sum < limit {
		m := rng.Intn(1000) + 1
		if sum+m > limit {
			break
		}
		cur = append(cur, m)
		sum += m
		if len(cur) >= 100 {
			tests = append(tests, testInput{ms: cur})
			cur = nil
		}
	}
	if len(cur) > 0 {
		tests = append(tests, testInput{ms: cur})
	}
	return tests
}
