package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n    int
	tree []int
	lazy []bool
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1 // arr is 1-indexed
	st := &SegTree{n: n, tree: make([]int, 4*n+5), lazy: make([]bool, 4*n+5)}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(id, l, r int, arr []int) {
	if l == r {
		st.tree[id] = arr[l]
		return
	}
	m := (l + r) >> 1
	st.build(id<<1, l, m, arr)
	st.build(id<<1|1, m+1, r, arr)
	st.tree[id] = st.tree[id<<1] + st.tree[id<<1|1]
}

func (st *SegTree) apply(id, l, r int) {
	st.tree[id] = (r - l + 1) - st.tree[id]
	st.lazy[id] = !st.lazy[id]
}

func (st *SegTree) push(id, l, r int) {
	if st.lazy[id] {
		m := (l + r) >> 1
		st.apply(id<<1, l, m)
		st.apply(id<<1|1, m+1, r)
		st.lazy[id] = false
	}
}

func (st *SegTree) update(id, l, r, ql, qr int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(id, l, r)
		return
	}
	st.push(id, l, r)
	m := (l + r) >> 1
	st.update(id<<1, l, m, ql, qr)
	st.update(id<<1|1, m+1, r, ql, qr)
	st.tree[id] = st.tree[id<<1] + st.tree[id<<1|1]
}

func (st *SegTree) query(id, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return st.tree[id]
	}
	st.push(id, l, r)
	m := (l + r) >> 1
	left := st.query(id<<1, l, m, ql, qr)
	right := st.query(id<<1|1, m+1, r, ql, qr)
	return left + right
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	g := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		g[p] = append(g[p], i)
	}

	t := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	tin := make([]int, n+1)
	tout := make([]int, n+1)
	order := make([]int, n+1)
	timer := 0
	var dfs func(int)
	dfs = func(v int) {
		timer++
		tin[v] = timer
		order[timer] = v
		for _, to := range g[v] {
			dfs(to)
		}
		tout[v] = timer
	}
	dfs(1)

	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[tin[i]] = t[i]
	}

	st := NewSegTree(arr)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var cmd string
		var v int
		fmt.Fscan(reader, &cmd, &v)
		l := tin[v]
		r := tout[v]
		if cmd[0] == 'p' { // pow
			st.update(1, 1, n, l, r)
		} else {
			res := st.query(1, 1, n, l, r)
			fmt.Fprintln(writer, res)
		}
	}
}
