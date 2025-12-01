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

const (
	refSource   = "./2136B.go"
	maxTests    = 300
	totalNLimit = 5000
)

type testCase struct {
	n int
	k int
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refDecisions, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candDecisions, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		ref := refDecisions[i]
		cand := candDecisions[i]
		if !ref.ok {
			// No solution according to reference; candidate must also say NO.
			if cand.ok {
				fmt.Fprintf(os.Stderr, "test %d: reference says NO but candidate outputs YES\n", i+1)
				os.Exit(1)
			}
			continue
		}
		// Reference says YES; candidate must provide a valid permutation.
		if !cand.ok {
			fmt.Fprintf(os.Stderr, "test %d: reference says YES but candidate outputs NO\n", i+1)
			os.Exit(1)
		}
		if err := validatePermutation(tc, cand.perm); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

type decision struct {
	ok   bool
	perm []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2136B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
		fmt.Fprintf(&b, "%s\n", tc.s)
	}
	return b.String()
}

func parseOutputs(out string, tests []testCase) ([]decision, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([]decision, len(tests))
	for i, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing YES/NO token", i+1)
		}
		token := strings.ToLower(tokens[idx])
		idx++
		if token == "no" {
			res[i] = decision{ok: false}
			continue
		}
		if token != "yes" {
			return nil, fmt.Errorf("test %d: invalid token %q (expected YES/NO)", i+1, token)
		}
		if idx+tc.n > len(tokens) {
			return nil, fmt.Errorf("test %d: not enough permutation numbers", i+1)
		}
		perm := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return nil, fmt.Errorf("test %d: failed to parse permutation value %q: %v", i+1, tokens[idx], err)
			}
			perm[j] = val
			idx++
		}
		res[i] = decision{ok: true, perm: perm}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens after parsing outputs (used %d of %d)", idx, len(tokens))
	}
	return res, nil
}

func validatePermutation(tc testCase, perm []int) error {
	if len(perm) != tc.n {
		return fmt.Errorf("permutation length mismatch: got %d expected %d", len(perm), tc.n)
	}
	seen := make([]bool, tc.n+1)
	for i, v := range perm {
		if v < 1 || v > tc.n {
			return fmt.Errorf("value out of range at position %d: %d", i+1, v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}

	prevGreater := make([]int, tc.n)
	nextGreater := make([]int, tc.n)
	stack := make([]int, 0)
	for i := 0; i < tc.n; i++ {
		for len(stack) > 0 && perm[stack[len(stack)-1]] < perm[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			prevGreater[i] = -1
		} else {
			prevGreater[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := tc.n - 1; i >= 0; i-- {
		for len(stack) > 0 && perm[stack[len(stack)-1]] < perm[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			nextGreater[i] = tc.n
		} else {
			nextGreater[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	for i := 0; i < tc.n; i++ {
		if tc.s[i] != '1' {
			continue
		}
		left := prevGreater[i] + 1
		right := nextGreater[i] - 1
		if right-left+1 >= tc.k {
			return fmt.Errorf("position %d with value %d can be maximum in interval length %d (>=k=%d)", i+1, perm[i], right-left+1, tc.k)
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Deterministic edge cases.
	add(testCase{n: 1, k: 1, s: "0"})
	add(testCase{n: 1, k: 1, s: "1"})
	add(testCase{n: 2, k: 2, s: "11"})
	add(testCase{n: 3, k: 2, s: "101"})
	add(testCase{n: 4, k: 3, s: "0011"})

	for len(tests) < maxTests && sumN < totalNLimit {
		n := rng.Intn(60) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		k := rng.Intn(n) + 1
		bytesS := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(4) == 0 {
				bytesS[i] = '1'
			} else {
				bytesS[i] = '0'
			}
		}
		// ensure at least one '1' occasionally
		if rng.Intn(5) == 0 {
			bytesS[rng.Intn(n)] = '1'
		}
		add(testCase{n: n, k: k, s: string(bytesS)})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, k: 1, s: "0"})
	}
	return tests
}
