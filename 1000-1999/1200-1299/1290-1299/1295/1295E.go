package main

import (
	"bufio"
	"fmt"
	"os"
)

// Segment tree for range add and range minimum query on int64 values

type SegTree struct {
	n    int
	tree []int64
	lazy []int64
}

func NewSegTree(a []int64) *SegTree {
	n := len(a) - 1 // a is 1-indexed
	size := 1
	for size < n {
		size <<= 1
	}
	st := &SegTree{n: n, tree: make([]int64, size*2), lazy: make([]int64, size*2)}
	var build func(int, int, int)
	build = func(id, l, r int) {
		if l == r {
			if l <= n {
				st.tree[id] = a[l]
			}
			return
		}
		m := (l + r) / 2
		build(id*2, l, m)
		build(id*2+1, m+1, r)
		if st.tree[id*2] < st.tree[id*2+1] {
			st.tree[id] = st.tree[id*2]
		} else {
			st.tree[id] = st.tree[id*2+1]
		}
	}
	build(1, 1, size)
	return st
}

func (st *SegTree) apply(id int, val int64) {
	st.tree[id] += val
	st.lazy[id] += val
}

func (st *SegTree) push(id int) {
	if st.lazy[id] != 0 {
		v := st.lazy[id]
		st.apply(id*2, v)
		st.apply(id*2+1, v)
		st.lazy[id] = 0
	}
}

func (st *SegTree) update(id, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(id, val)
		return
	}
	st.push(id)
	m := (l + r) / 2
	st.update(id*2, l, m, ql, qr, val)
	st.update(id*2+1, m+1, r, ql, qr, val)
	if st.tree[id*2] < st.tree[id*2+1] {
		st.tree[id] = st.tree[id*2]
	} else {
		st.tree[id] = st.tree[id*2+1]
	}
}

func (st *SegTree) Update(l, r int, val int64) {
	if l > r {
		return
	}
	st.update(1, 1, st.size(), l, r, val)
}

func (st *SegTree) size() int {
	return len(st.tree) / 2
}

func (st *SegTree) query(id, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return 1 << 60
	}
	if ql <= l && r <= qr {
		return st.tree[id]
	}
	st.push(id)
	m := (l + r) / 2
	left := st.query(id*2, l, m, ql, qr)
	right := st.query(id*2+1, m+1, r, ql, qr)
	if left < right {
		return left
	}
	return right
}

func (st *SegTree) Query(l, r int) int64 {
	return st.query(1, 1, st.size(), l, r)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
		pos[p[i]] = i
	}
	w := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &w[i])
	}

	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + w[i]
	}

	arr := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = pref[i]
	}

	st := NewSegTree(arr)
	ans := st.Query(1, n-1)

	total := int64(0)
	for x := 1; x <= n; x++ {
		idx := pos[x]
		val := w[idx]
		total += val
		st.Update(idx, n, -2*val)
		cur := st.Query(1, n-1) + total
		if cur < ans {
			ans = cur
		}
	}
	fmt.Fprintln(out, ans)
}
