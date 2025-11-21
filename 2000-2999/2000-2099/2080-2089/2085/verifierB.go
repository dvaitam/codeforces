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
	arr []int
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
	toks := strings.Fields(output)
	return &tokenReader{tokens: toks}
}

func (tr *tokenReader) nextInt() (int, error) {
	if tr.idx >= len(tr.tokens) {
		return 0, fmt.Errorf("unexpected end of output")
	}
	token := tr.tokens[tr.idx]
	tr.idx++
	val, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", token)
	}
	return val, nil
}

func (tr *tokenReader) hasMore() bool {
	return tr.idx < len(tr.tokens)
}

func (tr *tokenReader) peek() string {
	if !tr.hasMore() {
		return ""
	}
	return tr.tokens[tr.idx]
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

func mexSegment(seg []int) int {
	seen := make([]bool, len(seg)+2)
	for _, v := range seg {
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

func applyOperation(arr []int, l, r int) []int {
	mex := mexSegment(arr[l : r+1])
	newArr := make([]int, 0, len(arr)-(r-l))
	newArr = append(newArr, arr[:l]...)
	newArr = append(newArr, mex)
	newArr = append(newArr, arr[r+1:]...)
	return newArr
}

func verifyCase(idx int, tc testCase, reader *tokenReader) error {
	k, err := reader.nextInt()
	if err != nil {
		return fmt.Errorf("missing operation count: %v", err)
	}
	if k < 0 || k > tc.n {
		return fmt.Errorf("invalid number of operations k=%d, expected 0<=k<=%d", k, tc.n)
	}

	cur := make([]int, len(tc.arr))
	copy(cur, tc.arr)

	for op := 0; op < k; op++ {
		l, err := reader.nextInt()
		if err != nil {
			return fmt.Errorf("operation %d: missing l: %v", op+1, err)
		}
		r, err := reader.nextInt()
		if err != nil {
			return fmt.Errorf("operation %d: missing r: %v", op+1, err)
		}
		if len(cur) < 2 {
			return fmt.Errorf("operation %d: array length is %d, cannot apply further operations", op+1, len(cur))
		}
		if l < 1 || r < 1 || l >= r || r > len(cur) {
			return fmt.Errorf("operation %d: invalid indices l=%d r=%d for length %d", op+1, l, r, len(cur))
		}
		cur = applyOperation(cur, l-1, r-1)
	}

	if len(cur) != 1 || cur[0] != 0 {
		return fmt.Errorf("after %d operations final array is %v (expected [0])", k, cur)
	}
	return nil
}

func verifyOutput(t testInput, output string) error {
	reader := newTokenReader(output)
	for idx, tc := range t.cases {
		if err := verifyCase(idx, tc, reader); err != nil {
			return fmt.Errorf("test case %d: %v", idx+1, err)
		}
	}
	if reader.hasMore() {
		return fmt.Errorf("extra tokens at the end of output, starting with %q", reader.peek())
	}
	return nil
}

func cloneSlice(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func makeInput(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	stored := make([]testCase, len(cases))
	for i, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for j, v := range cs.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		stored[i] = testCase{n: cs.n, arr: cloneSlice(cs.arr)}
	}
	return testInput{text: sb.String(), cases: stored}
}

func fixedInputs() []testInput {
	return []testInput{
		makeInput([]testCase{
			{n: 4, arr: []int{1, 2, 3, 4}},
		}),
		makeInput([]testCase{
			{n: 4, arr: []int{0, 1, 0, 0}},
			{n: 6, arr: []int{0, 0, 0, 0, 0, 0}},
		}),
		makeInput([]testCase{
			{n: 5, arr: []int{5, 4, 3, 2, 1}},
		}),
		makeInput([]testCase{
			{n: 5000, arr: buildValueRange(5000, 0)},
		}),
	}
}

func buildValueRange(n int, offset int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = (i + offset) % (n + 1)
	}
	return arr
}

func randomInput(rng *rand.Rand) testInput {
	t := rng.Intn(5) + 1
	remaining := 5000
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		minNeeded := (t - i - 1) * 4
		maxLen := remaining - minNeeded
		if maxLen < 4 {
			maxLen = 4
		}
		length := 4 + rng.Intn(maxLen-4+1)
		remaining -= length
		arr := make([]int, length)
		for j := 0; j < length; j++ {
			arr[j] = rng.Intn(length + 1)
		}
		cases = append(cases, testCase{n: length, arr: arr})
	}
	return makeInput(cases)
}

func generateTests() []testInput {
	tests := fixedInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		tests = append(tests, randomInput(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for idx, test := range tests {
		output, err := runBinary(bin, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\nInput:\n%s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		if err := verifyOutput(test, output); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d failed: %v\nInput:\n%sOutput:\n%s\n", idx+1, err, test.text, output)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
