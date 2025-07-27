package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n    int
	cnt  []int
	lazy []bool
}

func newSegTree(s string) *SegTree {
	n := len(s)
	cnt := make([]int, 4*n)
	lazy := make([]bool, 4*n)
	st := &SegTree{n: n, cnt: cnt, lazy: lazy}
	st.build(1, 0, n-1, s)
	return st
}

func (st *SegTree) build(idx, l, r int, s string) {
	if l == r {
		if s[l] == 'R' {
			st.cnt[idx] = 1
		}
		return
	}
	m := (l + r) / 2
	st.build(idx*2, l, m, s)
	st.build(idx*2+1, m+1, r, s)
	st.cnt[idx] = st.cnt[idx*2] + st.cnt[idx*2+1]
}

func (st *SegTree) push(idx, l, r int) {
	if !st.lazy[idx] || l == r {
		return
	}
	m := (l + r) / 2
	st.apply(idx*2, l, m)
	st.apply(idx*2+1, m+1, r)
	st.lazy[idx] = false
}

func (st *SegTree) apply(idx, l, r int) {
	st.cnt[idx] = (r - l + 1) - st.cnt[idx]
	st.lazy[idx] = !st.lazy[idx]
}

func (st *SegTree) flipRange(idx, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		delta := (r - l + 1) - 2*st.cnt[idx]
		st.cnt[idx] = (r - l + 1) - st.cnt[idx]
		st.lazy[idx] = !st.lazy[idx]
		return delta
	}
	st.push(idx, l, r)
	m := (l + r) / 2
	d1 := st.flipRange(idx*2, l, m, ql, qr)
	d2 := st.flipRange(idx*2+1, m+1, r, ql, qr)
	st.cnt[idx] = st.cnt[idx*2] + st.cnt[idx*2+1]
	return d1 + d2
}

func (st *SegTree) flip(l, r int) int {
	return st.flipRange(1, 0, st.n-1, l, r)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}

	var left, right, top, bottom string
	fmt.Fscan(in, &left)
	fmt.Fscan(in, &right)
	fmt.Fscan(in, &top)
	fmt.Fscan(in, &bottom)

	stL := newSegTree(left)
	stR := newSegTree(right)
	stU := newSegTree(top)
	stD := newSegTree(bottom)

	totalRed := stL.cnt[1] + stR.cnt[1] + stU.cnt[1] + stD.cnt[1]
	totalPorts := 2 * (n + m)
	fmt.Fprintln(out, min(totalRed, totalPorts-totalRed))

	for i := 0; i < q; i++ {
		var s string
		var l, r int
		fmt.Fscan(in, &s, &l, &r)
		l--
		r--
		var delta int
		switch s {
		case "L":
			delta = stL.flip(l, r)
		case "R":
			delta = stR.flip(l, r)
		case "U":
			delta = stU.flip(l, r)
		case "D":
			delta = stD.flip(l, r)
		}
		totalRed += delta
		fmt.Fprintln(out, min(totalRed, totalPorts-totalRed))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
