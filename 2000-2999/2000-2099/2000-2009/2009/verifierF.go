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

const refSource2009F = "2000-2999/2000-2099/2000-2009/2009/2009F.go"

type query struct {
	l int64
	r int64
}

type testCase struct {
	n       int
	q       int
	a       []int64
	queries []query
}

type testInput struct {
	input string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2009F)
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
		expected, err := parseResults(refOut, totalQueries(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runCandidate(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		actual, err := parseResults(userOut, totalQueries(test.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx, tc := range test.cases {
			for qIdx := 0; qIdx < tc.q; qIdx++ {
				pos := prefixIndex(test.cases, caseIdx, qIdx)
				if expected[pos] != actual[pos] {
					fmt.Fprintf(os.Stderr, "test %d case %d query %d mismatch: expected %d, got %d\n", idx+1, caseIdx+1, qIdx+1, expected[pos], actual[pos])
					fmt.Fprintf(os.Stderr, "case data:\n%s", formatSingleCase(tc))
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2009F-ref-*")
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

func parseResults(out string, count int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, count)
	for i := 0; i < count; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return nil, fmt.Errorf("expected %d answers, got %d (%v)", count, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func totalQueries(cases []testCase) int {
	sum := 0
	for _, tc := range cases {
		sum += tc.q
	}
	return sum
}

func prefixIndex(cases []testCase, caseIdx int, qIdx int) int {
	idx := 0
	for i := 0; i < caseIdx; i++ {
		idx += cases[i].q
	}
	idx += qIdx
	return idx
}

func formatSingleCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.l, q.r))
	}
	return sb.String()
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTest())
	tests = append(tests, edgeCases())
	rng := rand.New(rand.NewSource(2009))
	tests = append(tests, randomTest(rng, []int{3, 5, 7}, []int{3, 4, 5}, 10))
	tests = append(tests, randomTest(rng, []int{100, 200}, []int{100, 200}, 1000000))
	tests = append(tests, randomTest(rng, []int{200000}, []int{200000}, 1000000))
	return tests
}

func sampleTest() testInput {
	cases := []testCase{
		{
			n:       3,
			q:       3,
			a:       []int64{1, 2, 3},
			queries: []query{{1, 9}, {3, 5}, {8, 8}},
		},
		{
			n:       5,
			q:       4,
			a:       []int64{8, 8, 5, 5, 4},
			queries: []query{{8, 3}, {2, 4}, {1, 14}, {3, 7}},
		},
	}
	return buildTestInput(cases)
}

func edgeCases() testInput {
	cases := []testCase{
		{
			n:       1,
			q:       3,
			a:       []int64{5},
			queries: []query{{1, 1}, {1, 1}, {1, 1}},
		},
		{
			n:       2,
			q:       4,
			a:       []int64{1, 2},
			queries: []query{{1, 1}, {2, 3}, {3, 4}, {1, 4}},
		},
	}
	return buildTestInput(cases)
}

func randomTest(rng *rand.Rand, ns []int, qs []int, maxVal int64) testInput {
	var cases []testCase
	sumN := 0
	sumQ := 0
	for i := 0; i < len(ns) && i < len(qs); i++ {
		n := ns[i]
		q := qs[i]
		sumN += n
		sumQ += q
		if sumN > 200000 {
			n -= sumN - 200000
			if n <= 0 {
				break
			}
		}
		if sumQ > 200000 {
			q -= sumQ - 200000
			if q <= 0 {
				break
			}
		}
		tc := testCase{
			n:       n,
			q:       q,
			a:       make([]int64, n),
			queries: make([]query, q),
		}
		for j := 0; j < n; j++ {
			tc.a[j] = rng.Int63n(maxVal) + 1
		}
		maxIndex := int64(n) * int64(n)
		for j := 0; j < q; j++ {
			l := rng.Int63n(maxIndex) + 1
			r := rng.Int63n(maxIndex) + 1
			if l > r {
				l, r = r, l
			}
			tc.queries[j] = query{l: l, r: r}
		}
		cases = append(cases, tc)
		if sumN >= 200000 || sumQ >= 200000 {
			break
		}
	}
	return buildTestInput(cases)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", tc.a[i]))
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q.l, q.r))
		}
	}
	return testInput{input: sb.String(), cases: cases}
}
