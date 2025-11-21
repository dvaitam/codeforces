package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := generateTests()
	for i, tc := range tests {
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
		if err := validate(tc.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(input, output string) error {
	fields := strings.Fields(input)
	ptr := 0
	readInt := func() (int, error) {
		if ptr >= len(fields) {
			return 0, fmt.Errorf("unexpected end of input while parsing")
		}
		var x int
		_, err := fmt.Sscan(fields[ptr], &x)
		ptr++
		return x, err
	}
	t, err := readInt()
	if err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	type test struct {
		n, k int
	}
	tests := make([]test, t)
	for i := 0; i < t; i++ {
		n, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read n for case %d: %v", i+1, err)
		}
		k, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read k for case %d: %v", i+1, err)
		}
		tests[i] = test{n: n, k: k}
	}

	lines := splitLines(output)
	if len(lines) != t {
		return fmt.Errorf("expected %d output lines, got %d", t, len(lines))
	}

	for i, tc := range tests {
		s := strings.TrimSpace(lines[i])
		if len(s) != tc.n {
			return fmt.Errorf("case %d: expected length %d, got %d", i+1, tc.n, len(s))
		}
		ones := strings.Count(s, "1")
		if ones != tc.k {
			return fmt.Errorf("case %d: expected %d ones, got %d", i+1, tc.k, ones)
		}
		for idx, ch := range s {
			if ch != '0' && ch != '1' {
				return fmt.Errorf("case %d: invalid character %q at position %d", i+1, ch, idx+1)
			}
		}
		// Count subsequences
		count101, count010 := countSubseq(s)
		if count101 != count010 {
			return fmt.Errorf("case %d: subseq counts differ (101=%d, 010=%d)", i+1, count101, count010)
		}
	}
	return nil
}

func splitLines(s string) []string {
	raw := strings.Split(s, "\n")
	var res []string
	for _, line := range raw {
		if line == "" {
			// allow trailing empty lines to be ignored
			continue
		}
		res = append(res, line)
	}
	return res
}

func countSubseq(s string) (int, int) {
	// Count 101 subsequences: for each '0', add onesBefore * onesAfter.
	totalOnes := strings.Count(s, "1")
	totalZeros := len(s) - totalOnes
	onesBefore := 0
	zerosBefore := 0
	var cnt101, cnt010 int
	for _, ch := range s {
		if ch == '0' {
			onesAfter := totalOnes - onesBefore
			cnt101 += onesBefore * onesAfter
			zerosBefore++
		} else { // '1'
			zerosAfter := totalZeros - zerosBefore
			cnt010 += zerosBefore * zerosAfter
			onesBefore++
		}
	}
	return cnt101, cnt010
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Sample from statement.
	tests = append(tests, testCase{
		input: "5\n4 2\n5 3\n5 5\n6 2\n1 1\n",
	})

	// Edge cases: all zeros, all ones, small sizes.
	tests = append(tests, buildSingle(1, 0))
	tests = append(tests, buildSingle(1, 1))
	tests = append(tests, buildSingle(2, 0))
	tests = append(tests, buildSingle(2, 2))

	// Random cases.
	for i := 0; i < 50; i++ {
		n := rng.Intn(100) + 1
		k := rng.Intn(n + 1)
		tests = append(tests, buildSingle(n, k))
	}

	return tests
}

func buildSingle(n, k int) testCase {
	return testCase{
		input: fmt.Sprintf("1\n%d %d\n", n, k),
	}
}
