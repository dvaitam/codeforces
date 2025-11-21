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

type caseData struct {
	n       int
	parents []int
	a       []int
}

type testCase struct {
	input   string
	caseCnt int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.caseCnt; caseIdx++ {
			if refVals[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, caseIdx+1, refVals[caseIdx], gotVals[caseIdx], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "2164F1_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2164F1.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	vals := make([]int64, expected)
	for i, token := range fields {
		v, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		vals[i] = v
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 40, 50)...)
	tests = append(tests, randomTests(rng, 25, 500)...)
	tests = append(tests, randomTests(rng, 15, 2000)...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	c1 := caseData{
		n:       5,
		parents: []int{1, 2, 3, 4},
		a:       []int{0, 1, 2, 3, 4},
	}
	c2 := caseData{
		n:       5,
		parents: []int{1, 2, 1, 1},
		a:       []int{0, 1, 0, 0, 0},
	}
	c3 := caseData{
		n:       8,
		parents: []int{1, 1, 3, 3, 4, 5, 7},
		a:       []int{0, 0, 1, 0, 1, 3, 3, 1},
	}
	c4 := chainCase(6)
	c5 := reverseChainCase(6)
	return []testCase{
		makeTestCase([]caseData{c1, c2, c3}),
		makeTestCase([]caseData{c4, c5}),
	}
}

func randomTests(rng *rand.Rand, batches int, maxN int) []testCase {
	const limit = 5000
	var tests []testCase
	for b := 0; b < batches; b++ {
		remaining := limit
		targetCases := rng.Intn(4) + 1
		var cases []caseData
		for len(cases) < targetCases && remaining >= 2 {
			limitN := maxN
			if limitN > remaining {
				limitN = remaining
			}
			if limitN < 2 {
				break
			}
			n := rng.Intn(limitN-1) + 2
			cases = append(cases, randomCase(rng, n))
			remaining -= n
		}
		if len(cases) == 0 {
			cases = append(cases, randomCase(rng, 2))
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests(rng *rand.Rand) []testCase {
	largeChain := chainCase(5000)
	star := starCase(4000, rng)
	randomLarge := randomCase(rng, 4500)
	return []testCase{
		makeTestCase([]caseData{largeChain}),
		makeTestCase([]caseData{star}),
		makeTestCase([]caseData{randomLarge}),
	}
}

func makeTestCase(cases []caseData) testCase {
	sumN := 0
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sumN += c.n
		if sumN > 5000 {
			panic("sum of n exceeds 5000 in a single test case bundle")
		}
		sb.WriteString(strconv.Itoa(c.n))
		sb.WriteByte('\n')
		if c.n > 1 {
			for i := 0; i < c.n-1; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(c.parents[i]))
			}
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
		for i := 0; i < c.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c.a[i]))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), caseCnt: len(cases)}
}

func randomCase(rng *rand.Rand, n int) caseData {
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rng.Intn(i-1) + 1
	}
	perm := rng.Perm(n)
	permutation := make([]int, n)
	for i := 0; i < n; i++ {
		permutation[i] = perm[i] + 1
	}
	return buildCaseFromPermutation(n, parents, permutation)
}

func chainCase(n int) caseData {
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = i - 1
	}
	permutation := make([]int, n)
	for i := 0; i < n; i++ {
		permutation[i] = i + 1
	}
	return buildCaseFromPermutation(n, parents, permutation)
}

func reverseChainCase(n int) caseData {
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = i - 1
	}
	permutation := make([]int, n)
	for i := 0; i < n; i++ {
		permutation[i] = n - i
	}
	return buildCaseFromPermutation(n, parents, permutation)
}

func starCase(n int, rng *rand.Rand) caseData {
	if n < 2 {
		n = 2
	}
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = 1
	}
	permutation := make([]int, n)
	if rng == nil {
		for i := 0; i < n; i++ {
			permutation[i] = i + 1
		}
	} else {
		perm := rng.Perm(n)
		for i := 0; i < n; i++ {
			permutation[i] = perm[i] + 1
		}
	}
	return buildCaseFromPermutation(n, parents, permutation)
}

func buildCaseFromPermutation(n int, parents []int, permutation []int) caseData {
	if len(parents) != n-1 {
		panic("invalid parents length")
	}
	if len(permutation) != n {
		panic("invalid permutation length")
	}
	parentArr := make([]int, n+1)
	parentArr[1] = 0
	for i := 2; i <= n; i++ {
		parentArr[i] = parents[i-2]
	}
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = permutation[i-1]
	}
	a := make([]int, n)
	for u := 1; u <= n; u++ {
		cnt := 0
		cur := parentArr[u]
		for cur != 0 {
			if val[cur] < val[u] {
				cnt++
			}
			cur = parentArr[cur]
		}
		a[u-1] = cnt
	}
	parentCopy := append([]int(nil), parents...)
	return caseData{
		n:       n,
		parents: parentCopy,
		a:       a,
	}
}
