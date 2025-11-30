package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `100
4 1 2
14 3 14
12 5 10
9 5 6
8 4 8
19 12 17
17 9 14
4 1 3
17 15 17
15 7 10
20 6 10
10 1 4
13 3 6
19 17 19
19 6 14
16 12 15
14 6 14
8 4 8
19 8 16
11 8 11
19 12 18
17 15 17
14 10 14
18 8 14
8 5 7
18 10 15
19 17 19
12 4 12
19 12 18
5 2 5
3 1 2
4 1 3
10 2 5
11 4 11
9 1 8
4 1 3
14 3 7
3 1 2
5 1 2
3 1 3
7 2 4
19 1 14
4 1 2
4 1 3
6 3 5
18 1 11
17 9 10
11 7 9
18 8 10
13 2 3
17 13 15
19 13 17
19 11 14
13 5 10
16 11 12
20 5 16
4 2 3
7 2 4
6 4 5
19 2 10
10 8 9
11 2 6
14 5 12
11 9 10
7 1 5
16 3 5
19 3 11
6 1 2
8 2 3
9 1 9
17 5 14
15 4 15
9 6 8
16 9 10
4 2 3
6 4 6
3 1 3
12 6 9
3 1 2
6 3 4
3 1 2
16 11 15
17 4 14
5 1 4
3 1 3
5 1 5
9 1 7
15 12 14
7 3 7
6 3 4
6 1 6
13 11 13
9 6 7
3 1 2
18 10 16
17 3 16
14 5 13
19 16 19
16 8 13
15 4 7
`

func solveCase(n, m, k int) string {
	parts := make([]string, 0, n)
	for i := n; i >= k; i-- {
		parts = append(parts, strconv.Itoa(i))
	}
	for i := m + 1; i < k; i++ {
		parts = append(parts, strconv.Itoa(i))
	}
	for i := 1; i <= m; i++ {
		parts = append(parts, strconv.Itoa(i))
	}
	return strings.Join(parts, " ")
}

type testCase struct {
	n int
	m int
	k int
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesC)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		m, err := nextInt()
		if err != nil {
			return nil, err
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n, m: m, k: k}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(lines))
		os.Exit(1)
	}
	for i, tc := range tests {
		gotFields := strings.Fields(lines[i])
		wantFields := strings.Fields(solveCase(tc.n, tc.m, tc.k))
		if len(gotFields) != len(wantFields) {
			fmt.Printf("case %d length mismatch: expected %d numbers got %d\n", i+1, len(wantFields), len(gotFields))
			os.Exit(1)
		}
		for j := range wantFields {
			if gotFields[j] != wantFields[j] {
				fmt.Printf("case %d failed at position %d\nexpected: %s\ngot: %s\n", i+1, j+1, wantFields[j], gotFields[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
