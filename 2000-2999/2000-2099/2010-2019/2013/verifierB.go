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

const refSource2013B = "./2013B.go"

type testCase struct {
	n   int
	arr []int64
}

type testInput struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2013B)
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
	tmp, err := os.CreateTemp("", "2013B-ref-*")
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
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, tinyCases())
	rng := rand.New(rand.NewSource(2013))
	tests = append(tests, randomTest(rng, []int{3, 4, 5, 6}, 20))
	tests = append(tests, randomTest(rng, []int{100, 200}, 1000))
	tests = append(tests, randomTest(rng, []int{200000}, 1000000000))
	return tests
}

func sampleTest() testInput {
	cases := []testCase{
		{n: 2, arr: []int64{3, 2}},
		{n: 3, arr: []int64{2, 2, 8}},
		{n: 4, arr: []int64{1, 2, 4, 3}},
		{n: 5, arr: []int64{1, 2, 3, 4, 5}},
		{n: 5, arr: []int64{3, 2, 4, 5, 4}},
	}
	return buildTestInput(cases)
}

func tinyCases() testInput {
	cases := []testCase{
		{n: 2, arr: []int64{1, 1}},
		{n: 2, arr: []int64{1, 1000000000}},
		{n: 3, arr: []int64{1, 1, 1}},
		{n: 3, arr: []int64{1000000000, 1, 1000000000}},
	}
	return buildTestInput(cases)
}

func randomTest(rng *rand.Rand, sizes []int, maxVal int64) testInput {
	var cases []testCase
	sumN := 0
	for _, n := range sizes {
		if n < 2 {
			n = 2
		}
		sumN += n
		if sumN > 200000 {
			n -= sumN - 200000
			if n < 2 {
				break
			}
		}
		tc := testCase{
			n:   n,
			arr: make([]int64, n),
		}
		for i := 0; i < n; i++ {
			tc.arr[i] = rng.Int63n(maxVal) + 1
		}
		cases = append(cases, tc)
		if sumN >= 200000 {
			break
		}
	}
	return buildTestInput(cases)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return testInput{input: sb.String(), cases: cases}
}
