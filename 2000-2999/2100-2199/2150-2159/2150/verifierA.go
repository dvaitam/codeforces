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
	refBin       = "./2150A_ref.bin"
	maxCoord     = int64(1_000_000_000)
	globalLimitN = 100000
	globalLimitM = 100000
)

type testCase struct {
	s     string
	black []int64
}

type testInput struct {
	name  string
	cases []testCase
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBin, "2150A.go")
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
		sb.WriteString(strconv.Itoa(len(tc.s)))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(len(tc.black)))
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		for i, v := range tc.black {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t testInput) ([][]int64, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([][]int64, 0, len(t.cases))

	readInt := func(tok string) (int64, error) {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	for caseIdx := range t.cases {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing output", caseIdx+1)
		}
		k64, err := readInt(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid k %q (%v)", caseIdx+1, tokens[idx], err)
		}
		if k64 < 0 {
			return nil, fmt.Errorf("test %d: negative k", caseIdx+1)
		}
		k := int(k64)
		idx++
		if idx+k > len(tokens) {
			return nil, fmt.Errorf("test %d: expected %d numbers, only %d tokens remain", caseIdx+1, k, len(tokens)-idx)
		}
		arr := make([]int64, k)
		for i := 0; i < k; i++ {
			val, err := readInt(tokens[idx+i])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid number %q (%v)", caseIdx+1, tokens[idx+i], err)
			}
			if val < 1 || val > maxCoord {
				return nil, fmt.Errorf("test %d: value %d out of bounds", caseIdx+1, val)
			}
			if i > 0 && val <= arr[i-1] {
				return nil, fmt.Errorf("test %d: values not strictly increasing at position %d", caseIdx+1, i+1)
			}
			arr[i] = val
		}
		idx += k
		res = append(res, arr)
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected (%d tokens)", len(tokens)-idx)
	}
	return res, nil
}

func compareOutputs(exp, got [][]int64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("expected %d test cases, got %d", len(exp), len(got))
	}
	for i := range exp {
		if len(exp[i]) != len(got[i]) {
			return fmt.Errorf("test %d: expected %d black cells, got %d", i+1, len(exp[i]), len(got[i]))
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
		name: "sample-like",
		cases: []testCase{
			{s: "BAB", black: []int64{2, 5}},
			{s: "ABA", black: []int64{1, 4, 9, 10}},
			{s: "ABABB", black: []int64{1, 7}},
		},
	}
}

func smallEdge() testInput {
	return testInput{
		name: "edge-small",
		cases: []testCase{
			{s: "A", black: []int64{1}},
			{s: "B", black: []int64{1}},
			{s: "AAAA", black: []int64{5}},
			{s: "BBBB", black: []int64{1, 3, 5, 7}},
		},
	}
}

func randomString(rng *rand.Rand, n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			bytes[i] = 'A'
		} else {
			bytes[i] = 'B'
		}
	}
	return string(bytes)
}

func randomBlack(rng *rand.Rand, m int64) []int64 {
	if m <= 0 {
		return nil
	}
	black := make([]int64, m)
	start := rng.Int63n(maxCoord/2) + 1
	black[0] = start
	for i := 1; i < int(m); i++ {
		// small step to keep values within bounds
		delta := rng.Int63n(9) + 1
		val := black[i-1] + delta
		if val > maxCoord {
			val = maxCoord - int64(len(black)-i)
		}
		black[i] = val
	}
	return black
}

func randomInput(rng *rand.Rand, name string, maxCases int, maxN, maxM int) testInput {
	cases := make([]testCase, 0, maxCases)
	sumN, sumM := 0, 0
	for len(cases) < maxCases && sumN < globalLimitN && sumM < globalLimitM {
		n := rng.Intn(maxN) + 1
		m := rng.Intn(maxM) + 1
		if sumN+n > globalLimitN {
			n = globalLimitN - sumN
		}
		if sumM+m > globalLimitM {
			m = globalLimitM - sumM
		}
		if n <= 0 || m <= 0 {
			break
		}
		tc := testCase{
			s:     randomString(rng, n),
			black: randomBlack(rng, int64(m)),
		}
		cases = append(cases, tc)
		sumN += n
		sumM += m
	}
	return testInput{name: name, cases: cases}
}

func stressInput(rng *rand.Rand) testInput {
	n := 50000
	m := 50000
	tc := testCase{
		s:     randomString(rng, n),
		black: randomBlack(rng, int64(m)),
	}
	return testInput{name: "stress", cases: []testCase{tc}}
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []testInput{
		sampleInput(),
		smallEdge(),
		randomInput(rng, "random-small", 8, 10, 10),
		randomInput(rng, "random-medium", 5, 500, 500),
		randomInput(rng, "random-large", 3, 5000, 5000),
		stressInput(rng),
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
