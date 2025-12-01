package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type caseData struct {
	n int
	s string
}

type caseResult struct {
	res  string
	perm []int
}

type testCase struct {
	desc  string
	input string
	cases []caseData
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		expRes, err := parseResults(exp, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", i+1, tc.desc, err, exp)
			os.Exit(1)
		}

		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		gotRes, err := parseResults(got, tc.cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d (%s): %v\noutput:\n%s\n", i+1, tc.desc, err, got)
			os.Exit(1)
		}

		for idx := range tc.cases {
			if expRes[idx].res != gotRes[idx].res {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected %s but got %s\ninput:\n%s", i+1, tc.desc, idx+1, expRes[idx].res, gotRes[idx].res, tc.input)
				os.Exit(1)
			}
			if expRes[idx].res == "YES" {
				if err := validatePermutation(tc.cases[idx], gotRes[idx].perm); err != nil {
					fmt.Fprintf(os.Stderr, "test %d (%s) case %d: invalid permutation: %v\ninput:\n%soutput:\n%s", i+1, tc.desc, idx+1, err, tc.input, got)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2146C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	src := filepath.Clean("./2146C.go")
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseInputCases(input string) ([]caseData, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &cases[i].n); err != nil {
			return nil, fmt.Errorf("case %d: failed to read n: %v", i+1, err)
		}
		if _, err := fmt.Fscan(reader, &cases[i].s); err != nil {
			return nil, fmt.Errorf("case %d: failed to read s: %v", i+1, err)
		}
		if len(cases[i].s) != cases[i].n {
			return nil, fmt.Errorf("case %d: string length %d does not match n=%d", i+1, len(cases[i].s), cases[i].n)
		}
	}
	return cases, nil
}

func parseResults(output string, cases []caseData) ([]caseResult, error) {
	reader := bufio.NewReader(strings.NewReader(output))
	results := make([]caseResult, len(cases))
	for i, c := range cases {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return nil, fmt.Errorf("case %d: expected YES/NO token: %v", i+1, err)
		}
		token = strings.ToUpper(token)
		if token != "YES" && token != "NO" {
			return nil, fmt.Errorf("case %d: invalid token %q", i+1, token)
		}
		results[i].res = token
		if token == "YES" {
			results[i].perm = make([]int, c.n)
			for j := 0; j < c.n; j++ {
				if _, err := fmt.Fscan(reader, &results[i].perm[j]); err != nil {
					return nil, fmt.Errorf("case %d: failed to read permutation element %d: %v", i+1, j+1, err)
				}
			}
		}
	}
	// ensure no extra tokens except whitespace
	if _, err := fmt.Fscan(reader, new(string)); err == nil {
		return nil, fmt.Errorf("extra output detected")
	}
	return results, nil
}

func validatePermutation(c caseData, perm []int) error {
	n := c.n
	if len(perm) != n {
		return fmt.Errorf("permutation length %d differs from n=%d", len(perm), n)
	}
	pos := make([]int, n+1)
	seen := make([]bool, n+1)
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range", v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
		pos[v] = i
	}

	prefMax := make([]int, n)
	curMax := 0
	for i := 0; i < n; i++ {
		prefMax[i] = curMax
		if perm[i] > curMax {
			curMax = perm[i]
		}
	}
	suffMin := make([]int, n)
	curMin := n + 1
	for i := n - 1; i >= 0; i-- {
		suffMin[i] = curMin
		if perm[i] < curMin {
			curMin = perm[i]
		}
	}

	for val := 1; val <= n; val++ {
		idx := pos[val]
		leftMax := prefMax[idx]
		rightMin := suffMin[idx]
		isStable := leftMax < val && rightMin > val
		expectStable := c.s[val-1] == '1'
		if isStable != expectStable {
			return fmt.Errorf("value %d stability mismatch (expected %v, got %v)", val, expectStable, isStable)
		}
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeTest("all_ones_small", "1\n3\n111\n"),
		makeTest("all_zeros_small", "1\n4\n0000\n"),
		makeTest("mixed_simple", "1\n5\n11011\n"),
		makeTest("zeros_blocks", "1\n8\n11001100\n"),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomTest(fmt.Sprintf("rand-%d", i+1), rng))
	}
	return tests
}

func makeTest(desc, input string) testCase {
	cases, err := parseInputCases(input)
	if err != nil {
		panic(fmt.Sprintf("invalid static test %s: %v", desc, err))
	}
	return testCase{desc: desc, input: input, cases: cases}
}

func randomTest(desc string, rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 2
		var bytes strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(4) == 0 {
				bytes.WriteByte('0')
			} else {
				bytes.WriteByte('1')
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n%s\n", n, bytes.String()))
	}
	cases, err := parseInputCases(sb.String())
	if err != nil {
		panic(fmt.Sprintf("random test generation failed: %v", err))
	}
	return testCase{desc: desc, input: sb.String(), cases: cases}
}
