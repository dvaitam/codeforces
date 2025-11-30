package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesA = `4 4 1 8 16 15 12 3 4 5 2 4 2 3 3 1 3 3 4 1 1
3 4 1 11 13 10 3 3 8 2 3 8 2 2 8 1 1 6 1 4
3 2 2 2 6 18 1 1 2 3 3 1 1 2 2 2
3 3 3 3 17 10 3 3 9 3 3 7 1 3 6 2 2 2 2 1 1
1 3 2 2 1 1 2 1 1 8 1 1 8 1 1 3 3
3 4 4 20 11 2 2 2 7 3 3 3 1 1 4 1 3 3 3 3 3 4 1 1 2 4
2 1 1 0 3 1 1 6 1 1
1 1 1 6 1 1 7 1 1
1 4 1 8 1 1 1 1 1 6 1 1 8 1 1 9 1 4
2 3 2 15 18 1 1 0 1 1 5 2 2 9 2 2 1 2
4 3 3 11 12 8 4 1 4 1 3 3 8 3 3 3 2 3 3 3 3 3
2 3 2 13 20 1 1 9 1 2 2 1 1 10 2 3 3 3
1 4 4 1 1 1 1 1 1 7 1 1 0 1 1 4 4 4 4 4 1 2 1 4
4 3 1 6 0 0 16 1 2 1 2 3 4 2 2 7 2 2
1 3 2 3 1 1 10 1 1 2 1 1 0 1 1 3 3
3 3 3 1 19 20 2 3 10 2 3 8 1 1 6 3 3 1 1 1 2
3 3 2 2 10 19 1 1 4 1 1 9 2 3 6 3 3 2 2
4 2 1 9 5 16 2 3 4 5 3 4 1 1 2
4 3 2 3 15 3 15 4 4 4 3 3 2 4 4 1 1 1 3 3
1 4 1 3 1 1 7 1 1 6 1 1 3 1 1 6 2 3
1 1 1 16 1 1 1 1 1
3 2 1 6 19 4 1 1 7 2 3 8 1 1
4 2 2 20 13 16 15 3 4 7 2 4 9 1 1 2 2
3 1 1 8 19 4 2 3 7 1 1
1 1 1 4 1 1 0 1 1
2 2 2 11 16 2 2 9 1 1 6 1 2 2 2
4 2 1 0 5 9 16 3 4 1 4 4 4 2 2
4 1 1 20 4 7 9 3 3 0 1 1
2 4 1 4 11 2 2 9 2 2 7 1 1 7 1 1 0 2 4
3 1 1 6 12 15 1 1 9 1 1
1 3 1 12 1 1 8 1 1 6 1 1 3 2 3
1 1 1 15 1 1 8 1 1
1 4 3 3 1 1 9 1 1 1 1 1 1 1 1 0 4 4 4 4 4 4
3 1 1 2 3 11 3 3 5 1 1
1 2 2 2 1 1 0 1 1 0 2 2 1 1
2 2 2 3 15 2 2 2 1 1 5 2 2 2 2
3 2 1 3 17 18 2 2 6 1 1 3 2 2
2 2 2 11 13 1 1 9 1 2 1 1 1 2 2
4 2 2 9 20 11 2 2 3 10 3 3 6 2 2 2 2
4 2 2 9 15 2 5 1 3 1 2 4 7 2 2 2 2
2 2 2 10 16 1 2 7 1 2 3 2 2 2 2
1 2 2 3 1 1 6 1 1 10 2 2 2 2
1 1 1 1 1 1 6 1 1
1 3 2 9 1 1 0 1 1 4 1 1 8 1 1 1 2
3 3 3 4 18 16 3 3 8 1 2 10 3 3 4 2 3 2 2 1 1
1 4 4 18 1 1 8 1 1 10 1 1 3 1 1 8 3 4 1 4 3 3 4 4
2 1 1 15 12 1 1 10 1 1
1 3 1 3 1 1 10 1 1 2 1 1 1 2 3
3 4 2 10 13 20 3 3 2 2 2 8 2 2 3 1 2 5 4 4 4 4
2 2 2 6 18 1 2 0 1 1 2 2 2 1 1
3 1 1 11 13 14 1 3 8 1 1
4 3 3 15 6 10 8 1 1 0 2 3 0 3 3 2 1 2 3 3 3 3
2 4 2 10 19 1 1 5 2 2 5 2 2 8 1 1 5 1 2 3 3
2 3 3 9 9 2 2 7 2 2 0 2 2 0 2 3 3 3 1 2
4 4 4 3 2 2 7 1 2 6 2 3 9 1 4 8 4 4 2 2 3 2 2 3 4 3 4
1 3 3 17 1 1 7 1 1 4 1 1 1 3 3 1 3 2 2
2 3 3 0 17 2 2 1 2 2 1 2 2 10 3 3 1 3 1 1
3 3 1 11 20 15 2 2 2 1 3 8 2 3 9 3 3
2 4 2 5 16 1 1 3 2 2 3 1 1 6 2 2 10 4 4 1 3
4 4 3 0 7 17 20 2 4 7 4 4 1 3 3 9 4 4 5 3 4 1 2 1 3
2 3 2 7 8 2 2 4 1 2 9 1 1 5 1 2 3 3
1 2 1 0 1 1 1 1 1 10 2 2
2 4 4 4 11 2 2 10 2 2 6 1 2 4 2 2 9 1 4 4 4 1 2 2 2
3 2 1 7 0 5 3 3 1 2 2 9 2 2
1 3 2 12 1 1 7 1 1 4 1 1 3 3 3 1 3
4 3 2 15 12 0 9 3 4 0 3 3 7 4 4 6 2 3 1 1
3 1 1 18 9 6 3 3 6 1 1
1 3 1 18 1 1 2 1 1 9 1 1 5 2 3
3 3 2 12 19 13 1 1 2 3 3 7 1 2 0 3 3 3 3
2 1 1 8 5 1 1 10 1 1
4 4 3 9 9 0 13 3 4 8 3 4 3 4 4 0 2 4 9 4 4 4 4 4 4
1 3 3 5 1 1 10 1 1 3 1 1 2 2 2 3 3 1 3
4 1 1 5 14 2 13 4 4 4 1 1
3 1 1 0 15 0 2 2 6 1 1
4 1 1 11 18 4 18 3 4 0 1 1
2 1 1 18 11 2 2 1 1 1
4 1 1 0 10 12 4 3 4 10 1 1
4 3 3 1 5 4 4 4 4 8 1 4 2 3 3 1 3 3 2 2 2 3
3 3 3 14 4 14 3 3 0 3 3 10 3 3 5 1 1 3 3 3 3
4 2 2 20 4 15 17 1 1 8 3 3 1 1 2 2 2
3 4 3 20 3 4 2 2 2 3 3 0 2 2 0 2 2 4 4 4 2 3 1 4
1 4 3 19 1 1 9 1 1 2 1 1 4 1 1 1 3 3 1 1 3 3
2 4 1 10 14 2 2 1 1 2 4 1 1 1 2 2 1 4 4
2 3 1 0 12 2 2 8 1 1 2 1 1 9 2 2
4 4 1 7 7 20 9 3 3 3 3 3 2 4 4 5 2 3 9 1 2
1 4 2 4 1 1 5 1 1 0 1 1 2 1 1 6 4 4 2 2
3 2 2 2 4 13 3 3 1 2 3 3 2 2 1 2
2 4 4 18 2 2 2 8 2 2 9 2 2 4 1 2 6 1 1 2 4 4 4 4 4
4 1 1 19 14 5 5 3 4 6 1 1
1 3 2 4 1 1 10 1 1 4 1 1 4 2 2 2 3
1 4 3 12 1 1 2 1 1 9 1 1 0 1 1 7 1 1 4 4 2 2
2 2 1 0 16 1 2 5 2 2 10 1 2
3 2 2 20 16 10 3 3 6 3 3 3 2 2 2 2
3 4 4 3 4 4 1 3 9 3 3 10 3 3 9 3 3 4 2 3 1 1 1 1 3 4
3 1 1 11 18 8 3 3 3 1 1
2 1 1 16 13 1 2 10 1 1
1 2 2 5 1 1 6 1 1 5 1 2 1 2
1 3 3 15 1 1 4 1 1 8 1 1 4 3 3 2 2 1 2
4 4 1 8 15 7 3 3 3 5 1 2 10 1 2 9 2 2 0 4 4
1 4 3 10 1 1 3 1 1 4 1 1 0 1 1 8 3 3 1 3 4 4`

