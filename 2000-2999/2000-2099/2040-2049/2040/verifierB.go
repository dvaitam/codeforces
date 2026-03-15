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
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	for idx, tc := range tests {
		refVals := solveAll(tc.input)

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, fmtVals(refVals), candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, fmtVals(refVals), candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

// minOps computes the correct answer for CF 2040/B.
// After the first type-1 op covers 1 cell, each subsequent op can extend
// coverage from c to 2*c+2 (place a 1 at position c + (c+2) = 2c+2,
// then the segment [1, 2c+2] has c+1 ones out of 2c+2 = ceil((2c+2)/2)).
func minOps(n int64) int64 {
	ops := int64(0)
	covered := int64(0)
	for covered < n {
		if covered == 0 {
			covered = 1
		} else {
			covered = covered*2 + 2
		}
		ops++
	}
	return ops
}

func solveAll(input string) []int64 {
	fields := strings.Fields(strings.TrimSpace(input))
	t, _ := strconv.Atoi(fields[0])
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		n, _ := strconv.ParseInt(fields[i+1], 10, 64)
		res[i] = minOps(n)
	}
	return res
}

func fmtVals(vals []int64) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(parts, " ")
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, input string) ([]int64, error) {
	fields := strings.Fields(strings.TrimSpace(output))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	lines := strings.Fields(strings.TrimSpace(input))
	if len(lines) == 0 {
		return nil, fmt.Errorf("invalid input")
	}
	t, err := strconv.Atoi(lines[0])
	if err != nil || t <= 0 {
		return nil, fmt.Errorf("invalid test case count")
	}
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		res[i] = val
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("single", []int64{1}),
		makeManual("small", []int64{1, 2, 3, 4}),
		makeManual("powers", []int64{1, 2, 3, 7, 8, 9}),
		makeManual("examples", []int64{1, 2, 4, 20}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManual(name string, ns []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ns)))
	for _, n := range ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{
		name:  name,
		input: sb.String(),
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := randInt64(rng, 1, 1_000_000_000_000)
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: sb.String(),
	}
}

func randInt64(rng *rand.Rand, low, high int64) int64 {
	return low + rng.Int63n(high-low+1)
}
