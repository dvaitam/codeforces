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

type testCase struct {
	input string
	t     int
	nVals []int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", out, "2151C.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, nVals []int) ([][]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(nVals) {
		return nil, fmt.Errorf("expected %d lines, got %d", len(nVals), len(lines))
	}
	results := make([][]string, len(nVals))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != nVals[i] {
			return nil, fmt.Errorf("case %d: expected %d values, got %d", i+1, nVals[i], len(fields))
		}
		results[i] = fields
	}
	return results, nil
}

func deterministicTests() []testCase {
	tests := []testCase{
		formatTest([]testInput{
			{n: 1, times: []int64{1, 2}},
			{n: 2, times: []int64{1, 2, 3, 4}},
		}),
		formatTest([]testInput{
			{n: 3, times: []int64{1, 2, 3, 4, 5, 6}},
		}),
		formatTest([]testInput{
			{n: 2, times: []int64{1, 3, 7, 10}},
			{n: 1, times: []int64{5, 10}},
		}),
	}
	return tests
}

type testInput struct {
	n     int
	times []int64
}

func formatTest(cases []testInput) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	nVals := make([]int, len(cases))
	for i, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.times {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		nVals[i] = tc.n
	}
	return testCase{input: sb.String(), t: len(cases), nVals: nVals}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalN := 0
	for len(tests) < count && totalN < 200000 {
		t := rnd.Intn(4) + 1
		cases := make([]testInput, t)
		nVals := make([]int, t)
		for i := 0; i < t; i++ {
			if totalN >= 200000 {
				t = i
				break
			}
			n := rnd.Intn(2000) + 1
			if totalN+n > 200000 {
				n = 200000 - totalN
			}
			totalN += n
			times := make([]int64, 2*n)
			cur := int64(rnd.Intn(10) + 1)
			for j := 0; j < 2*n; j++ {
				cur += int64(rnd.Intn(10) + 1)
				times[j] = cur
			}
			cases[i] = testInput{n: n, times: times}
			nVals[i] = n
		}
		formatted := formatTest(cases)
		formatted.nVals = nVals
		tests = append(tests, formatted)
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, tc.nVals)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.nVals)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.t; caseIdx++ {
			for i := 0; i < tc.nVals[caseIdx]; i++ {
				if gotVals[caseIdx][i] != expVals[caseIdx][i] {
					fmt.Fprintf(os.Stderr, "case %d test %d k=%d mismatch: expected %s got %s\n", idx+1, caseIdx+1, i+1, expVals[caseIdx][i], gotVals[caseIdx][i])
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
