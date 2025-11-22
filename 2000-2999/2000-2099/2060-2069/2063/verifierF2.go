package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type pair struct {
	l, r int
}

type testCase struct {
	n     int
	pairs []pair
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n%s", err, refOut)
		os.Exit(1)
	}
	refVals, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candVals, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range refVals {
		if refVals[i] != candVals[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d token %d: expected %d got %d\n", findTest(tests, i)+1, findTokenIndex(tests, i)+1, refVals[i], candVals[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2063F2-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2063F2")
	cmd := exec.Command("go", "build", "-o", outPath, "2063F2.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
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
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func parseOutputs(out string, tests []testCase) ([]int64, error) {
	tokens := strings.Fields(out)
	total := 0
	for _, tc := range tests {
		total += tc.n + 1
	}
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d tokens, got %d", total, len(tokens))
	}
	res := make([]int64, total)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tok)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	totalN := 0
	const maxSumN = 300000

	add := func(tc testCase) {
		if totalN+tc.n > maxSumN {
			return
		}
		tests = append(tests, tc)
		totalN += tc.n
	}

	// Small deterministic sequences (n >= 2)
	add(makeFromSequence("(())"))
	add(makeFromSequence("()()"))
	add(makeFromSequence("(()())"))
	add(makeFromSequence("((()))"))
	add(makeFromSequence("()(())"))

	for len(tests) < 120 && totalN < maxSumN {
		n := rng.Intn(5000) + 2
		if totalN+n > maxSumN {
			n = maxSumN - totalN
		}
		seq := randomBalanced(rng, n)
		tc := makeFromSequence(seq)
		shufflePairs(rng, tc.pairs)
		add(tc)
	}

	return tests
}

func makeFromSequence(seq string) testCase {
	stack := make([]int, 0)
	pairs := make([]pair, 0, len(seq)/2)
	for i, ch := range seq {
		pos := i + 1
		if ch == '(' {
			stack = append(stack, pos)
		} else {
			open := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pairs = append(pairs, pair{l: open, r: pos})
		}
	}
	return testCase{n: len(seq) / 2, pairs: pairs}
}

func randomBalanced(rng *rand.Rand, n int) string {
	seq := make([]byte, 0, 2*n)
	balance := 0
	for i := 0; i < 2*n; i++ {
		rem := 2*n - i
		if balance == rem {
			seq = append(seq, ')')
			balance--
			continue
		}
		if balance == 0 {
			seq = append(seq, '(')
			balance++
			continue
		}
		if rng.Intn(2) == 0 {
			seq = append(seq, '(')
			balance++
		} else {
			seq = append(seq, ')')
			balance--
		}
	}
	return string(seq)
}

func shufflePairs(rng *rand.Rand, arr []pair) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for _, p := range tc.pairs {
			fmt.Fprintf(&b, "%d %d\n", p.l, p.r)
		}
	}
	return b.String()
}

func findTest(tests []testCase, idx int) int {
	acc := 0
	for t, tc := range tests {
		lenT := tc.n + 1
		if idx < acc+lenT {
			return t
		}
		acc += lenT
	}
	return len(tests) - 1
}

func findTokenIndex(tests []testCase, idx int) int {
	acc := 0
	for _, tc := range tests {
		lenT := tc.n + 1
		if idx < acc+lenT {
			return idx - acc
		}
		acc += lenT
	}
	return 0
}
