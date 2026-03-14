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
	name  string
	input string
	t     int
}

type parsedCase struct {
	n    int
	s    string
	k    int
	idxs []int // 1-indexed positions
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

// parseInput extracts (n, s) pairs from the test input.
func parseInput(input string) []parsedCase {
	fields := strings.Fields(input)
	t, _ := strconv.Atoi(fields[0])
	cases := make([]parsedCase, t)
	idx := 1
	for i := 0; i < t; i++ {
		n, _ := strconv.Atoi(fields[idx])
		s := fields[idx+1]
		cases[i] = parsedCase{n: n, s: s}
		idx += 2
	}
	return cases
}

// parseOutput parses candidate output for problem 2162/B.
// Each test case produces either:
//   - one line with k, followed by a second line with k space-separated indices (when k > 0)
//   - one line with 0 (no indices line) when k == 0
//   - one line with -1 when no valid answer exists
func parseOutput(out string, expected int) ([]parsedCase, error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	trimmed := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed = append(trimmed, strings.TrimSpace(line))
	}
	for len(trimmed) > 0 && trimmed[len(trimmed)-1] == "" {
		trimmed = trimmed[:len(trimmed)-1]
	}

	result := make([]parsedCase, 0, expected)
	i := 0
	casesRead := 0
	for i < len(trimmed) && casesRead < expected {
		line := trimmed[i]
		if line == "" {
			return nil, fmt.Errorf("unexpected empty line at position %d", i)
		}
		kVal, err := strconv.Atoi(strings.Fields(line)[0])
		if err != nil {
			return nil, fmt.Errorf("invalid integer on line %d: %q", i+1, line)
		}
		i++

		var idxs []int
		if kVal > 0 {
			if i >= len(trimmed) {
				return nil, fmt.Errorf("expected indices line after k=%d for case %d", kVal, casesRead+1)
			}
			idxLine := trimmed[i]
			fields := strings.Fields(idxLine)
			for _, f := range fields {
				val, err := strconv.Atoi(f)
				if err != nil {
					return nil, fmt.Errorf("invalid integer %q in indices line for case %d", f, casesRead+1)
				}
				idxs = append(idxs, val)
			}
			if len(idxs) != kVal {
				return nil, fmt.Errorf("case %d: k=%d but got %d indices", casesRead+1, kVal, len(idxs))
			}
			i++
		} else if kVal == 0 {
			// k == 0: some solutions output an empty line; consume it if present
			if i < len(trimmed) && trimmed[i] == "" {
				i++
			}
		}
		// kVal == -1: no indices line

		result = append(result, parsedCase{k: kVal, idxs: idxs})
		casesRead++
	}
	if casesRead != expected {
		return nil, fmt.Errorf("expected %d cases, got %d", expected, casesRead)
	}
	return result, nil
}

// validateAnswer checks that a candidate answer for one test case is semantically correct.
// It verifies:
//  1. If k == -1, there truly is no valid answer (brute force check for small n).
//  2. If k >= 0, the selected indices form a non-decreasing subsequence
//     and the remaining characters form a palindrome.
func validateAnswer(input parsedCase, answer parsedCase) error {
	n := input.n
	s := input.s
	k := answer.k

	if k == -1 {
		// Verify no solution exists by brute force (n <= 10, so 2^10 = 1024 is fine)
		if hasSolution(s) {
			return fmt.Errorf("candidate says -1 but a valid solution exists for %q", s)
		}
		return nil
	}

	if k < 0 || k > n {
		return fmt.Errorf("invalid k=%d for n=%d", k, n)
	}

	idxs := answer.idxs
	if len(idxs) != k {
		return fmt.Errorf("k=%d but %d indices given", k, len(idxs))
	}

	// Check indices are valid, in range [1, n], and strictly increasing
	removed := make(map[int]bool)
	for i, idx := range idxs {
		if idx < 1 || idx > n {
			return fmt.Errorf("index %d out of range [1, %d]", idx, n)
		}
		if i > 0 && idx <= idxs[i-1] {
			return fmt.Errorf("indices not strictly increasing: %d after %d", idx, idxs[i-1])
		}
		removed[idx] = true
	}

	// Check: subsequence p (the removed chars) must be non-decreasing
	var subseq []byte
	for _, idx := range idxs {
		subseq = append(subseq, s[idx-1]) // 1-indexed to 0-indexed
	}
	for i := 1; i < len(subseq); i++ {
		if subseq[i-1] > subseq[i] {
			return fmt.Errorf("subsequence not non-decreasing: %q", string(subseq))
		}
	}

	// Check: remaining string x must be a palindrome
	var remaining []byte
	for i := 0; i < n; i++ {
		if !removed[i+1] { // 1-indexed
			remaining = append(remaining, s[i])
		}
	}
	for l, r := 0, len(remaining)-1; l < r; l, r = l+1, r-1 {
		if remaining[l] != remaining[r] {
			return fmt.Errorf("remaining string %q is not a palindrome", string(remaining))
		}
	}

	return nil
}

// hasSolution brute-forces whether any valid subsequence exists for a given string.
func hasSolution(s string) bool {
	n := len(s)
	for mask := 0; mask < (1 << n); mask++ {
		var subseq []byte
		var remaining []byte
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				subseq = append(subseq, s[i])
			} else {
				remaining = append(remaining, s[i])
			}
		}
		// Check non-decreasing
		ok := true
		for i := 1; i < len(subseq); i++ {
			if subseq[i-1] > subseq[i] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// Check palindrome
		pal := true
		for l, r := 0, len(remaining)-1; l < r; l, r = l+1, r-1 {
			if remaining[l] != remaining[r] {
				pal = false
				break
			}
		}
		if pal {
			return true
		}
	}
	return false
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_char", input: "1\n1 0\n", t: 1},
		{name: "small_palindrome", input: "1\n4 0110\n", t: 1},
		{name: "multiple_cases", input: "2\n3 010\n4 0100\n", t: 2},
		{name: "example_from_problem", input: "4\n3 010\n5 00011\n5 10100\n6 100101\n", t: 4},
		{name: "all_same", input: "2\n3 000\n3 111\n", t: 2},
		{name: "alternating", input: "2\n4 0101\n4 1010\n", t: 2},
	}
}

func randomTests(count int) []testCase {
	tests := make([]testCase, 0, count)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			n := rng.Intn(7) + 1
			var s strings.Builder
			for k := 0; k < n; k++ {
				s.WriteByte(byte('0' + rng.Intn(2)))
			}
			sb.WriteString(fmt.Sprintf("%d %s\n", n, s.String()))
		}
		tests = append(tests, testCase{
			name:  fmt.Sprintf("random_%d", i+1),
			input: sb.String(),
			t:     t,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}

	tests := append(manualTests(), randomTests(50)...)
	for idx, tc := range tests {
		// Parse the input to get (n, s) for each sub-case
		inputs := parseInput(tc.input)

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		answers, err := parseOutput(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		for ci := 0; ci < tc.t; ci++ {
			inputs[ci].k = 0 // just used for n and s
			if err := validateAnswer(inputs[ci], answers[ci]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s), case %d: %v\ninput: n=%d s=%q\ncandidate output:\n%s",
					idx+1, tc.name, ci+1, err, inputs[ci].n, inputs[ci].s, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
