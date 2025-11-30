package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesB = `100
1 2
1
1 1
5 5
5 2 5 1 5
2 4 4 5 3
8 5
1 1 6 8 6 7 7 3
3 4 4 1 3
6 3
2 5 5 3 5 6
5 2 4
7 6
7 5 3 3 7 4 2
7 4 6 6 4 6
4 8
3 4 3 4
4 3 4 4 2 3 2 3
8 5
5 7 5 4 8 6 2 6
1 4 2 1 1
5 4
1 5 2 3 2
2 1 4 1
1 6
1
1 1 1 1 1 1
1 1
1
1
5 3
2 2 5 1 4
5 1 2
3 1
1 2 3
3
2 5
2 2
1 2 2 1 2
7 3
4 2 1 6 6 3 7
1 1 4
3 7
2 3 2
1 2 2 2 3 2 3
1 3
1
1 1 1
3 3
1 2 3
1 3 3
1 4
1
1 1 1 1
4 6
3 4 3 1
2 1 4 4 2 1
2 4
1 1
1 1 1 1
4 1
4 4 3 4
2
4 7
4 1 1 4
2 1 4 3 1 1 3
5 6
3 1 4 1 1
3 2 1 4 1 4
8 8
4 2 1 5 1 6 5 2
4 8 4 2 6 7 8 3
6 7
1 3 1 1 1 5
3 6 4 2 6 1 1
8 1
8 5 6 8 3 6 5 8
8
7 8
7 6 3 4 2 2 4
5 3 5 4 6 6 6 1
2 2
2 1
1 2
2 2
1 1
2 2
4 6
4 2 3 1
2 4 1 3 2 3
3 3
2 3 1
2 2 3
3 8
2 3 1
3 2 3 1 2 3 1 2
5 7
3 4 4 3 5
3 1 2 5 5 2 4
7 1
3 4 5 6 4 6 2
7
2 1
2 1
2
4 2
4 2 3 2
4 2
1 7
1
1 1 1 1 1 1 1
8 2
8 6 8 5 8 1 2 6
3 7
5 3
1 2 4 4 4
3 2 1
5 8
1 3 1 5 4
5 4 2 3 4 2 4 5
5 2
3 3 3 3 3
4 5
2 4
2 1
1 2 1 1
8 4
5 1 2 2 7 6 4 6
6 2 6 8
6 3
4 4 3 4 2 6
4 6 2
5 6
2 1 2 4 2
3 2 3 2 2 2
5 7
4 3 3 5 5
5 3 4 3 5 5 1
6 5
4 4 2 3 3 4
4 1 2 3 4
3 1
1 2 1
2
2 7
1 2
1 2 2 2 1 2 2
4 8
1 2 2 3
3 2 4 3 3 1 3 2
4 4
1 3 3 1
2 2 1 4
8 5
3 6 2 6 7 4 1 7
8 8 6 2 8
7 8
2 4 4 5 4 1 1
4 5 2 1 6 5 2 1
7 5
4 7 6 1 3 1 6
3 2 2 1 2
7 2
3 7 6 4 1 7 4
2 1
8 3
1 3 1 1 4 1 6 4
3 6 8
1 3
1
1 1 1
4 1
2 4 3 4
2
2 3
1 1
2 1 2
6 7
2 4 2 4 3 3
1 5 1 4 3 3 5
8 5
4 7 3 2 1 4 4 4
7 1 3 1 5
8 1
4 3 6 1 4 2 3 3
2
8 5
4 3 6 5 2 7 7 1
8 5 2 5 1
4 7
3 3 4 2
4 2 2 4 4 3 4
8 5
4 8 8 4 8 6 5 2
3 6 8 4 3
5 4
5 3 1 1 1
2 3 1 3
5 6
4 1 4 4 1
3 5 5 2 2 2
3 7
3 1 3
3 2 2 3 1 2 2
4 3
3 2 2 3
4 4 2
4 1
3 2 2 4
4
2 8
2 2
2 2 2 1 1 1 1 1
1 1
1
1
7 4
5 3 1 3 5 3 5
7 4 5 1
8 4
5 1 1 8 1 3 3 4
6 4 1 3
5 2
5 5 1 2 4
2 1
5 5
4 1 5 3 3
1 5 3 1 5
6 6
6 3 4 3 5 5
2 1 1 3 4 6
1 6
1
1 1 1 1 1 1
3 1
1 3 3
2
3 5
1 1 1
3 1 2 1 2
2 2
1 1
1 1
1 7
1
1 1 1 1 1 1 1
7 3
4 6 3 1 5 2 5
7 4 4
8 3
8 3 6 3 1 5 3 3
7 5 8
8 8
4 7 7 5 4 6 1 7
1 7 5 1 8 5 5 4
8 8
6 8 4 3 8 5 6 7
2 4 7 2 7 8 8 2
7 8
5 7 6 6 6 3 7
5 3 2 2 2 6 6 2
7 8
4 6 6 7 2 4 3
3 5 4 3 6 6 7 3
6 1
4 5 1 2 4 1
1
1 6
1
1 1 1 1 1 1
7 4
4 1 5 1 3 1 3
4 2 6 6
1 1
1
1
8 1
8 4 4 6 3 6 5 7
8
8 5
2 4 6 4 6 6 3 4
4 4 8 4 7
5 4
5 2 1 4 4
3 2 2 5
3 8
3 1 3
1 1 3 1 1 3 1 2
4 3
3 4 1 1
3 1 3
1 6
1
1 1 1 1 1 1
8 2
5 3 8 2 5 5 1 4
2 7
3 3
3 2 2
1 1 3`

// Embedded solver from 599B.go.
func expected(n int, f []int, b []int) string {
	pos := make([][]int, n+1)
	for i, v := range f {
		pos[v] = append(pos[v], i+1)
	}
	ans := make([]int, len(b))
	ambiguous := false
	for i, v := range b {
		if len(pos[v]) == 0 {
			return "Impossible"
		}
		if len(pos[v]) > 1 {
			ambiguous = true
		}
		ans[i] = pos[v][0]
	}
	if ambiguous {
		return "Ambiguity"
	}
	var sb strings.Builder
	sb.WriteString("Possible\n")
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

type testCase struct {
	n int
	m int
	f []int
	b []int
}

func parseCases() ([]testCase, error) {
	fields := strings.Fields(testcasesB)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	idx := 0
	readInt := func() (int, error) {
		if idx >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[idx])
		idx++
		return v, err
	}
	t, err := readInt()
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err1 := readInt()
		m, err2 := readInt()
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("case %d: bad n or m", i+1)
		}
		f := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad f[%d]", i+1, j+1)
			}
			f[j] = val
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			val, err := readInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: bad b[%d]", i+1, j+1)
			}
			b[j] = val
		}
		cases = append(cases, testCase{n: n, m: m, f: f, b: b})
	}
	return cases, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := expected(tc.n, tc.f, tc.b)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.f {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
