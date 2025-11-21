package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	min  int64
	max  int64
	best int64
}

const INF int64 = 1 << 60

type SegTree struct {
	n          int
	tree       []Node
	crossRight bool
}

func NewSegTree(values []int64, crossRight bool) *SegTree {
	st := &SegTree{
		n:          len(values),
		tree:       make([]Node, len(values)*4),
		crossRight: crossRight,
	}
	st.build(1, 0, st.n-1, values)
	return st
}

func (st *SegTree) build(node, l, r int, vals []int64) {
	if l == r {
		st.tree[node] = Node{min: vals[l], max: vals[l], best: -INF}
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, vals)
	st.build(node<<1|1, mid+1, r, vals)
	st.tree[node] = st.merge(st.tree[node<<1], st.tree[node<<1|1])
}

func (st *SegTree) update(pos int, val int64) {
	st.updateRec(1, 0, st.n-1, pos, val)
}

func (st *SegTree) updateRec(node, l, r, pos int, val int64) {
	if l == r {
		st.tree[node] = Node{min: val, max: val, best: -INF}
		return
	}
	mid := (l + r) >> 1
	if pos <= mid {
		st.updateRec(node<<1, l, mid, pos, val)
	} else {
		st.updateRec(node<<1|1, mid+1, r, pos, val)
	}
	st.tree[node] = st.merge(st.tree[node<<1], st.tree[node<<1|1])
}

func (st *SegTree) merge(left, right Node) Node {
	res := Node{}
	if left.min < right.min {
		res.min = left.min
	} else {
		res.min = right.min
	}
	if left.max > right.max {
		res.max = left.max
	} else {
		res.max = right.max
	}
	if left.best > right.best {
		res.best = left.best
	} else {
		res.best = right.best
	}
	var cross int64
	if st.crossRight {
		cross = right.max - left.min
	} else {
		cross = left.max - right.min
	}
	if cross > res.best {
		res.best = cross
	}
	return res
}

func (st *SegTree) queryBest() int64 {
	if st.n <= 1 {
		return 0
	}
	return st.tree[1].best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		aValues := make([]int64, n)
		bValues := make([]int64, n)
		for i := 0; i < n; i++ {
			aValues[i] = a[i] - int64(i)
			bValues[i] = a[i] + int64(i)
		}
		segA := NewSegTree(aValues, true)
		segB := NewSegTree(bValues, false)

		printAns := func() {
			best := segA.queryBest()
			if segB.queryBest() > best {
				best = segB.queryBest()
			}
			if best < 0 {
				best = 0
			}
			fmt.Fprintln(out, best)
		}

		printAns()
		for ; q > 0; q-- {
			var p int
			var x int64
			fmt.Fscan(in, &p, &x)
			p--
			a[p] = x
			segA.update(p, a[p]-int64(p))
			segB.update(p, a[p]+int64(p))
			printAns()
		}
	}
}
