package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource2002A = "./2002A.go"

type testCase struct {
	n int64
	m int64
	k int64
}

type testInput struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2002A)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, test := range tests {
		refOut, err := runProgram(refBin, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		expected, err := parseOutputs(refOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseOutputs(userOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx := range test.cases {
			if expected[caseIdx] != actual[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d, got %d\ncase input:\n%s",
					idx+1, caseIdx+1, expected[caseIdx], actual[caseIdx], formatCase(test.cases[caseIdx]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2002A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(source))
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, cases int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, cases)
	for i := 0; i < cases; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d numbers, got %d (%v)", cases, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func formatCase(tc testCase) string {
	return fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.m, tc.k)
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, edgeTests())
	rng := rand.New(rand.NewSource(2002))
	tests = append(tests, randomTest(rng, []int{5, 7, 9}, 10000))
	tests = append(tests, randomTest(rng, []int{1000, 2000}, 10000))
	tests = append(tests, randomTest(rng, []int{1000}, 10000))
	return tests
}

func sampleTests() testInput {
	cases := []testCase{
		{n: 3, m: 3, k: 2},
		{n: 5, m: 1, k: 10000},
		{n: 7, m: 3, k: 4},
		{n: 3, m: 2, k: 7},
		{n: 8, m: 9, k: 6},
		{n: 2, m: 5, k: 4},
	}
	return buildTestInput(cases)
}

func edgeTests() testInput {
	cases := []testCase{
		{n: 1, m: 1, k: 1},
		{n: 1, m: 10000, k: 1},
		{n: 10000, m: 1, k: 1},
		{n: 10000, m: 10000, k: 10000},
		{n: 9999, m: 10000, k: 2},
	}
	return buildTestInput(cases)
}

func randomTest(rng *rand.Rand, sizes []int, maxVal int64) testInput {
	var cases []testCase
	for _, t := range sizes {
		for i := 0; i < t; i++ {
			cases = append(cases, testCase{
				n: rng.Int63n(maxVal) + 1,
				m: rng.Int63n(maxVal) + 1,
				k: rng.Int63n(maxVal) + 1,
			})
		}
	}
	return buildTestInput(cases)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	}
	return testInput{input: sb.String(), cases: cases}
}