// Embedded solver from 295A.go.
func solveCase(n, m, k int, a []int64, l, r []int, d []int64, queries [][2]int) []int64 {
	opCount := make([]int64, m+3)
	for _, q := range queries {
		x, y := q[0], q[1]
		if x < 1 {
			x = 1
		}
		if y > m {
			y = m
		}
		opCount[x]++
		if y+1 < len(opCount) {
			opCount[y+1]--
		}
	}
	for i := 1; i <= m; i++ {
		opCount[i] += opCount[i-1]
	}
	delta := make([]int64, n+3)
	for i := 1; i <= m; i++ {
		if opCount[i] == 0 {
			continue
		}
		add := opCount[i] * d[i]
		delta[l[i]] += add
		delta[r[i]+1] -= add
	}
	var curr int64
	res := make([]int64, n)
	for i := 1; i <= n; i++ {
		curr += delta[i]
		a[i] += curr
		res[i-1] = a[i]
	}
	return res
}

type testCase struct {
	n, m, k int
	a       []int64
	l, r    []int
	d       []int64
	q       [][2]int
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesA)
	if data == "" {
		return nil, fmt.Errorf("no testcases provided")
	}
	lines := strings.Split(data, "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d malformed", i+1)
		}
		pos := 0
		readInt := func() (int, error) {
			if pos >= len(fields) {
				return 0, fmt.Errorf("line %d: not enough fields", i+1)
			}
			v, err := strconv.Atoi(fields[pos])
			pos++
			return v, err
		}
		n, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", i+1)
		}
		m, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m", i+1)
		}
		k, err := readInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k", i+1)
		}
		if len(fields) != 3+n+3*m+2*k {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", i+1, 3+n+3*m+2*k, len(fields))
		}
		a := make([]int64, n+1)
		for j := 1; j <= n; j++ {
			v, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad a[%d]", i+1, j)
			}
			a[j] = v
			pos++
		}
		l := make([]int, m+1)
		r := make([]int, m+1)
		d := make([]int64, m+1)
		for j := 1; j <= m; j++ {
			v1, err1 := strconv.Atoi(fields[pos])
			v2, err2 := strconv.Atoi(fields[pos+1])
			v3, err3 := strconv.ParseInt(fields[pos+2], 10, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("line %d: bad op %d", i+1, j)
			}
			l[j], r[j], d[j] = v1, v2, v3
			pos += 3
		}
		q := make([][2]int, k)
		for j := 0; j < k; j++ {
			x, err1 := strconv.Atoi(fields[pos])
			y, err2 := strconv.Atoi(fields[pos+1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("line %d: bad query %d", i+1, j+1)
			}
			q[j] = [2]int{x, y}
			pos += 2
		}
		cases = append(cases, testCase{n: n, m: m, k: k, a: a, l: l, r: r, d: d, q: q})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(tc.a[i], 10))
	}
	sb.WriteByte('\n')
	for i := 1; i <= tc.m; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.l[i], tc.r[i], tc.d[i])
	}
	for i := 0; i < tc.k; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.q[i][0], tc.q[i][1])
	}
	return sb.String()
}

func expectedOutput(tc testCase) string {
	res := solveCase(tc.n, tc.m, tc.k, append([]int64{}, tc.a...), tc.l, tc.r, tc.d, tc.q)
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
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
		fmt.Fprintln(os.Stderr, "usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expect := expectedOutput(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
