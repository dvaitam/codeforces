package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const valueLimit = 1 << 30

type testCase struct {
	n int
	x int
}

type testInput struct {
	text  string
	cases []testCase
}

type tokenReader struct {
	tokens []string
	idx    int
}

func newTokenReader(output string) *tokenReader {
	return &tokenReader{
		tokens: strings.Fields(output),
	}
}

func (tr *tokenReader) nextInt64() (int64, error) {
	if tr.idx >= len(tr.tokens) {
		return 0, fmt.Errorf("unexpected end of output")
	}
	token := tr.tokens[tr.idx]
	tr.idx++
	val, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", token)
	}
	return val, nil
}

func (tr *tokenReader) hasMore() bool {
	return tr.idx < len(tr.tokens)
}

func (tr *tokenReader) peek() string {
	if tr.hasMore() {
		return tr.tokens[tr.idx]
	}
	return ""
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2070-2079", "2072")
	tmp, err := os.CreateTemp("", "ref2072C")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2072C.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutput(output string, cases []testCase) ([][]int, error) {
	reader := newTokenReader(output)
	result := make([][]int, len(cases))
	for idx, tc := range cases {
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := reader.nextInt64()
			if err != nil {
				return nil, fmt.Errorf("test %d: %v", idx+1, err)
			}
			arr[i] = int(val)
		}
		result[idx] = arr
	}
	if reader.hasMore() {
		return nil, fmt.Errorf("extra tokens in output, starting with %q", reader.peek())
	}
	return result, nil
}

func mexValue(arr []int) int {
	seen := make([]bool, len(arr)+2)
	for _, v := range arr {
		if v >= 0 && v < len(seen) {
			seen[v] = true
		}
	}
	for i, ok := range seen {
		if !ok {
			return i
		}
	}
	return len(seen)
}

func checkArray(tc testCase, arr []int) (int, error) {
	if len(arr) != tc.n {
		return 0, fmt.Errorf("expected %d numbers, got %d", tc.n, len(arr))
	}
	orVal := 0
	for i, v := range arr {
		if v < 0 || v >= valueLimit {
			return 0, fmt.Errorf("element %d out of range at position %d", v, i+1)
		}
		orVal |= v
	}
	if orVal != tc.x {
		return 0, fmt.Errorf("OR is %d, expected %d", orVal, tc.x)
	}
	return mexValue(arr), nil
}

func makeInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
	}
	return testInput{
		text:  sb.String(),
		cases: append([]testCase(nil), cases...),
	}
}

func fixedInputs() []testInput {
	return []testInput{
		makeInput([]testCase{
			{n: 1, x: 69},
			{n: 7, x: 7},
			{n: 5, x: 7},
			{n: 7, x: 38},
			{n: 2, x: 3},
		}),
		makeInput([]testCase{
			{n: 10, x: 0},
			{n: 4, x: 0},
			{n: 1, x: 0},
		}),
		makeInput([]testCase{
			{n: 5, x: 1},
			{n: 8, x: 123456},
			{n: 12, x: 1<<29 - 1},
		}),
		makeInput([]testCase{
			{n: 200000, x: (1 << 29) - 7},
		}),
	}
}

func randomTest(rng *rand.Rand) testInput {
	t := rng.Intn(20) + 1
	remaining := 200000
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		casesLeft := t - i
		maxLen := remaining - (casesLeft - 1)
		if maxLen < 1 {
			maxLen = 1
		}
		n := 1 + rng.Intn(maxLen)
		if n > 200000 {
			n = 200000
		}
		remaining -= n
		x := rng.Intn(1 << 30)
		cases = append(cases, testCase{n: n, x: x})
	}
	return makeInput(cases)
}

func generateTests() []testInput {
	tests := fixedInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		refOut, err := runBinary(refBin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		refArrays, err := parseOutput(refOut, input.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		expectedMex := make([]int, len(refArrays))
		for i, arr := range refArrays {
			mex, err := checkArray(input.cases[i], arr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "reference invalid on test %d case %d: %v\n", idx+1, i+1, err)
				os.Exit(1)
			}
			expectedMex[i] = mex
		}

		candOut, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		candArrays, err := parseOutput(candOut, input.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		for i, arr := range candArrays {
			mex, err := checkArray(input.cases[i], arr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "test %d case %d invalid: %v\ninput n=%d x=%d\n", idx+1, i+1, err, input.cases[i].n, input.cases[i].x)
				os.Exit(1)
			}
			if mex != expectedMex[i] {
				fmt.Fprintf(os.Stderr, "test %d case %d failed: mex=%d expected=%d\ninput n=%d x=%d\ncandidate:%v\n", idx+1, i+1, mex, expectedMex[i], input.cases[i].n, input.cases[i].x, arr)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
