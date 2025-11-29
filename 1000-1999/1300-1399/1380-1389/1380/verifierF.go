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

const testcasesRaw = `100
5 3
80730
2 1
3 7
2 6
2 5
30
1 6
2 2
2 2
1 2
2 2
3 1
033
1 2
5 3
38323
4 4
1 5
4 2
3 3
154
3 9
1 9
3 5
2 3
54
2 5
1 7
2 2
1 3
0
1 6
1 8
1 5
7 5
0702931
7 3
7 7
3 8
3 8
3 7
2 5
54
1 6
1 3
2 8
2 2
2 4
2 3
54
1 1
1 4
2 2
1 1
9
1 0
4 5
5476
2 0
1 7
3 3
2 9
2 6
2 2
65
1 0
2 4
3 4
928
2 7
3 5
2 4
2 7
7 2
1682752
1 7
3 8
6 1
590458
6 4
8 3
45290784
6 4
8 4
6 5
5 3
65275
3 8
2 8
2 3
6 4
416299
5 6
3 9
5 4
6 0
4 2
9792
2 2
1 7
4 2
0215
2 7
2 8
1 4
7
1 6
1 9
1 3
1 0
6 4
461850
5 9
3 1
3 8
5 5
5 3
26698
3 7
2 2
5 6
8 2
29150615
3 5
7 6
7 2
7476629
5 4
6 4
8 3
60547427
1 1
8 3
5 0
3 4
078
3 4
1 7
1 3
2 4
3 3
479
2 8
3 9
3 1
1 2
4
1 8
1 9
5 5
07559
2 0
1 4
5 7
1 8
2 0
7 4
9976763
3 7
7 1
3 0
6 6
5 4
42278
4 5
5 2
4 9
5 0
2 2
41
1 0
2 6
2 4
70
1 1
1 9
1 2
2 7
3 2
926
1 1
3 3
1 5
1
1 1
1 0
1 6
1 0
1 1
6 5
775502
6 4
2 9
6 8
3 8
5 9
4 3
1834
3 8
2 3
3 7
7 2
2051901
7 1
5 9
8 2
67751808
7 0
3 6
4 1
1732
4 5
4 3
5956
4 2
3 4
4 5
1 5
9
1 2
1 1
1 0
1 2
1 3
1 4
7
1 0
1 7
1 9
1 5
8 4
13221567
3 8
5 1
5 2
5 3
1 4
0
1 5
1 0
1 7
1 2
2 3
47
1 2
1 3
1 9
4 1
9654
2 7
6 3
172323
1 9
4 4
1 7
1 4
6
1 0
1 8
1 1
1 3
8 1
70307600
5 6
8 3
00328638
4 1
6 1
2 8
3 5
139
1 7
3 5
2 6
3 6
3 2
1 3
7
1 8
1 2
1 7
2 3
65
1 2
1 7
1 1
7 5
3407303
7 4
3 1
1 6
7 5
4 8
1 3
2
1 7
1 0
1 6
5 1
68251
3 2
3 3
754
1 1
3 8
1 0
2 1
25
2 4
1 2
8
1 2
1 3
6 3
029921
3 8
5 2
6 2
8 4
61557379
7 2
8 0
3 8
8 7
1 1
2
1 0
4 1
0937
3 6
8 2
34348853
6 5
4 0
8 4
52237601
3 6
2 4
4 5
8 7
2 2
60
1 8
1 6
1 3
0
1 8
1 8
1 8
7 4
9808265
5 0
6 8
2 4
3 1
1 1
5
1 4
1 3
3
1 0
1 4
1 9
4 4
8113
3 8
4 6
1 6
4 2
4 5
5651
3 3
1 3
3 4
3 9
1 2
2 4
23
1 9
2 6
1 6
2 8
4 3
6419
3 2
4 0
2 1
5 5
89926
2 5
5 1
3 1
1 7
2 8
2 3
37
2 3
2 1
1 3
3 1
152
2 9
2 5
75
1 8
1 5
1 1
1 3
2 5
3 4
136
3 7
2 0
1 9
2 2
1 4
6
1 1
1 4
1 1
1 5
4 3
6247
1 9
1 4
4 5
2 5
60
1 6
1 4
1 2
2 1
2 6
3 5
631
2 5
1 1
1 2
2 0
1 4
8 4
81033707
2 4
8 8
1 0
3 6
7 5
9903066
7 6
4 0
7 2
3 2
3 3
7 1
5366342
4 5
5 1
12711
4 0
7 2
3075972
4 3
1 5
1 1
2
1 1
5 2
55324
5 8
1 4
2 1
65
2 5
8 5
97392408
1 6
3 3
8 0
1 0
4 6
4 5
9340
2 0
4 5
1 8
3 7
4 1
2 3
98
2 9
2 6
2 2
2 4
03
2 1
2 8
2 6
1 7
6 1
709939
4 0
1 1
1
1 2`

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

func parseTestcases(raw string) ([]testCaseF, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("empty test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count")
	}
	tests := make([]testCaseF, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing n for case %d", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid n for case %d", i+1)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing m for case %d", i+1)
		}
		m, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid m for case %d", i+1)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing string for case %d", i+1)
		}
		s := sc.Text()
		ops := make([][2]int, m)
		for j := 0; j < m; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("missing op %d for case %d", j+1, i+1)
			}
			x, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid op %d for case %d", j+1, i+1)
			}
			if !sc.Scan() {
				return nil, fmt.Errorf("missing op %d value for case %d", j+1, i+1)
			}
			d, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("invalid op %d value for case %d", j+1, i+1)
			}
			ops[j] = [2]int{x, d}
		}
		tests = append(tests, testCaseF{n: n, m: m, s: s, ops: ops})
	}
	return tests, nil
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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		for _, op := range tc.ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(tc)
		lines := strings.Fields(got)
		if len(lines) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d outputs got %d\n", idx+1, len(expected), len(lines))
			os.Exit(1)
		}
		for i := 0; i < len(expected); i++ {
			val, err := strconv.ParseInt(lines[i], 10, 64)
			if err != nil || val != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d line %d expected %d got %s\n", idx+1, i+1, expected[i], lines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
