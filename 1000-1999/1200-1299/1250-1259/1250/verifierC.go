package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type proj struct {
	r   int
	p   int64
	idx int
}

type projectData struct {
	l int
	r int
	p int64
}

type testcase struct {
	n    int
	k    int64
	data []projectData
}

// Segment tree supporting range add and range max with index.
type SegTree struct {
	n    int
	tree []int64
	idx  []int
	lazy []int64
}

func NewSegTree(n int, k int64) *SegTree {
	st := &SegTree{
		n:    n,
		tree: make([]int64, 4*(n+2)),
		idx:  make([]int, 4*(n+2)),
		lazy: make([]int64, 4*(n+2)),
	}
	var build func(p, l, r int)
	build = func(p, l, r int) {
		if l == r {
			st.tree[p] = -k * int64(l)
			st.idx[p] = l
			return
		}
		mid := (l + r) >> 1
		build(p<<1, l, mid)
		build(p<<1|1, mid+1, r)
		st.pull(p)
	}
	build(1, 1, n)
	return st
}

func (st *SegTree) pull(p int) {
	if st.tree[p<<1] >= st.tree[p<<1|1] {
		st.tree[p] = st.tree[p<<1]
		st.idx[p] = st.idx[p<<1]
	} else {
		st.tree[p] = st.tree[p<<1|1]
		st.idx[p] = st.idx[p<<1|1]
	}
}

func (st *SegTree) apply(p int, val int64) {
	st.tree[p] += val
	st.lazy[p] += val
}

func (st *SegTree) push(p int) {
	if st.lazy[p] != 0 {
		v := st.lazy[p]
		st.apply(p<<1, v)
		st.apply(p<<1|1, v)
		st.lazy[p] = 0
	}
}

func (st *SegTree) update(p, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(p, val)
		return
	}
	st.push(p)
	mid := (l + r) >> 1
	if ql <= mid {
		st.update(p<<1, l, mid, ql, qr, val)
	}
	if qr > mid {
		st.update(p<<1|1, mid+1, r, ql, qr, val)
	}
	st.pull(p)
}

