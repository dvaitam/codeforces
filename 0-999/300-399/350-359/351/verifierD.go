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
const testcasesRaw = `4 19 18 5 12 5 4 4 1 4 3 3 2 4 4 4
7 5 8 5 17 13 1 3 2 7 7 3 3
5 16 20 13 14 13 5 4 4 3 3 1 2 4 4 3 5
7 10 14 17 13 19 12 18 5 4 5 3 3 7 7 5 7 6 6
6 18 19 19 4 7 19 3 3 3 1 4 6 6
2 12 3 4 1 1 2 2 2 2 1 1
7 19 11 18 9 17 8 2 3 1 1 1 5 5 5
4 14 10 20 9 2 1 3 3 4
3 13 13 15 5 2 2 3 3 2 2 2 3 2 3
9 11 1 14 19 11 1 13 20 19 2 1 6 8 9
6 20 9 16 1 19 2 1 3 5
8 10 19 20 11 6 12 6 11 3 5 7 7 7 1 3
5 17 8 9 8 11 2 4 4 1 5
6 11 8 15 6 3 11 2 5 6 3 4
2 2 17 2 2 2 2 2
2 20 12 5 1 2 2 2 2 2 2 2 2 2
1 14 2 1 1 1 1
9 8 2 15 17 10 18 11 8 3 5 5 5 4 4 1 9 4 7 1 1
8 4 6 17 10 8 1 17 18 4 1 2 6 6 5 8 1 6
4 7 4 18 4 2 2 3 2 2
8 19 13 2 9 8 9 20 17 5 7 7 8 8 1 1 3 3 2 2
2 16 2 1 2 2
3 11 3 12 4 3 3 3 3 2 3 1 2
7 4 5 18 1 13 3 19 2 1 3 4 7
1 20 4 1 1 1 1 1 1 1 1
4 7 18 9 19 1 4 4
7 5 1 11 12 18 9 4 4 6 6 7 7 6 6 6 7
10 18 4 19 1 16 5 8 13 2 17 1 10 10
7 6 1 11 4 1 4 16 3 5 6 7 7 1 7
10 17 17 8 4 18 4 18 2 18 11 5 3 4 4 5 4 7 10 10 5 7
10 13 12 18 14 3 13 17 8 14 6 4 10 10 3 9 3 5 2 9
8 17 15 19 6 5 9 7 5 5 6 6 5 8 5 6 5 5 5 8
7 7 6 19 12 8 11 16 2 4 7 6 6
8 19 18 1 16 3 13 2 15 2 4 4 4 6
4 7 9 5 6 5 1 3 2 2 3 3 4 4 1 1
2 9 10 1 2 2
10 11 1 1 11 11 14 13 16 3 7 5 8 9 3 8 2 6 2 8 2 9
9 9 4 17 12 12 15 10 9 4 3 9 9 8 9 1 5
10 6 5 6 12 15 4 4 18 5 11 5 7 9 3 10 8 9 3 4 2 4
9 18 19 13 12 4 9 9 13 2 2 1 8 9 9
4 17 12 11 13 4 1 3 4 4 2 3 1 1
10 4 6 7 19 14 13 5 19 20 5 4 4 8 9 9 10 10 4 10
5 12 10 1 15 14 4 3 5 5 5 4 5 4 4
10 7 1 4 8 16 6 17 15 7 7 5 4 4 9 10 2 6 3 5 8 8
10 2 1 12 20 8 17 3 16 18 1 3 6 8 6 7 2 2
2 11 7 1 1 2
4 16 11 4 2 4 1 2 2 3 4 4 1 4
4 16 10 1 15 4 4 4 2 3 1 3 3 4
8 17 12 20 13 8 1 7 9 3 3 6 4 5 4 4
3 19 13 17 2 3 3 1 1
10 6 15 16 6 2 1 13 15 11 14 1 1 4
7 2 13 16 1 8 8 4 4 4 5 2 4 5 5 3 3
10 2 10 9 15 10 16 8 18 9 1 3 6 8 2 2 7 7
10 20 1 4 1 3 1 6 17 2 16 1 4 9
9 11 7 16 11 16 12 2 13 10 5 7 7 5 6 7 7 9 9 9 9
9 13 6 13 18 12 6 12 14 15 2 8 9 6 8
3 17 20 13 4 1 1 1 3 1 2 1 3
2 11 8 5 1 1 2 2 2 2 1 1 1 2
2 11 19 3 1 2 1 1 1 1
3 5 19 1 1 1 2
4 8 18 17 14 5 3 3 4 4 1 1 1 1 1 2
5 3 4 15 13 8 5 1 4 3 4 5 5 1 3 5 5
7 7 4 18 1 15 10 3 3 3 4 4 4 5 7
6 14 3 20 17 7 8 3 1 3 2 5 4 4
5 7 11 6 7 7 5 4 5 3 3 5 5 4 5 4 4
5 1 6 4 1 5 3 5 5 4 4 2 3
5 16 14 2 12 15 2 3 3 1 4
5 14 15 3 7 5 4 3 4 3 3 4 5 4 5
9 18 8 12 10 10 1 15 12 12 3 4 8 1 1 3 8
9 5 18 1 6 2 1 7 15 12 3 9 9 8 8 4 4
5 14 11 2 20 18 1 4 5
5 8 16 14 9 11 1 1 4
1 6 5 1 1 1 1 1 1 1 1 1 1
8 2 15 15 17 20 17 14 12 5 3 5 3 3 3 7 2 5 6 7
8 9 9 15 10 17 5 19 11 2 1 7 8 8
8 19 20 9 1 11 19 20 18 1 8 8
5 9 4 14 3 12 1 5 5
8 7 10 12 6 13 13 11 2 3 4 4 6 7 7 8
1 5 4 1 1 1 1 1 1 1 1
2 1 10 1 2 2
8 15 9 10 18 18 2 6 8 4 3 4 3 8 3 8 8 8
1 5 4 1 1 1 1 1 1 1 1
9 5 18 7 13 4 14 13 6 1 3 2 4 2 4 5 6
7 12 20 3 7 1 12 5 4 2 2 3 7 4 4 6 7
8 1 12 17 15 19 14 15 18 5 5 8 3 7 8 8 4 8 5 6
5 6 11 9 7 5 1 5 5
7 20 6 4 19 1 6 4 4 5 7 3 7 6 7 1 4
9 19 15 11 5 19 20 7 11 15 5 6 8 6 8 6 8 5 7 3 3
10 17 8 11 8 10 14 9 15 5 16 3 9 9 10 10 10 10
2 14 14 5 2 2 1 2 2 2 1 2 1 1
9 16 4 14 13 18 12 1 10 12 5 6 9 8 9 2 9 3 5 4 4
6 3 20 1 5 3 7 3 4 5 2 2 1 5
6 18 15 12 7 15 11 5 1 5 4 4 4 4 3 5 5 5
2 6 20 1 2 2
6 13 1 16 6 16 9 2 4 4 5 6
10 10 14 10 4 3 3 15 11 3 1 3 8 9 2 7 7 8
6 7 13 3 2 17 1 5 5 5 5 5 2 3 4 6 4 5
6 15 17 8 7 19 10 3 5 5 5 6 5 6
6 19 17 7 15 2 14 3 6 6 4 4 2 2
8 4 17 13 7 11 6 2 11 4 4 7 1 4 7 8 7 7`

