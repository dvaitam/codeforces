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

const refSource = "2000-2999/2100-2199/2140-2149/2148/2148F.go"

type testCase struct {
	name    string
	input   string
	t       int
	arrays  [][]int
	maxLens []int
}

func runProgram(bin, input string) (string, error) {
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
	return stdout.String(), nil
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	result := make([]string, 0, expected)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, line)
	}
	if len(result) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(result))
	}
	return result, nil
}

func manualTests() []testCase {
	return []testCase{
		buildTestCase("single_short", [][]int{{1, 2, 3}}, 1),
		buildTestCase("two_arrays", [][]int{{1, 3, 5}, {2, 4}}, 1),
		buildTestCase("multiple", [][]int{{1}, {2, 2}, {3, 3, 3}}, 1),
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		t := rng.Intn(3) + 1
		cases := make([][]int, t)
		for j := 0; j < t; j++ {
			length := rng.Intn(5) + 1
			arr := make([]int, length)
			for k := 0; k < length; k++ {
				arr[k] = rng.Intn(10)
			}
			cases[j] = arr
		}
		tests = append(tests, buildTestCase(fmt.Sprintf("random_%d", i+1), cases, t))
	}
	return tests
}

func buildTestCase(name string, cases [][]int, t int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	maxLens := make([]int, t)
	for i := 0; i < t; i++ {
		n := len(cases[i])
		maxLens[i] = n
		sb.WriteString(fmt.Sprintf("%d %d\n", n, n))
		for j, v := range cases[i] {
			sb.WriteString(strconv.Itoa(v))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:    name,
		input:   sb.String(),
		t:       t,
		arrays:  cases,
		maxLens: maxLens,
	}
}

func parseSequence(line string, expectedLen int) ([]int, error) {
	fields := strings.Fields(line)
	if len(fields) != expectedLen {
		return nil, fmt.Errorf("expected %d integers, got %d in line %q", expectedLen, len(fields), line)
	}
	result := make([]int, expectedLen)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		result[i] = val
	}
	return result, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}

	refBin, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	tests := append(manualTests(), randomTests(100)...)
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refLines, err := parseOutput(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candLines, err := parseOutput(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.t; caseIdx++ {
			refSeq, err := parseSequence(refLines[caseIdx], tc.maxLens[caseIdx])
			if err != nil {
				fmt.Fprintf(os.Stderr, "reference sequence parse error on test %d case %d: %v\n", idx+1, caseIdx+1, err)
				os.Exit(1)
			}
			candSeq, err := parseSequence(candLines[caseIdx], tc.maxLens[caseIdx])
			if err != nil {
				fmt.Fprintf(os.Stderr, "candidate sequence parse error on test %d case %d: %v\n", idx+1, caseIdx+1, err)
				os.Exit(1)
			}
			for i := range refSeq {
				if refSeq[i] != candSeq[i] {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d mismatch at position %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
						idx+1, tc.name, caseIdx+1, i+1, refSeq[i], candSeq[i], tc.input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
