package main

import (
	"bufio"
	"fmt"
	"os"
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

// Segment tree supporting range add and range max with index

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	maxDay := 200000
	events := make([][]proj, maxDay+2)
	data := make([]projectData, n+1)
	for i := 1; i <= n; i++ {
		var l, r int
		var p int64
		fmt.Fscan(in, &l, &r, &p)
		events[l] = append(events[l], proj{r: r, p: p, idx: i})
		data[i] = projectData{l: l, r: r, p: p}
	}

	st := NewSegTree(maxDay, k)

	bestProfit := int64(0)
	bestL, bestR := -1, -1
	for L := maxDay; L >= 1; L-- {
		for _, pr := range events[L] {
			st.update(1, 1, maxDay, pr.r, maxDay, pr.p)
		}
		val, idx := st.query(1, 1, maxDay, L, maxDay)
		profit := val + k*int64(L-1)
		if profit > bestProfit {
			bestProfit = profit
			bestL = L
			bestR = idx
		}
	}

	if bestProfit <= 0 {
		fmt.Fprintln(out, 0)
		return
	}

	chosen := make([]int, 0)
	for i := 1; i <= n; i++ {
		if data[i].l >= bestL && data[i].r <= bestR {
			chosen = append(chosen, i)
		}
	}

	fmt.Fprintf(out, "%d %d %d %d\n", bestProfit, bestL, bestR, len(chosen))
	for i, idx := range chosen {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, idx)
	}
	if len(chosen) > 0 {
		fmt.Fprintln(out)
	}
}