type query struct{ l, r int }

type testCase struct {
	m int
	b []int
	q []query
}

// Fenwick tree for point updates and prefix sums.
type fenwick struct {
	n    int
	tree []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, tree: make([]int, n+1)}
}

func (f *fenwick) add(i, v int) {
	for x := i; x <= f.n; x += x & -x {
		f.tree[x] += v
	}
}

func (f *fenwick) sum(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += f.tree[x]
	}
	return s
}

func (f *fenwick) query(l, r int) int {
	if l > r {
		return 0
	}
	return f.sum(r) - f.sum(l-1)
}

// solveCase mirrors 351D.go.
func solveCase(tc testCase) []int {
	m := tc.m
	b := tc.b
	qs := make([][]struct{ l, idx int }, m+1)
	for i, qu := range tc.q {
		qs[qu.r] = append(qs[qu.r], struct{ l, idx int }{qu.l, i})
	}
	ans := make([]int, len(tc.q))
	fw := newFenwick(m)
	maxVal := 0
	for i := 1; i <= m; i++ {
		if b[i] > maxVal {
			maxVal = b[i]
		}
	}
	last := make([]int, maxVal+1)
	for i := 1; i <= m; i++ {
		val := b[i]
		if val < len(last) && last[val] != 0 {
			fw.add(last[val], -1)
		}
		fw.add(i, 1)
		if val < len(last) {
			last[val] = i
		}
		for _, qu := range qs[i] {
			ans[qu.idx] = fw.query(qu.l, i)
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		ptr := 0
		m, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %w", idx+1, err)
		}
		ptr++
		if len(fields) < ptr+m+1 {
			return nil, fmt.Errorf("line %d: not enough b values", idx+1)
		}
		b := make([]int, m+1)
		for i := 1; i <= m; i++ {
			val, err := strconv.Atoi(fields[ptr])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b%d: %w", idx+1, i, err)
			}
			b[i] = val
			ptr++
		}
		if ptr >= len(fields) {
			return nil, fmt.Errorf("line %d: missing q", idx+1)
		}
		q, err := strconv.Atoi(fields[ptr])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse q: %w", idx+1, err)
		}
		ptr++
		if len(fields) != ptr+2*q {
			return nil, fmt.Errorf("line %d: expected %d query values got %d", idx+1, 2*q, len(fields)-ptr)
		}
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			l, err := strconv.Atoi(fields[ptr])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse l%d: %w", idx+1, i+1, err)
			}
			r, err := strconv.Atoi(fields[ptr+1])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse r%d: %w", idx+1, i+1, err)
			}
			queries[i] = query{l: l, r: r}
			ptr += 2
		}
		cases = append(cases, testCase{m: m, b: b, q: queries})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.m))
	for i := 1; i <= tc.m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.b[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.q)))
	for _, qu := range tc.q {
		fmt.Fprintf(&sb, "%d %d\n", qu.l, qu.r)
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		expectVals := solveCase(tc)
		expectStrs := make([]string, len(expectVals))
		for i, v := range expectVals {
			expectStrs[i] = strconv.Itoa(v)
		}
		expect := strings.Join(expectStrs, "\n")
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
