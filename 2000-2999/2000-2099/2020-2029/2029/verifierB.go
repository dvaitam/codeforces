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

const refSource2029B = "2029B.go"

type testCase struct {
	n int
	s string
	r string
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

	refBin, err := buildReference(refSource2029B)
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
		expected, err := parseAnswers(refOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseAnswers(userOut, len(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx, ans := range actual {
			exp := expected[caseIdx]
			if ans != exp {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %s, got %s\ninput:\n%s",
					idx+1, caseIdx+1, exp, ans, formatSingleCase(test.cases[caseIdx]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2029B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	sourcePath := filepath.Join(".", source)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), sourcePath)
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

func parseAnswers(out string, cases int) ([]string, error) {
	reader := strings.NewReader(out)
	res := make([]string, cases)
	for i := 0; i < cases; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d answers, got %d (%v)", cases, i, err)
		}
		res[i] = strings.ToUpper(strings.TrimSpace(res[i]))
		if res[i] != "YES" && res[i] != "NO" {
			return nil, fmt.Errorf("invalid answer %q", res[i])
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
	sb.WriteString(fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.r))
	return sb.String()
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, edgeTest())
	rng := rand.New(rand.NewSource(2029))
	tests = append(tests, randomTest(rng, []int{2, 3, 4, 5, 6}, 20))
	tests = append(tests, randomTest(rng, []int{10, 15, 20, 30, 40}, 200))
	tests = append(tests, randomTest(rng, []int{500, 600, 700}, 2000))
	tests = append(tests, randomTest(rng, []int{100000}, 100000))
	return tests
}

func sampleTest() testInput {
	cases := []testCase{
		{n: 2, s: "11", r: "0"},
		{n: 2, s: "01", r: "1"},
		{n: 3, s: "101", r: "010"},
		{n: 3, s: "111", r: "100"},
		{n: 6, s: "101001", r: "010110"},
		{n: 8, s: "10010010", r: "0010010"},
	}
	return buildTestInput(cases)
}

func edgeTest() testInput {
	cases := []testCase{
		{n: 2, s: "01", r: "0"},
		{n: 2, s: "10", r: "1"},
		{n: 3, s: "010", r: "11"},
		{n: 3, s: "101", r: "00"},
		{n: 4, s: "0101", r: "101"},
	}
	return buildTestInput(cases)
}

func randomTest(rng *rand.Rand, sizes []int, maxN int) testInput {
	var cases []testCase
	totalN := 0
	for _, n := range sizes {
		if n < 2 {
			n = 2
		}
		totalN += n
		if totalN > maxN {
			n -= totalN - maxN
			if n < 2 {
				break
			}
		}
		s := randomBinaryString(rng, n)
		r := randomBinaryString(rng, n-1)
		cases = append(cases, testCase{n: n, s: s, r: r})
		if totalN >= maxN {
			break
		}
	}
	return buildTestInput(cases)
}

func randomBinaryString(rng *rand.Rand, n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = '0'
		} else {
			bytes[i] = '1'
		}
	}
	return string(bytes)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.r)
		sb.WriteByte('\n')
	}
	return testInput{input: sb.String(), cases: cases}
}
