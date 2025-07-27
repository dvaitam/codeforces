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

const MOD int64 = 998244353

type Matrix struct {
	a11, a12, a21, a22 int64
}

func mul(a, b Matrix) Matrix {
	return Matrix{
		a11: (a.a11*b.a11 + a.a12*b.a21) % MOD,
		a12: (a.a11*b.a12 + a.a12*b.a22) % MOD,
		a21: (a.a21*b.a11 + a.a22*b.a21) % MOD,
		a22: (a.a21*b.a12 + a.a22*b.a22) % MOD,
	}
}

func countSum(s int) int64 {
	if s <= 9 {
		return int64(s + 1)
	}
	return int64(19 - s)
}

type SegTree struct {
	n    int
	tree []Matrix
	arr  []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1
	st := &SegTree{n: n, tree: make([]Matrix, 4*n+4), arr: arr}
	st.build(1, 1, n)
	return st
}

func (st *SegTree) build(node, l, r int) {
	if l == r {
		st.tree[node] = st.makeMatrix(l)
		return
	}
	m := (l + r) >> 1
	st.build(node<<1, l, m)
	st.build(node<<1|1, m+1, r)
	st.tree[node] = mul(st.tree[node<<1|1], st.tree[node<<1])
}

func (st *SegTree) makeMatrix(i int) Matrix {
	cnt1 := countSum(st.arr[i])
	cnt2 := int64(0)
	if i > 1 {
		val := st.arr[i-1]*10 + st.arr[i]
		if val >= 10 && val <= 18 {
			cnt2 = countSum(val)
		}
	}
	return Matrix{a11: cnt1 % MOD, a12: cnt2 % MOD, a21: 1, a22: 0}
}

func (st *SegTree) update(pos int) {
	st.updateRec(1, 1, st.n, pos)
}

func (st *SegTree) updateRec(node, l, r, pos int) {
	if l == r {
		st.tree[node] = st.makeMatrix(l)
		return
	}
	m := (l + r) >> 1
	if pos <= m {
		st.updateRec(node<<1, l, m, pos)
	} else {
		st.updateRec(node<<1|1, m+1, r, pos)
	}
	st.tree[node] = mul(st.tree[node<<1|1], st.tree[node<<1])
}

func (st *SegTree) query() Matrix {
	return st.tree[1]
}

type testCaseF struct {
	n, m int
	s    string
	ops  [][2]int
}

func parseTestcases(path string) ([]testCaseF, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseF, T)
	for i := 0; i < T; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var s string
		fmt.Fscan(in, &s)
		ops := make([][2]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &ops[j][0], &ops[j][1])
		}
		cases[i] = testCaseF{n: n, m: m, s: s, ops: ops}
	}
	return cases, nil
}

func solveCase(tc testCaseF) []int64 {
	arr := make([]int, tc.n+1)
	for i := 1; i <= tc.n; i++ {
		arr[i] = int(tc.s[i-1] - '0')
	}
	st := NewSegTree(arr)
	res := make([]int64, tc.m)
	for idx, op := range tc.ops {
		x, d := op[0], op[1]
		arr[x] = d
		st.update(x)
		if x < tc.n {
			st.update(x + 1)
		}
		r := st.query()
		res[idx] = r.a11 % MOD
	}
	return res
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		for _, op := range tc.ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
		}
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected := solveCase(tc)
		scanner := bufio.NewScanner(strings.NewReader(outStr))
		for i := 0; i < len(expected); i++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "case %d: not enough output lines\n", idx+1)
				os.Exit(1)
			}
			got, err := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
			if err != nil || got != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d line %d expected %d got %s\n", idx+1, i+1, expected[i], strings.TrimSpace(scanner.Text()))
				os.Exit(1)
			}
		}
		if scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