func (st *SegTree) query(p, l, r, ql, qr int) (int64, int) {
	if ql > r || qr < l {
		return int64(-1 << 60), -1
	}
	if ql <= l && r <= qr {
		return st.tree[p], st.idx[p]
	}
	st.push(p)
	mid := (l + r) >> 1
	lv, li := st.query(p<<1, l, mid, ql, qr)
	rv, ri := st.query(p<<1|1, mid+1, r, ql, qr)
	if lv >= rv {
		return lv, li
	}
	return rv, ri
}

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `5 5 12 17 12 10 18 9 5 5 16 7 16 13 11 18 10
4 5 16 18 11 12 18 13 19 20 17 10 20 19
2 7 14 15 9 17 18 17
2 20 17 20 11 9 15 10
4 17 2 8 9 15 20 3 1 7 20 4 12 1
2 6 20 20 12 20 20 1
1 1 20 20 10
2 13 16 18 12 3 16 4
1 6 1 20 18
3 12 7 20 13 7 12 6 6 13 5
5 19 17 20 19 8 8 2 14 18 2 5 7 14 10 15 20
2 2 11 13 11 9 16 16
5 11 9 19 3 15 16 4 17 19 18 13 20 14 6 15 14
2 4 19 19 6 16 18 4
4 9 9 10 12 17 17 1 2 2 5 4 17 13
4 1 11 13 3 13 16 2 7 13 1 2 14 5
5 10 14 17 18 2 16 18 5 9 7 19 19 3 19 20 10
3 9 7 13 1 2 18 3 1 10 5
2 17 5 8 11 1 11 5
3 19 3 5 12 7 11 15 11 18 16
5 17 3 5 11 15 18 13 18 19 18 14 18 3 8 10 13
5 11 5 7 5 5 6 13 19 19 3 17 19 19 10 15 16
4 19 11 20 15 17 19 8 17 17 13 12 17 17
1 8 19 19 16
1 12 14 20 20
2 18 1 17 17 4 5 19
3 2 12 14 11 9 16 16 14 14 8
4 12 6 20 2 16 20 6 13 17 19 16 17 1
4 20 20 20 20 15 17 19 11 17 19 5 18 13
2 12 13 17 8 7 11 8
5 19 12 19 8 20 20 20 5 15 8 5 13 18 6 14 20
1 15 16 17 6
5 12 20 20 3 18 19 7 10 16 9 17 18 12 14 16 1
1 7 7 16 6
5 15 3 11 17 10 17 14 11 17 2 4 5 11 6 18 9
5 18 15 17 18 14 20 16 12 15 3 3 5 17 1 19 17
3 20 4 11 1 2 5 10 3 15 5
5 20 17 18 9 10 12 6 19 19 5 5 12 9 16 18 5
5 9 7 8 5 14 15 7 19 19 12 18 20 16 5 13 3
5 9 20 20 1 5 11 19 1 14 5 5 20 18 2 18 11
2 6 15 17 15 8 18 9
2 4 12 19 9 12 20 19
3 7 10 18 13 3 8 13 4 7 17
2 15 5 18 9 18 18 6
1 9 15 17 13
5 2 12 18 13 18 18 18 19 19 2 10 19 15 17 20 7
4 7 13 14 4 18 19 4 17 20 6 2 8 8
2 4 20 20 14 3 20 7
5 7 7 14 3 7 20 4 17 18 20 7 15 19 17 20 20
4 9 16 20 1 12 17 7 2 8 15 7 15 7
1 5 8 15 16
1 20 9 19 10
4 11 10 20 3 20 20 16 12 13 8 5 18 8
1 16 10 12 17
3 5 4 20 11 11 12 12 7 10 10
5 14 13 19 2 5 19 3 10 20 5 17 18 20 16 18 5
5 5 19 20 9 7 8 6 5 19 18 14 15 2 10 17 19
1 12 4 7 16
2 17 2 2 19 3 17 9
5 9 18 19 6 8 18 4 19 19 16 8 8 8 5 9 5
5 20 6 7 4 9 10 9 20 20 19 9 18 14 2 4 5
4 2 12 17 9 14 17 1 10 20 8 7 19 3
5 5 8 16 7 15 17 1 9 15 19 17 18 11 19 20 7
3 4 15 20 9 8 20 9 10 11 12
2 15 1 13 15 2 11 17
5 14 15 19 19 17 18 6 1 9 4 13 17 7 16 19 9
5 13 5 20 9 6 17 1 7 12 9 4 4 14 10 15 20
2 11 14 20 5 4 11 12
4 15 6 16 12 14 20 14 14 16 19 4 6 18
4 5 4 4 17 13 20 8 13 14 14 2 11 20
2 10 9 17 4 18 20 11
1 5 1 3 11
3 15 7 10 12 19 20 9 18 18 15
4 1 4 4 4 11 14 10 18 18 17 10 11 4
5 20 10 10 8 4 9 15 6 8 15 9 20 6 15 20 19
1 18 13 13 4
1 10 3 15 7
3 13 8 15 20 17 19 6 14 19 1
5 14 20 20 10 12 20 11 16 19 9 12 16 6 6 20 8
4 11 1 12 16 20 20 8 1 4 15 17 18 7
4 20 11 19 5 2 17 14 15 18 12 1 7 19
2 3 16 17 3 17 20 20
3 7 3 16 1 18 19 11 1 7 20
3 14 15 20 17 10 11 2 3 14 13
4 13 8 12 16 13 13 14 10 20 16 5 17 14
4 19 20 20 16 18 19 10 14 20 19 6 11 2
4 15 15 17 3 8 13 8 13 13 20 11 13 6
4 3 4 9 20 2 9 6 1 17 14 14 15 10
4 15 10 18 2 6 13 6 4 7 16 10 14 10
4 2 19 20 11 15 15 2 10 13 15 10 18 7
4 3 14 14 2 3 14 1 10 14 2 1 2 20
1 5 20 20 16
1 19 5 5 18
2 4 3 5 20 4 10 4
4 14 3 6 16 15 19 6 18 18 5 12 18 17
2 3 17 17 18 8 9 13
3 17 11 19 19 18 19 10 8 17 3
2 20 13 18 10 12 20 12
1 1 1 5 9
2 5 6 7 16 16 17 9`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		parts := strings.Fields(ln)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid line %d", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("parse n on line %d: %w", idx+1, err)
		}
		kVal, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse k on line %d: %w", idx+1, err)
		}
		expectedFields := 2 + 3*n
		if len(parts) != expectedFields {
			return nil, fmt.Errorf("line %d: expected %d fields, got %d", idx+1, expectedFields, len(parts))
		}
		data := make([]projectData, 0, n)
		for i := 0; i < n; i++ {
			l, err := strconv.Atoi(parts[2+3*i])
			if err != nil {
				return nil, fmt.Errorf("parse l on line %d: %w", idx+1, err)
			}
			r, err := strconv.Atoi(parts[3+3*i])
			if err != nil {
				return nil, fmt.Errorf("parse r on line %d: %w", idx+1, err)
			}
			p, err := strconv.ParseInt(parts[4+3*i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse p on line %d: %w", idx+1, err)
			}
			data = append(data, projectData{l: l, r: r, p: p})
		}
		cases = append(cases, testcase{n: n, k: kVal, data: data})
	}
	return cases, nil
}

// solve implements the logic from 1250C.go for a single test case.
func solve(tc testcase) string {
	const maxDay = 200000
	events := make([][]proj, maxDay+2)
	for i, d := range tc.data {
		events[d.l] = append(events[d.l], proj{r: d.r, p: d.p, idx: i + 1})
	}

	st := NewSegTree(maxDay, tc.k)

	var bestProfit int64
	bestL, bestR := -1, -1
	for L := maxDay; L >= 1; L-- {
		for _, pr := range events[L] {
			st.update(1, 1, maxDay, pr.r, maxDay, pr.p)
		}
		val, idx := st.query(1, 1, maxDay, L, maxDay)
		profit := val + tc.k*int64(L-1)
		if profit > bestProfit {
			bestProfit = profit
			bestL = L
			bestR = idx
		}
	}

	if bestProfit <= 0 {
		return "0"
	}

	chosen := make([]int, 0)
	for i, d := range tc.data {
		if d.l >= bestL && d.r <= bestR {
			chosen = append(chosen, i+1)
		}
	}

	var out strings.Builder
	fmt.Fprintf(&out, "%d %d %d %d\n", bestProfit, bestL, bestR, len(chosen))
	for i, idx := range chosen {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", idx)
	}
	if len(chosen) > 0 {
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for _, d := range tc.data {
			fmt.Fprintf(&sb, "%d %d %d\n", d.l, d.r, d.p)
		}
		want := solve(tc)
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
