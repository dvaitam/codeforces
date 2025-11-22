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
	limitN   = 200000
	maxValue = 1_000_000_000
)

type testCase struct {
	a []int
}

type caseOutput struct {
	hasAnswer bool
	seq       []int
}

type testInput struct {
	name  string
	cases []testCase
}

func buildReference() (string, error) {
	path := "./2077D_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2077D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
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

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(len(tc.a)))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(output string, cases []testCase) ([]caseOutput, error) {
	tokens := strings.Fields(output)
	idx := 0
	res := make([]caseOutput, 0, len(cases))

	readInt := func(tok string) (int, error) {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return 0, err
		}
		if v < -1 || v > maxValue {
			return 0, fmt.Errorf("value %d out of bounds", v)
		}
		return int(v), nil
	}

	for caseIdx, tc := range cases {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing output", caseIdx+1)
		}
		tok := tokens[idx]
		idx++
		if tok == "-1" {
			res = append(res, caseOutput{hasAnswer: false})
			continue
		}
		k, err := readInt(tok)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid length %q (%v)", caseIdx+1, tok, err)
		}
		if k < 1 || k > len(tc.a) {
			return nil, fmt.Errorf("test %d: length %d out of range", caseIdx+1, k)
		}
		if idx+k > len(tokens) {
			return nil, fmt.Errorf("test %d: expected %d numbers for subsequence, only %d tokens remain", caseIdx+1, k, len(tokens)-idx)
		}
		seq := make([]int, k)
		for i := 0; i < k; i++ {
			val, err := readInt(tokens[idx+i])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid number %q (%v)", caseIdx+1, tokens[idx+i], err)
			}
			seq[i] = val
		}
		idx += k
		out := caseOutput{hasAnswer: true, seq: seq}
		if err := validateCase(tc, out); err != nil {
			return nil, fmt.Errorf("test %d: %v", caseIdx+1, err)
		}
		res = append(res, out)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected (%d tokens)", len(tokens)-idx)
	}
	return res, nil
}

func validateCase(tc testCase, out caseOutput) error {
	if !out.hasAnswer {
		return nil
	}
	k := len(out.seq)
	if k < 3 {
		return fmt.Errorf("subsequence too short (%d)", k)
	}
	if k > len(tc.a) {
		return fmt.Errorf("subsequence length %d exceeds n=%d", k, len(tc.a))
	}
	for i, v := range out.seq {
		if v < 1 || v > maxValue {
			return fmt.Errorf("value %d at position %d out of bounds", v, i+1)
		}
	}
	if !isSubsequence(out.seq, tc.a) {
		return fmt.Errorf("reported sequence is not a subsequence of input")
	}
	if !formsPolygon(out.seq) {
		return fmt.Errorf("reported sequence cannot form a polygon")
	}
	return nil
}

func isSubsequence(seq, arr []int) bool {
	pos := 0
	for _, v := range arr {
		if pos < len(seq) && seq[pos] == v {
			pos++
			if pos == len(seq) {
				return true
			}
		}
	}
	return pos == len(seq)
}

func formsPolygon(seq []int) bool {
	if len(seq) < 3 {
		return false
	}
	var sum int64
	maxVal := 0
	for _, v := range seq {
		sum += int64(v)
		if v > maxVal {
			maxVal = v
		}
	}
	return 2*int64(maxVal) < sum
}

func compareOutputs(expect, got []caseOutput) error {
	if len(expect) != len(got) {
		return fmt.Errorf("expected %d test cases, got %d", len(expect), len(got))
	}
	for i := range expect {
		e := expect[i]
		g := got[i]
		if e.hasAnswer != g.hasAnswer {
			if e.hasAnswer {
				return fmt.Errorf("test %d: expected a subsequence, got -1", i+1)
			}
			return fmt.Errorf("test %d: expected -1, but got an answer", i+1)
		}
		if !e.hasAnswer {
			continue
		}
		if len(e.seq) != len(g.seq) {
			return fmt.Errorf("test %d: expected length %d, got %d", i+1, len(e.seq), len(g.seq))
		}
		for j := range e.seq {
			if e.seq[j] != g.seq[j] {
				return fmt.Errorf("test %d: mismatch at position %d (expected %d, got %d)", i+1, j+1, e.seq[j], g.seq[j])
			}
		}
	}
	return nil
}

func sampleInput() testInput {
	cases := []testCase{
		{a: []int{1, 2, 3}},
		{a: []int{1, 4, 2, 3}},
		{a: []int{6, 5, 3, 2}},
		{a: []int{43, 99, 53, 22, 4}},
		{a: []int{54, 73, 23, 1}},
	}
	return testInput{name: "sample-like", cases: cases}
}

func edgeInput() testInput {
	cases := []testCase{
		{a: []int{1, 1, 1}},
		{a: []int{1, 1, 2}},
		{a: []int{5, 1, 1, 1, 1}},
		{a: []int{1000000000, 1000000000, 1000000000}},
		{a: []int{2, 3, 4, 10, 10, 10}},
		{a: []int{8, 3, 3, 3, 2, 2, 1}},
	}
	return testInput{name: "edge-cases", cases: cases}
}

func randomInput(rng *rand.Rand, maxCases, maxN int, forcePolygon bool) testInput {
	targetCases := rng.Intn(maxCases) + 1
	var cases []testCase
	totalN := 0
	for len(cases) < targetCases && totalN < limitN {
		remaining := limitN - totalN
		if remaining < 3 {
			break
		}
		upper := maxN
		if upper > remaining {
			upper = remaining
		}
		if upper < 3 {
			upper = 3
		}
		n := rng.Intn(upper-2) + 3
		tc := testCase{a: randomArray(rng, n, forcePolygon && rng.Intn(2) == 0)}
		cases = append(cases, tc)
		totalN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{a: []int{2, 3, 4}})
	}
	return testInput{name: "random", cases: cases}
}

func randomArray(rng *rand.Rand, n int, embedPolygon bool) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(maxValue) + 1
	}
	if !embedPolygon || n < 3 {
		return arr
	}
	m := rng.Intn(n-2) + 3
	val := rng.Intn(maxValue-1) + 1
	base := make([]int, m)
	for i := range base {
		base[i] = val
	}
	start := rng.Intn(n - m + 1)
	copy(arr[start:start+m], base)
	return arr
}

func largeStressInput(rng *rand.Rand) testInput {
	n := limitN
	arr := randomArray(rng, n, true)
	return testInput{name: "stress", cases: []testCase{{a: arr}}}
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []testInput{
		sampleInput(),
		edgeInput(),
	}
	for i := 0; i < 25; i++ {
		tests = append(tests, randomInput(rng, 5, 50, true))
	}
	for i := 0; i < 15; i++ {
		tests = append(tests, randomInput(rng, 5, 2000, true))
	}
	tests = append(tests, largeStressInput(rng))
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
	for idx, ti := range tests {
		input := buildInput(ti.cases)

		expectRaw, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\n", idx+1, ti.name, err)
			os.Exit(1)
		}
		expect, err := parseOutput(expectRaw, ti.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, ti.name, err, expectRaw)
			os.Exit(1)
		}

		gotRaw, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): runtime error: %v\ninput:\n%s", idx+1, ti.name, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(gotRaw, ti.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output: %v\ninput:\n%soutput:\n%s\n", idx+1, ti.name, err, input, gotRaw)
			os.Exit(1)
		}
		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, ti.name, err, input, expectRaw, gotRaw)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
