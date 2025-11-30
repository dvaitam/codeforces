package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	a []int64
	m int
	b []int64
}

// solve embeds the logic from 1704G.go (current placeholder) by returning the
// expected output for the provided testcase.
func solve(idx int) string {
	return expectedOutputs[idx]
}

// Embedded copy of testcasesG.txt.
const testcaseData = `
5 5 -1 -5 1 0 5 2 1 -2 -5 -1
3 -3 -3 0 3 4 4 2
5 5 -1 0 0 -1 3 2 1 0
2 4 -2 2 5 5
3 -2 0 -3 3 1 -1 -5
5 0 -3 0 -4 -2 5 -2 1 4 3 3
2 5 5 2 -2 -4
3 1 -4 1 2 2 2
2 -3 -5 2 -5 4
5 -3 2 1 -2 -2 2 0 -2
5 -1 -3 3 2 3 5 -5 -3 -2 -5 -1
5 -4 -3 5 -5 1 4 3 -3 -1 -3
3 -1 3 -2 3 2 4 -2
3 -5 -2 4 3 0 0 5
5 -5 5 1 -4 -3 5 -3 0 -1 -3 -5
5 2 0 0 0 4 3 3 -3 -2
5 5 -4 -2 -1 -1 5 0 5 2 -5 2
5 -3 -5 1 1 -3 4 3 1 -2 4
5 -5 -2 5 2 0 4 2 -1 -4 -5
4 -3 -3 -5 0 3 -2 4 5
2 -5 -5 2 2 2
5 -2 1 -3 -4 -2 4 -1 -2 -5 -4
3 1 5 -1 3 -1 -1 -4
5 -3 -1 0 -1 4 4 -5 -1 -2 -1
5 3 3 -2 -5 0 5 -4 -4 1 0 -1
2 -3 5 2 0 0
5 -3 -5 2 0 -3 3 2 -3 4
5 3 1 -4 -1 -4 5 -3 4 1 2 5
3 -4 -5 3 3 4 4 3
3 -2 -5 3 3 1 5 4
5 -1 4 0 3 0 4 4 -5 -3 1
3 2 5 0 3 -2 -4 4
4 -1 3 -5 -1 3 2 -3 -4
3 -5 4 4 3 -1 -4 3
3 4 5 -3 2 5 1
4 5 -2 4 2 4 -3 -5 -2 -5
2 -4 1 2 1 -2
3 -5 1 -5 2 -5 1
2 0 0 2 1 3
2 0 2 2 2 -2
5 3 -4 1 4 4 5 -5 0 3 -1 -2
5 2 1 0 -5 -5 5 4 1 4 -3 -3
4 -2 -3 2 1 3 -5 0 0
4 -4 1 -3 5 2 -3 -5
4 -1 0 -5 -4 2 -3 -2
4 -2 -3 2 -2 3 0 -3 -1
3 2 -5 4 2 5 -4
2 -3 1 2 3 -2
2 -4 4 2 -1 2
5 2 5 4 4 -4 4 -4 -1 -2 -5
3 -4 -2 -3 3 1 5 -3
3 3 5 -2 3 4 5 -3
2 1 2 2 5 -5
5 5 2 -5 -2 4 4 -3 3 4 2
2 0 -4 2 3 -1
3 3 5 -5 3 -4 -3 3
3 -3 5 -5 2 -5 -4
4 2 -5 5 3 4 5 3 -4 2
5 4 5 -5 -1 0 5 -4 0 -4 -3 -3
2 4 1 2 4 5
3 -3 2 -4 3 3 5 2
2 4 3 2 3 -1
2 -1 -4 2 4 0
5 4 4 -1 2 -1 5 3 2 4 -3 -1
4 -4 5 -1 -5 4 1 -3 2 4
4 5 -1 -4 4 3 1 -2 -2
5 3 4 -4 1 5 3 -1 -5 5
4 -2 0 -4 -1 3 -4 -3 -2
5 0 2 1 -4 3 5 -5 1 -2 1 4
4 4 4 -3 -3 3 -5 -5 5
2 2 -5 2 3 2
4 -3 5 -2 2 4 -2 5 -2 1
4 2 -1 -1 -3 3 5 -4 4
2 3 4 2 4 4
3 -4 2 -4 3 3 3 -3
2 5 0 2 4 -2
2 2 -2 2 -4 2
4 4 3 -1 -5 3 4 5 1
3 1 -3 3 3 -4 2 -1
5 3 4 3 -4 -4 5 -4 -1 1 -5 -4
4 -4 1 5 2 3 5 5 5
3 1 -5 -1 3 5 1 1
5 -5 5 -2 3 -3 2 2 1
3 -4 2 -2 3 5 1 1
2 -1 -5 2 2 1
3 4 -3 -2 3 -1 -4 3
2 1 -3 2 -3 -1
2 5 -1 2 -5 4
4 -1 -5 2 -4 2 5 5
5 5 2 -1 2 -5 5 -5 -1 -5 0 -2
4 -3 -2 -3 -4 2 -4 -3
4 0 2 4 4 4 2 -5 1 4
5 -2 3 -3 3 -4 5 3 1 -2 3 -5
4 1 1 3 5 4 -2 -3 1 -3
2 4 -4 2 -3 0
5 2 3 -3 4 -1 3 -1 -4 3
4 -2 5 4 3 2 -1 -3
3 -5 2 0 2 4 -2
4 -5 1 2 2 4 -5 -3 3 1
5 5 0 -1 -1 -3 4 5 3 0 -4
`

// Expected outputs for each testcase (trimmed).
var expectedOutputs = []string{
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"3\n3 2 1",
	"-1",
	"-1",
	"-1",
	"-1",
	"0",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"2\n3 2",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		pos := 0
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		pos++
		if len(fields) < pos+n+1 {
			return nil, fmt.Errorf("line %d: not enough values for array a", i+1)
		}
		a := make([]int64, n)
		for j := 0; j < n; j++ {
			val, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad a value %d: %v", i+1, j+1, err)
			}
			a[j] = val
		}
		pos += n
		m, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		pos++
		if len(fields) < pos+m {
			return nil, fmt.Errorf("line %d: expected at least %d values for b, got %d", i+1, m, len(fields)-pos)
		}
		b := make([]int64, m)
		for j := 0; j < m; j++ {
			val, err := strconv.ParseInt(fields[pos+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad b value %d: %v", i+1, j+1, err)
			}
			b[j] = val
		}
		tests = append(tests, testCase{n: n, a: a, m: m, b: b})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	fmt.Fprintf(&input, "%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(v, 10))
	}
	input.WriteByte('\n')
	fmt.Fprintf(&input, "%d\n", tc.m)
	for i, v := range tc.b {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.FormatInt(v, 10))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(expected) != got {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(i)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
