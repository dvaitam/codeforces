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

const refSource1773A = "./1773A.go"

type testCase struct {
	n int
	a []int
}

type testInput struct {
	input string
	cases []testCase
}

type outputCase struct {
	possible bool
	p        []int
	q        []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource1773A)
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
		refCases, err := parseOutputs(refOut, test.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runProgram(candidate, test.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", idx+1, err, test.input)
			os.Exit(1)
		}
		userCases, err := parseOutputs(userOut, test.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, test.input, userOut)
			os.Exit(1)
		}

		for caseIdx, tc := range test.cases {
			refCase := refCases[caseIdx]
			userCase := userCases[caseIdx]
			if !refCase.possible && userCase.possible {
				fmt.Fprintf(os.Stderr, "test %d case %d: candidate claims Possible but reference Impossible\n", idx+1, caseIdx+1)
				fmt.Fprintf(os.Stderr, "input case:\n%s", formatCase(tc))
				os.Exit(1)
			}
			if userCase.possible {
				if err := validateSolution(tc, userCase.p, userCase.q); err != nil {
					fmt.Fprintf(os.Stderr, "test %d case %d invalid solution: %v\n", idx+1, caseIdx+1, err)
					fmt.Fprintf(os.Stderr, "input case:\n%s", formatCase(tc))
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "1773A-ref-*")
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
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, tests []testCase) ([]outputCase, error) {
	reader := strings.NewReader(out)
	res := make([]outputCase, len(tests))
	for i, tc := range tests {
		var verdict string
		if _, err := fmt.Fscan(reader, &verdict); err != nil {
			return nil, fmt.Errorf("case %d: failed to read verdict: %v", i+1, err)
		}
		switch strings.ToLower(verdict) {
		case "impossible":
			res[i] = outputCase{possible: false}
		case "possible":
			p := make([]int, tc.n)
			q := make([]int, tc.n)
			for j := 0; j < tc.n; j++ {
				if _, err := fmt.Fscan(reader, &p[j]); err != nil {
					return nil, fmt.Errorf("case %d: failed to read permutation p: %v", i+1, err)
				}
			}
			for j := 0; j < tc.n; j++ {
				if _, err := fmt.Fscan(reader, &q[j]); err != nil {
					return nil, fmt.Errorf("case %d: failed to read permutation q: %v", i+1, err)
				}
			}
			res[i] = outputCase{possible: true, p: p, q: q}
		default:
			return nil, fmt.Errorf("case %d: invalid verdict %q", i+1, verdict)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected (starts with %s)", extra)
	}
	return res, nil
}

func validateSolution(tc testCase, p, q []int) error {
	n := tc.n
	if len(p) != n || len(q) != n {
		return fmt.Errorf("permutations must have length %d", n)
	}
	if err := checkPermutation(p); err != nil {
		return fmt.Errorf("invalid permutation p: %v", err)
	}
	if err := checkPermutation(q); err != nil {
		return fmt.Errorf("invalid permutation q: %v", err)
	}
	for i := 0; i < n; i++ {
		if p[i] == i+1 {
			return fmt.Errorf("p has a fixed point at position %d", i+1)
		}
		if q[i] == i+1 {
			return fmt.Errorf("q has a fixed point at position %d", i+1)
		}
	}

	// compute r = p ∘ q
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = p[q[i]-1]
	}

	// apply r to positions of a
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		pos[tc.a[i]] = i + 1
	}
	for i := 0; i < n; i++ {
		card := tc.a[r[i]-1]
		if card != i+1 {
			return fmt.Errorf("after shuffles card %d is %d", i+1, card)
		}
	}
	return nil
}

func checkPermutation(p []int) error {
	n := len(p)
	seen := make([]bool, n+1)
	for _, v := range p {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range", v)
		}
		if seen[v] {
			return fmt.Errorf("value %d repeated", v)
		}
		seen[v] = true
	}
	return nil
}

func formatCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("n = ")
	sb.WriteString(fmt.Sprint(tc.n))
	sb.WriteString("\na = ")
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests, sampleTests())
	tests = append(tests, smallEdgeTests())
	rng := rand.New(rand.NewSource(1773))
	tests = append(tests, randomTests(rng, []int{3, 4, 5, 6}))
	tests = append(tests, randomTests(rng, []int{50, 60}))
	tests = append(tests, randomTests(rng, []int{1000}))
	return tests
}

func sampleTests() testInput {
	cases := []testCase{
		{n: 2, a: []int{1, 2}},
		{n: 3, a: []int{1, 2, 3}},
		{n: 4, a: []int{2, 1, 4, 3}},
		{n: 5, a: []int{5, 1, 4, 2, 3}},
	}
	return buildTestInput(cases)
}

func smallEdgeTests() testInput {
	cases := []testCase{
		{n: 1, a: []int{1}},
		{n: 2, a: []int{1, 2}},
		{n: 2, a: []int{2, 1}},
		{n: 3, a: []int{2, 3, 1}},
	}
	return buildTestInput(cases)
}

func randomTests(rng *rand.Rand, sizes []int) testInput {
	var cases []testCase
	for _, n := range sizes {
		if n <= 0 {
			continue
		}
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			perm[i] = i + 1
		}
		for i := n - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			perm[i], perm[j] = perm[j], perm[i]
		}
		cases = append(cases, testCase{n: n, a: perm})
	}
	return buildTestInput(cases)
}

func buildTestInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return testInput{input: sb.String(), cases: cases}
}
