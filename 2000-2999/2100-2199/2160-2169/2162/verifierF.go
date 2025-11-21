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

type interval struct {
	l int
	r int
}

type testCase struct {
	n         int
	intervals []interval
}

type testInput struct {
	raw   string
	cases []testCase
}

func buildReference() (string, error) {
	path := "./2162F_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2162F.go")
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInputFromCases(cases []testCase) testInput {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(len(tc.intervals)))
		sb.WriteByte('\n')
		for _, seg := range tc.intervals {
			sb.WriteString(strconv.Itoa(seg.l))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(seg.r))
			sb.WriteByte('\n')
		}
	}
	copied := make([]testCase, len(cases))
	for i := range cases {
		copied[i] = testCase{
			n:         cases[i].n,
			intervals: append([]interval(nil), cases[i].intervals...),
		}
	}
	return testInput{raw: sb.String(), cases: copied}
}

func manualTests() []testInput {
	return []testInput{
		buildInputFromCases([]testCase{
			{n: 3, intervals: []interval{{1, 2}}},
		}),
		buildInputFromCases([]testCase{
			{n: 5, intervals: []interval{{1, 5}, {2, 4}, {3, 3}}},
			{n: 4, intervals: []interval{{1, 4}}},
		}),
		buildInputFromCases([]testCase{
			{n: 3, intervals: []interval{{1, 1}, {2, 2}, {3, 3}}},
		}),
	}
}

func randomInterval(n int, rng *rand.Rand) interval {
	l := rng.Intn(n) + 1
	r := rng.Intn(n-l+1) + l
	return interval{l: l, r: r}
}

func randomCase(n, m int, rng *rand.Rand) testCase {
	segs := make([]interval, m)
	for i := 0; i < m; i++ {
		segs[i] = randomInterval(n, rng)
	}
	return testCase{
		n:         n,
		intervals: segs,
	}
}

func randomTestInput(rng *rand.Rand, maxCases, maxN int) testInput {
	totalN := 0
	totalM := 0
	cases := []testCase{}
	for len(cases) < maxCases && totalN < 3000 && totalM < 3000 {
		remainN := 3000 - totalN
		if remainN < 3 {
			break
		}
		maxAllowedN := maxN
		if maxAllowedN > remainN {
			maxAllowedN = remainN
		}
		if maxAllowedN < 3 {
			maxAllowedN = 3
		}
		n := rng.Intn(maxAllowedN-2) + 3

		remainM := 3000 - totalM
		if remainM <= 0 {
			break
		}
		maxAllowedM := remainM
		if maxAllowedM > n*2 {
			maxAllowedM = n * 2
		}
		if maxAllowedM < 1 {
			maxAllowedM = 1
		}
		m := rng.Intn(maxAllowedM) + 1

		cases = append(cases, randomCase(n, m, rng))
		totalN += n
		totalM += m
	}
	if len(cases) == 0 {
		cases = append(cases, randomCase(3, 1, rng))
	}
	return buildInputFromCases(cases)
}

func buildMaxTest() testInput {
	n := 3000
	m := 3000
	segs := make([]interval, m)
	for i := 0; i < m; i++ {
		l := (i % n) + 1
		r := n
		segs[i] = interval{l: l, r: r}
	}
	return buildInputFromCases([]testCase{{n: n, intervals: segs}})
}

func buildTests() []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTestInput(rng, 3, 20))
	}
	for i := 0; i < 60; i++ {
		tests = append(tests, randomTestInput(rng, 5, 200))
	}
	tests = append(tests, buildMaxTest())
	return tests
}

func parsePermutations(out string, cases []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	perms := make([][]int, len(cases))
	idx := 0
	for i, tc := range cases {
		if idx+tc.n > len(fields) {
			return nil, fmt.Errorf("not enough numbers for test case %d", i+1)
		}
		perm := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q in test case %d", fields[idx], i+1)
			}
			perm[j] = val
			idx++
		}
		perms[i] = perm
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra output values detected (processed %d of %d)", idx, len(fields))
	}
	return perms, nil
}

func mexOfPermutation(tc testCase, perm []int) (int, error) {
	if len(perm) != tc.n {
		return 0, fmt.Errorf("expected %d values, got %d", tc.n, len(perm))
	}
	used := make([]bool, tc.n)
	for pos, v := range perm {
		if v < 0 || v >= tc.n {
			return 0, fmt.Errorf("value %d out of range at position %d", v, pos+1)
		}
		if used[v] {
			return 0, fmt.Errorf("value %d repeated", v)
		}
		used[v] = true
	}

	marks := make([]int, tc.n+2)
	mexValues := make([]int, len(tc.intervals))
	curMark := 1
	for i, seg := range tc.intervals {
		l := seg.l - 1
		r := seg.r - 1
		if l < 0 || r >= tc.n || l > r {
			return 0, fmt.Errorf("invalid interval [%d,%d]", seg.l, seg.r)
		}
		for pos := l; pos <= r; pos++ {
			val := perm[pos]
			marks[val] = curMark
		}
		mex := 0
		for mex <= tc.n && marks[mex] == curMark {
			mex++
		}
		mexValues[i] = mex
		curMark++
	}

	mexSeen := make([]bool, tc.n+2)
	for _, v := range mexValues {
		if v < len(mexSeen) {
			mexSeen[v] = true
		} else {
			return 0, fmt.Errorf("mex value %d out of supported range", v)
		}
	}
	final := 0
	for final < len(mexSeen) && mexSeen[final] {
		final++
	}
	return final, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var refPath string
	fail := func(format string, args ...interface{}) {
		if refPath != "" {
			_ = os.Remove(refPath)
		}
		fmt.Fprintf(os.Stderr, format+"\n", args...)
		os.Exit(1)
	}

	var err error
	refPath, err = buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, test := range tests {
		refOut, err := runProgram(refPath, test.raw)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		candOut, err := runProgram(bin, test.raw)
		if err != nil {
			fail("test %d: runtime error: %v\ninput:\n%s", idx+1, err, test.raw)
		}

		refPerms, err := parsePermutations(refOut, test.cases)
		if err != nil {
			fail("reference output parse error on test %d: %v\noutput:\n%s", idx+1, err, refOut)
		}
		candPerms, err := parsePermutations(candOut, test.cases)
		if err != nil {
			fail("candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
		}

		for caseIdx, tc := range test.cases {
			refMex, err := mexOfPermutation(tc, refPerms[caseIdx])
			if err != nil {
				fail("reference produced invalid permutation on test %d case %d: %v", idx+1, caseIdx+1, err)
			}
			candMex, err := mexOfPermutation(tc, candPerms[caseIdx])
			if err != nil {
				fail("test %d case %d invalid permutation: %v\ninput:\n%s", idx+1, caseIdx+1, err, test.raw)
			}
			if candMex != refMex {
				fail("test %d case %d failed: expected mex(M)=%d, got %d\ninput:\n%s\ncandidate permutation:%v", idx+1, caseIdx+1, refMex, candMex, test.raw, candPerms[caseIdx])
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
