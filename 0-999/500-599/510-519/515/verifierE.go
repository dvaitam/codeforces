package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const negInf = int64(-4e18)

// testcasesE will be filled with the contents of testcasesE.txt.
const testcasesE = `
100
6 3
1 1 1 3 4 4
4 2 3 4 4 5
1 1
4 6
3 6
6 2
4 4 4 5 5 3
4 5 3 3 4 5
6 2
3 4
6 4
3 1 2 3 2 4
1 5 3 2 2 4
1 1
2 5
1 2
6 2
3 4
4 2 3
4 4 3
3 3
1 1
3 3
2 2
3 1
3 2 4
3 5 2
2 2
6 1
1 1 3 1 3 3
1 1 1 5 1 2
5 2
6 1
5 2 1 5 2 1
1 1 5 1 2 1
4 4
3 1
4 4 5
2 5 1
1 1
3 1
5 2 1
4 3 5
2 2
5 2
5 1 2 5 4
3 3 3 5 3
1 1
1 2
5 2
2 2 1 4 4
2 3 3 2 1
2 3
3 3
3 4
2 3 3
2 1 5
3 3
1 1
2 2
3 3
3 3
2 4 1
4 5 3
1 1
3 3
3 3
6 1
3 2 2 3 3 5
1 2 4 2 5 5
5 1
6 3
1 3 1 1 5 2
3 2 4 2 5 5
2 2
6 6
2 4
6 3
3 3 1 1 4 3
2 4 2 4 5 1
6 2
5 6
4 4
6 2
2 4 4 4 2 3
1 2 5 1 3 2
1 2
3 3
6 3
5 4 3 5 4 4
4 1 3 5 3 3
3 6
3 4
5 5
3 4
5 5 2
4 3 4
3 3
1 1
1 1
3 3
5 2
3 4 1 2 1
2 5 4 2 4
3 4
2 4
5 4
3 5 2 1 3
5 3 4 5 5
1 3
2 2
2 2
3 3
3 1
4 4 1
5 5 1
1 1
6 2
5 3 1 3 3 2
4 3 1 3 1 1
5 6
1 3
5 3
2 1 4 4 4
2 1 4 3 4
4 4
4 4
2 4
4 3
1 2 4 3
3 3 4 3
2 2
4 4
2 2
6 3
2 1 3 2 3 3
1 1 4 3 3 3
3 6
6 3
6 2
4 2
4 2 1 5
3 1 1 4
1 1
4 4
4 2
1 1 2 1
4 5 1 3
1 1
2 3
6 2
5 4 3 3 2 4
3 2 2 1 2 5
3 6
6 6
3 3
5 4 3
4 1 2
2 2
2 2
1 1
4 4
4 1 3 2
1 3 2 2
2 3
4 4
4 4
3 3
4 1
3 4 2 1
1 4 5 5
4 4
5 3
2 3 4 2 5
1 4 2 1 1
5 1
2 2
1 3
5 2
1 1 2 2 3
4 5 5 3 4
2 4
1 1
4 1
4 1 4 2
2 1 5 2
3 3
5 4
3 5 4 4 3
3 5 3 4 5
3 4
3 4
3 4
3 4
6 2
2 4 2 2 5 2
1 3 5 5 2 5
2 4
1 4
4 4
2 1 1 5
1 1 1 5
2 2
2 3
3 4
4 4
5 1
2 4 1 2 2
2 3 3 3 2
4 4
3 1
4 1 1
2 3 4
1 1
5 1
4 5 4 3 2
4 1 3 3 2
1 1
5 4
1 5 2 2 5
5 3 4 1 4
3 4
3 4
1 1
2 3
5 3
5 2 4 4 1
5 3 2 4 1
3 3
5 2
1 2
6 3
2 5 3 5 2 1
1 2 1 2 4 2
1 4
3 3
3 5
5 4
5 1 3 4 4
2 4 1 5 5
4 4
2 3
3 3
2 3
4 1
3 4 2 3
3 5 2 4
1 1
4 4
4 1 1 5
2 5 1 2
3 4
1 2
2 2
4 4
4 1
4 1 1 1
5 2 2 4
2 3
3 2
3 5 1
5 3 3
3 3
2 2
6 3
1 3 3 3 1 3
5 2 3 1 1 5
4 1
4 4
5 5
5 1
2 5 4 5 1
1 1 1 2 4
4 5
4 2
5 1 1 1
4 4 3 2
4 1
3 4
4 3
5 2 1 4
4 5 1 2
3 4
2 2
2 2
4 1
3 3 3 3
2 2 5 4
3 4
3 4
2 4 2
4 4 2
3 3
2 2
1 1
3 3
4 2
5 4 4 2
5 1 5 1
2 2
3 3
3 2
5 5 5
2 5 4
2 2
2 2
3 2
4 5 3
4 1 2
3 3
1 1
5 4
3 3 1 3 1
3 3 2 3 2
2 2
2 4
1 1
3 3
5 4
1 5 1 5 4
1 3 2 5 3
1 1
5 5
5 2
2 2
5 3
2 4 3 2 5
5 2 1 4 2
3 4
3 3
1 2
6 3
5 2 1 4 4 4
3 5 3 4 2 3
1 4
6 1
4 5
4 1
1 4 2 3
1 4 5 3
1 1
6 4
5 3 2 2 2 1
1 5 5 2 3 3
4 5
2 4
2 5
5 2
5 2
1 3 2 1 5
3 5 2 1 5
1 3
3 5
3 2
1 4 3
5 3 2
3 3
2 2
4 1
1 4 3 3
3 4 1 3
3 3
5 3
3 4 2 1 4
5 5 3 4 5
1 1
1 3
1 2
3 2
1 5 3
2 3 2
3 3
2 2
5 1
1 4 3 4 2
2 4 1 2 1
2 3
3 2
2 4 3
3 5 3
3 3
3 3
6 1
2 2 2 2 1 2
1 1 5 5 1 5
6 6
6 4
4 4 1 5 5 3
5 2 2 5 3 5
3 6
4 5
5 1
4 5
4 4
3 5 3 1
3 5 3 5
3 4
4 1
3 3
1 1
5 1
1 3 5 1 5
2 5 1 1 1
3 4
3 1
4 2 3
5 3 3
1 1
6 4
4 3 3 1 5 5
5 3 1 3 1 2
4 4
3 3
3 5
6 6
4 2
5 4 4 3
4 1 1 2
1 1
4 4
4 3
4 1 5 5
3 3 3 5
3 3
1 1
2 3
5 3
4 1 4 2 2
3 2 1 1 4
4 4
5 1
3 4
5 2
2 4 1 4 5
1 3 1 2 1
1 3
4 5
5 1
3 4 5 3 3
3 2 1 5 2
4 5
3 4
3 1 5
2 1 5
1 1
3 3
2 2
3 3
6 4
4 2 4 3 2 1
3 3 3 3 2 4
1 1
6 1
2 5
3 5
4 1
1 5 2 5
2 1 2 2
2 3
5 3
5 1 2 5 1
3 3 3 1 5
3 3
1 2
3 5
5 2
2 5 2 5 2
5 1 3 5 5
5 5
1 1
4 3
3 4 3 4
2 4 4 3
1 2
3 3
3 4
3 3
3 3 3
2 4 2
1 1
3 3
3 3
4 1
3 2 1 4
4 5 1 1
4 4
4 4
1 5 2 4
5 4 4 5
4 4
4 1
3 3
2 3
3 4
2 5 4
2 2 3
1 1
2 2
3 3
2 2
4 3
5 5 2 2
4 5 4 3
4 1
4 4
3 3
4 3
4 3 4 5
2 4 2 2
4 4
1 1
2 3
6 1
3 3 2 2 2 3
2 2 3 2 5 5
2 2
3 2
3 5 3
3 2 1
2 2
1 1
3 1
3 4 5
4 3 5
1 1
6 4
1 5 3 5 3 2
1 1 2 4 3 1
3 3
5 5
5 5
5 6
3 1
5 5 2
2 3 5
2 2
3 1
1 3 2
4 5 1
2 2

`

