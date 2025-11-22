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

const (
	refBin   = "./2129F1_ref.bin"
	maxN     = 845
	maxCases = 40
)

type testCase struct {
	perm []int
}

type testInput struct {
	name  string
	cases []testCase
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBin, "2129F1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(t testInput) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(t.cases)))
	sb.WriteByte('\n')
	for _, tc := range t.cases {
		sb.WriteString(strconv.Itoa(len(tc.perm)))
		sb.WriteByte('\n')
		for i, v := range tc.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t testInput) ([][]int, error) {
	tokens := strings.Fields(out)
	idx := 0
	result := make([][]int, 0, len(t.cases))

	readInt := func(tok string) (int, error) {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return 0, err
		}
		if v < 1 {
			return 0, fmt.Errorf("value %d out of bounds", v)
		}
		return int(v), nil
	}

	for caseIdx, tc := range t.cases {
		n := len(tc.perm)
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing output", caseIdx+1)
		}
		if tokens[idx] != "!" {
			return nil, fmt.Errorf("test %d: expected '!' but got %q", caseIdx+1, tokens[idx])
		}
		idx++
		if idx+n > len(tokens) {
			return nil, fmt.Errorf("test %d: expected %d integers, only %d tokens remain", caseIdx+1, n, len(tokens)-idx)
		}
		perm := make([]int, n)
		seen := make([]bool, n+1)
		for i := 0; i < n; i++ {
			val, err := readInt(tokens[idx+i])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q (%v)", caseIdx+1, tokens[idx+i], err)
			}
			if val < 1 || val > n {
				return nil, fmt.Errorf("test %d: value %d out of range 1..%d", caseIdx+1, val, n)
			}
			if seen[val] {
				return nil, fmt.Errorf("test %d: value %d appears multiple times", caseIdx+1, val)
			}
			seen[val] = true
			perm[i] = val
		}
		for v := 1; v <= n; v++ {
			if !seen[v] {
				return nil, fmt.Errorf("test %d: value %d missing from permutation", caseIdx+1, v)
			}
		}
		idx += n
		result = append(result, perm)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected (%d tokens)", len(tokens)-idx)
	}
	return result, nil
}

func compareOutputs(exp, got [][]int) error {
	if len(exp) != len(got) {
		return fmt.Errorf("expected %d cases, got %d", len(exp), len(got))
	}
	for i := range exp {
		if len(exp[i]) != len(got[i]) {
			return fmt.Errorf("test %d: expected length %d, got %d", i+1, len(exp[i]), len(got[i]))
		}
		for j := range exp[i] {
			if exp[i][j] != got[i][j] {
				return fmt.Errorf("test %d: mismatch at position %d (expected %d, got %d)", i+1, j+1, exp[i][j], got[i][j])
			}
		}
	}
	return nil
}

func sampleInput() testInput {
	return testInput{
		name: "sample",
		cases: []testCase{
			{perm: []int{3, 1, 2}},
			{perm: []int{2, 1}},
		},
	}
}

func smallCases() testInput {
	return testInput{
		name: "small",
		cases: []testCase{
			{perm: []int{1, 2}},
			{perm: []int{2, 1, 3, 4}},
			{perm: []int{4, 3, 2, 1}},
			{perm: []int{1, 3, 2, 4, 5, 6}},
		},
	}
}

func randomPerm(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	return perm
}

func randomInput(rng *rand.Rand, name string, cases int, nMin, nMax int) testInput {
	all := make([]testCase, 0, cases)
	for i := 0; i < cases; i++ {
		n := rng.Intn(nMax-nMin+1) + nMin
		if n < 2 {
			n = 2
		}
		if n > maxN {
			n = maxN
		}
		all = append(all, testCase{perm: randomPerm(rng, n)})
	}
	return testInput{name: name, cases: all}
}

func largeStress(rng *rand.Rand) testInput {
	stressCases := make([]testCase, 0, 3)
	sizes := []int{maxN, maxN - 1, maxN/2 + 1}
	for _, n := range sizes {
		stressCases = append(stressCases, testCase{perm: randomPerm(rng, n)})
	}
	return testInput{name: "stress", cases: stressCases}
}

func mixedInput(rng *rand.Rand) testInput {
	cases := make([]testCase, 0)
	sizes := []int{3, 4, 5, 10, 20, 50, 100, 200}
	for _, n := range sizes {
		cases = append(cases, testCase{perm: randomPerm(rng, n)})
	}
	return testInput{name: "mixed", cases: cases}
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testInput{
		sampleInput(),
		smallCases(),
		mixedInput(rng),
		randomInput(rng, "random-small", 10, 2, 20),
		randomInput(rng, "random-mid", 8, 50, 200),
		randomInput(rng, "random-large", 6, 300, maxN),
		largeStress(rng),
	}
	// Ensure total cases do not exceed limit.
	for i := range tests {
		if len(tests[i].cases) > maxCases {
			tests[i].cases = tests[i].cases[:maxCases]
		}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, t := range tests {
		input := buildInput(t)

		expRaw, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", idx+1, t.name, err)
			os.Exit(1)
		}
		exp, err := parseOutput(expRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, t.name, err, expRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): runtime error: %v\ninput:\n%s", idx+1, t.name, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(gotRaw, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, t.name, err, input, gotRaw)
			os.Exit(1)
		}

		if err := compareOutputs(exp, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, t.name, err, input, expRaw, gotRaw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
