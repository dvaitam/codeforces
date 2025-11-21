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

type testCase struct {
	n   int
	arr []int64
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
	return &tokenReader{tokens: strings.Fields(output)}
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
	refDir := filepath.Join("2000-2999", "2100-2199", "2160-2169", "2169")
	tmp, err := os.CreateTemp("", "ref2169C")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2169C.go")
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

func parseOutput(output string, t int) ([]int64, error) {
	reader := newTokenReader(output)
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		val, err := reader.nextInt64()
		if err != nil {
			return nil, fmt.Errorf("test %d: %v", i+1, err)
		}
		ans[i] = val
	}
	if reader.hasMore() {
		return nil, fmt.Errorf("extra tokens in output starting with %q", reader.peek())
	}
	return ans, nil
}

func computeMaxSum(tc testCase) int64 {
	n := tc.n
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + tc.arr[i-1]
	}
	const inf = int64(1 << 60)
	minG := inf
	best := int64(0)
	for i := 1; i <= n; i++ {
		ii := int64(i)
		g := ii*ii - ii - pref[i-1]
		if g < minG {
			minG = g
		}
		candidate := ii*ii + ii - pref[i] - minG
		if candidate > best {
			best = candidate
		}
	}
	return pref[n] + best
}

func makeInput(cases []testCase) testInput {
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
	copyCases := make([]testCase, len(cases))
	for i, tc := range cases {
		arrCopy := make([]int64, len(tc.arr))
		copy(arrCopy, tc.arr)
		copyCases[i] = testCase{n: tc.n, arr: arrCopy}
	}
	return testInput{text: sb.String(), cases: copyCases}
}

func fixedInputs() []testInput {
	return []testInput{
		makeInput([]testCase{
			{n: 3, arr: []int64{2, 5, 1}},
			{n: 4, arr: []int64{4, 4, 1, 3}},
		}),
		makeInput([]testCase{
			{n: 1, arr: []int64{0}},
			{n: 2, arr: []int64{0, 0}},
			{n: 10, arr: []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		}),
		makeInput([]testCase{
			{n: 5, arr: []int64{0, 10, 0, 10, 0}},
			{n: 6, arr: []int64{6, 6, 6, 6, 6, 6}},
		}),
		makeInput([]testCase{
			{n: 200000, arr: buildRepeated(200000, 400000)},
		}),
	}
}

func buildRepeated(n int, val int64) []int64 {
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = val
	}
	return arr
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
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(2*n + 1))
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return makeInput(cases)
}

func generateTests() []testInput {
	tests := fixedInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 100 {
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
		expect, err := parseOutput(refOut, len(input.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input.text)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, len(input.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for caseIdx, tc := range input.cases {
			expectedValue := computeMaxSum(tc)
			if got[caseIdx] != expectedValue || expect[caseIdx] != expectedValue {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch. expected %d, candidate %d, reference %d\ninput: n=%d\n", idx+1, caseIdx+1, expectedValue, got[caseIdx], expect[caseIdx], tc.n)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