type node struct {
	mxL, mxR, best int64
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func merge(a, b node) node {
	return node{
		mxL:  max(a.mxL, b.mxL),
		mxR:  max(a.mxR, b.mxR),
		best: max(max(a.best, b.best), a.mxL+b.mxR),
	}
}

func build(seg []node, L, R []int64, p, l, r int) {
	if l == r {
		seg[p] = node{mxL: L[l], mxR: R[l], best: negInf}
		return
	}
	m := (l + r) >> 1
	build(seg, L, R, p<<1, l, m)
	build(seg, L, R, p<<1|1, m+1, r)
	seg[p] = merge(seg[p<<1], seg[p<<1|1])
}

func query(seg []node, p, l, r, ql, qr int) node {
	if qr < l || r < ql {
		return node{mxL: negInf, mxR: negInf, best: negInf}
	}
	if ql <= l && r <= qr {
		return seg[p]
	}
	m := (l + r) >> 1
	left := query(seg, p<<1, l, m, ql, qr)
	right := query(seg, p<<1|1, m+1, r, ql, qr)
	return merge(left, right)
}

type solver struct {
	n   int
	N   int
	seg []node
}

func newSolver(n int, d, h []int64) *solver {
	N := 2 * n
	C := make([]int64, N+2)
	for i := 1; i <= N; i++ {
		di := d[(i-1)%n+1]
		C[i] = C[i-1] + di
	}
	Larr := make([]int64, N+2)
	Rarr := make([]int64, N+2)
	for i := 1; i <= N; i++ {
		hi := h[(i-1)%n+1]
		Larr[i] = 2*hi - C[i-1]
		Rarr[i] = 2*hi + C[i-1]
	}
	seg := make([]node, 4*(N+2))
	build(seg, Larr, Rarr, 1, 1, N)
	return &solver{n: n, N: N, seg: seg}
}

func (s *solver) answer(a, b int) int64 {
	var l, r int
	if a <= b {
		l = b + 1
		r = a - 1 + s.n
	} else {
		l = b + 1
		r = a - 1
	}
	res := query(s.seg, 1, 1, s.N, l, r)
	return res.best
}

type testCase struct {
	n       int
	d, h    []int64
	queries [][2]int
}

func loadTests() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(strings.TrimSpace(testcasesE)))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	tests := make([]testCase, 0, t)
	for idx := 0; idx < t; idx++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing n", idx+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d: missing m", idx+1)
		}
		m, _ := strconv.Atoi(scan.Text())
		d := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing d[%d]", idx+1, i)
			}
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			d[i] = val
		}
		h := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing h[%d]", idx+1, i)
			}
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			h[i] = val
		}
		queries := make([][2]int, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing a for query %d", idx+1, i+1)
			}
			a, _ := strconv.Atoi(scan.Text())
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d: missing b for query %d", idx+1, i+1)
			}
			b, _ := strconv.Atoi(scan.Text())
			queries[i] = [2]int{a, b}
		}
		tests = append(tests, testCase{n: n, d: d, h: h, queries: queries})
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", tc.n, len(tc.queries))
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.FormatInt(tc.d[i], 10))
	}
	buf.WriteByte('\n')
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.FormatInt(tc.h[i], 10))
	}
	buf.WriteByte('\n')
	for _, q := range tc.queries {
		fmt.Fprintf(&buf, "%d %d\n", q[0], q[1])
	}
	return buf.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadTests()
	if err != nil {
		fmt.Println("failed to load testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		s := newSolver(tc.n, tc.d, tc.h)
		for qIdx, q := range tc.queries {
			if !outScan.Scan() {
				fmt.Printf("case %d query %d missing output\n", idx+1, qIdx+1)
				os.Exit(1)
			}
			got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
			expect := s.answer(q[0], q[1])
			if got != expect {
				fmt.Printf("case %d query %d failed: expected %d got %d\n", idx+1, qIdx+1, expect, got)
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("case %d extra output detected\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
