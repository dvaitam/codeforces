package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource2000C = "./2000C.go"

type testCase2000C struct {
	name   string
	input  string
	output []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference2000C()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests2000C()
	for idx, tc := range tests {
		refOut, err := runProgram2000C(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refAns := parseAnswerLines(refOut)
		if len(refAns) != len(tc.output) {
			fmt.Fprintf(os.Stderr, "reference output mismatch in verifier for test %d (%s)\ninput:\n%sreference:\n%s",
				idx+1, tc.name, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram2000C(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns := parseAnswerLines(candOut)
		if len(candAns) != len(refAns) {
			fmt.Fprintf(os.Stderr, "candidate printed %d answers but expected %d on test %d (%s)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				len(candAns), len(refAns), idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}

		for i := range refAns {
			if !equalYesNo(refAns[i], candAns[i]) {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on query %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference2000C() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2000C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2000C.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2000C)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram2000C(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswerLines(out string) []string {
	var res []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		res = append(res, line)
	}
	return res
}

func equalYesNo(a, b string) bool {
	return normalizeYesNo(a) == normalizeYesNo(b)
}

func normalizeYesNo(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "YES", "Y", "TRUE":
		return "YES"
	case "NO", "N", "FALSE":
		return "NO"
	default:
		return s
	}
}

func buildTests2000C() []testCase2000C {
	tests := []testCase2000C{
		manualTest2000C("single_char_yes", []int{7}, []string{"a"}, []string{"YES"}),
		manualTest2000C("single_char_no", []int{3}, []string{"bb"}, []string{"NO"}),
		manualTest2000C("example_from_statement", []int{3, 5, 2, 1, 3}, []string{"abfda", "afbfa"}, []string{"YES", "NO"}),
		manualTest2000C("multiple_patterns", []int{1, 2, 1, 3}, []string{"abab", "abca", "zzxz", "xyyx"}, []string{"YES", "NO", "YES", "NO"}),
		manualTest2000C("distinct_numbers", []int{1, 2, 3, 4}, []string{"abcd", "aabb"}, []string{"YES", "NO"}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest2000C(rng, i))
	}
	return tests
}

func manualTest2000C(name string, pattern []int, queries []string, answers []string) testCase2000C {
	input := formatInput2000C(pattern, queries)
	return testCase2000C{
		name:   name,
		input:  input,
		output: answers,
	}
}

func formatInput2000C(pattern []int, queries []string) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(pattern)))
	for i, val := range pattern {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", val))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(q)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTest2000C(rng *rand.Rand, idx int) testCase2000C {
	n := rng.Intn(6) + 1
	patternVals := make([]int, n)
	for i := 0; i < n; i++ {
		patternVals[i] = rng.Intn(3)
	}
	m := rng.Intn(4) + 1
	queries := make([]string, m)
	for i := 0; i < m; i++ {
		var sb strings.Builder
		for j := 0; j < n; j++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		queries[i] = sb.String()
	}
	input := formatInput2000C(patternVals, queries)
	return testCase2000C{
		name:   fmt.Sprintf("random_%d", idx+1),
		input:  input,
		output: make([]string, len(queries)), // placeholder; actual answers obtained from reference.
	}
}
