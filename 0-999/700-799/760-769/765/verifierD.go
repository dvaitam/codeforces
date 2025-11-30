package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesD.txt so the verifier is self-contained.
const testcasesD = `15
2
2 1
2
1 1
4
3 2 4 1
5
1 5 5 3 4
4
3 1 2 2
6
4 6 4 5 3 6
1
1
2
2 1
3
2 2 1
2
1 2
6
6 2 3 4 6 6
6
6 2 3 3 6 2
3
1 2 3
5
1 2 3 1 4
6
1 1 2 1 1 2`

// Embedded solver from 765D.go.
func solveCase(n int, f []int) string {
	for i := 1; i <= n; i++ {
		if f[f[i]] != f[i] {
			return "-1"
		}
	}
	p := make(map[int]int)
	g := make([]int, n+1)
	h := make([]int, 0, n)
	s := 0
	for i := 1; i <= n; i++ {
		fi := f[i]
		idx, ok := p[fi]
		if !ok {
			s++
			idx = s
			p[fi] = s
			h = append(h, fi)
		}
		g[i] = idx
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", s)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(g[i]))
	}
	b.WriteByte('\n')
	for i := 0; i < s; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(h[i]))
	}
	return strings.TrimSpace(b.String())
}

type testCase struct {
	n int
	f []int
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesD), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	lineIdx := 0
	// read t
	for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
		lineIdx++
	}
	if lineIdx >= len(lines) {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[lineIdx]))
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	lineIdx++
	cases := make([]testCase, 0, t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
			lineIdx++
		}
		if lineIdx >= len(lines) {
			return nil, fmt.Errorf("case %d: missing n", caseNum)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[lineIdx]))
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n", caseNum)
		}
		lineIdx++
		for lineIdx < len(lines) && strings.TrimSpace(lines[lineIdx]) == "" {
			lineIdx++
		}
		if lineIdx >= len(lines) {
			return nil, fmt.Errorf("case %d: missing permutation", caseNum)
		}
		numFields := strings.Fields(lines[lineIdx])
		if len(numFields) != n {
			return nil, fmt.Errorf("case %d: expected %d numbers got %d", caseNum, n, len(numFields))
		}
		f := make([]int, n+1)
		for i, s := range numFields {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("case %d: bad number", caseNum)
			}
			f[i+1] = v
		}
		cases = append(cases, testCase{n: n, f: f})
		lineIdx++
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solveCase(tc.n, tc.f)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.f[i]))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
