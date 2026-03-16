package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	name string
	a    []int64
	s    string
}

func formatTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		n := len(tc.a)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			name: "sample1",
			a:    []int64{3, 5, 1, 4, 3, 2},
			s:    "LRLLLR",
		},
		{
			name: "sample2",
			a:    []int64{2, 8},
			s:    "LR",
		},
		{
			name: "sample3",
			a:    []int64{3, 9},
			s:    "RL",
		},
		{
			name: "sample4",
			a:    []int64{1, 2, 3, 4, 5},
			s:    "LRLRR",
		},
		{
			name: "allL_then_R",
			a:    []int64{5, 4, 3, 2, 1},
			s:    "LLLRR",
		},
		{
			name: "allR_then_L",
			a:    []int64{5, 4, 3, 2, 1},
			s:    "RRRLL",
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for i := 0; i < 80; i++ {
		n := rng.Intn(50) + 2
		if totalN+n > 200000 {
			break
		}
		totalN += n
		a := make([]int64, n)
		builder := strings.Builder{}
		for j := 0; j < n; j++ {
			a[j] = int64(rng.Intn(100000) + 1)
			if rng.Intn(2) == 0 {
				builder.WriteByte('L')
			} else {
				builder.WriteByte('R')
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			a:    a,
			s:    builder.String(),
		})
	}
	return tests
}

// solveCase computes the correct answer for one test case.
// The problem: pair up L positions (from left) with R positions (from right)
// greedily, and sum the subarray between each pair.
func solveCase(a []int64, s string) int64 {
	n := len(a)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}

	var ans int64
	l := 0
	r := n - 1
	for l < r {
		if s[l] == 'L' && s[r] == 'R' {
			ans += prefix[r+1] - prefix[l]
			l++
			r--
		} else if s[l] != 'L' {
			l++
		} else {
			r--
		}
	}
	return ans
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, t int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	results := make([]int64, t)
	for i := 0; i < t; i++ {
		var val int64
		if _, err := fmt.Sscan(tokens[i], &val); err != nil {
			return nil, fmt.Errorf("failed to parse integer %q: %v", tokens[i], err)
		}
		results[i] = val
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra tokens after reading %d answers", t)
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatTests(tests)

	// Compute expected answers using built-in solver
	expected := make([]int64, len(tests))
	for i, tc := range tests {
		expected[i] = solveCase(tc.a, tc.s)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if got[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "test %s (%d) failed: expected %d, got %d\ninput:\n%s", tc.name, i+1, expected[i], got[i], formatTests([]testCase{tc}))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
