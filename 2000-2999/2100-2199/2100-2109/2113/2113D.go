package main

import (
	"bufio"
	"fmt"
	"os"
)

type segTree struct {
	n    int
	tree []int
}

func newSegTree(arr []int) *segTree {
	n := len(arr)
	tree := make([]int, 4*n)
	st := &segTree{n: n, tree: tree}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segTree) build(node, l, r int, arr []int) {
	if l == r {
		st.tree[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, arr)
	st.build(node*2+1, mid+1, r, arr)
	if st.tree[node*2] >= st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *segTree) update(node, l, r, idx, val int) {
	if l == r {
		st.tree[node] = val
		return
	}
	mid := (l + r) / 2
	if idx <= mid {
		st.update(node*2, l, mid, idx, val)
	} else {
		st.update(node*2+1, mid+1, r, idx, val)
	}
	if st.tree[node*2] >= st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *segTree) queryFirstGE(node, l, r, start, need int) int {
	if st.tree[node] < need || r < start {
		return -1
	}
	if l == r {
		return l
	}
	mid := (l + r) / 2
	if start <= mid {
		res := st.queryFirstGE(node*2, l, mid, start, need)
		if res != -1 {
			return res
		}
	}
	return st.queryFirstGE(node*2+1, mid+1, r, start, need)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		st := newSegTree(a)
		score := 0

		for i := 0; i < n; i++ {
			idx := st.queryFirstGE(1, 0, n-1, 0, b[i])
			if idx == -1 {
				continue
			}
			st.update(1, 0, n-1, idx, 0)
			score++
		}

		fmt.Fprintln(out, score)
	}
}
