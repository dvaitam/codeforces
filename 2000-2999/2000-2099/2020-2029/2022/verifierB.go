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

const refSource2022B = "2000-2999/2000-2099/2020-2029/2022/2022B.go"

type testCase struct {
	n int
	x int64
	a []int64
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

	refBin, err := buildReference(refSource2022B)
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
		expected, err := parseInts(refOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseInts(userOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx := range test.cases {
			if expected[caseIdx] != actual[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d, got %d\ninput:\n%s",
					idx+1, caseIdx+1, expected[caseIdx], actual[caseIdx], formatSingleCase(test.cases[caseIdx]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2022B-ref-*")
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

func parseInts(out string, cases int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, cases)
	for i := 0; i < cases; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d answers, got %d (%v)", cases, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func formatSingleCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
	for i, v := range tc.a {
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
	tests = append(tests, edgeCases())
	rng := rand.New(rand.NewSource(2022))
	tests = append(tests, randomTest(rng, []int{5, 7, 9}, 10))
	tests = append(tests, randomTest(rng, []int{100, 150, 200}, 1000))
	tests = append(tests, randomTest(rng, []int{1000, 2000}, 100000))
	tests = append(tests, randomTest(rng, []int{500000}, 1000000000))
	return tests
}

func sampleTest() testInput {
	cases := []testCase{
		{n: 3, x: 2, a: []int64{3, 1, 2}},
		{n: 3, x: 3, a: []int64{2, 1, 3}},
		{n: 5, x: 3, a: []int64{2, 2, 1, 9, 2}},
		{n: 7, x: 4, a: []int64{2, 5, 3, 3, 5, 2, 5}},
	}
	return buildTestInput(cases)
}

func edgeCases() testInput {
	cases := []testCase{
		{n: 1, x: 1, a: []int64{1}},
		{n: 1, x: 10, a: []int64{100}},
		{n: 2, x: 1, a: []int64{1, 1}},
		{n: 2, x: 10, a: []int64{1000000000, 1000000000}},
	}
	return buildTestInput(cases)
}

func randomTest(rng *rand.Rand, sizes []int, maxA int64) testInput {
	var cases []testCase
	sumN := 0
	for _, n := range sizes {
		sumN += n
		if sumN > 500000 {
			n -= sumN - 500000
			if n <= 0 {
				break
			}
		}
		tc := testCase{
			n: n,
			x: int64(rng.Intn(10) + 1),
			a: make([]int64, n),
		}
		for i := 0; i < n; i++ {
			tc.a[i] = int64(rng.Int63n(maxA) + 1)
		}
		cases = append(cases, tc)
		if sumN >= 500000 {
			break
		}
	}
	return buildTestInput(cases)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return testInput{input: sb.String(), cases: cases}
}
