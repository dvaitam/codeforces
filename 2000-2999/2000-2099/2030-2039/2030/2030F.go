package main

import (
	"bufio"
	"fmt"
	"os"
)

type RangeMaxTree struct {
	tree []int
	lazy []int
	n    int
}

func NewRangeMaxTree(n int) *RangeMaxTree {
	size := 4 * (n + 2)
	return &RangeMaxTree{
		tree: make([]int, size),
		lazy: make([]int, size),
		n:    n,
	}
}

func (st *RangeMaxTree) apply(node, val int) {
	if val > st.tree[node] {
		st.tree[node] = val
	}
	if val > st.lazy[node] {
		st.lazy[node] = val
	}
}

func (st *RangeMaxTree) push(node int) {
	if st.lazy[node] > 0 {
		st.apply(node*2, st.lazy[node])
		st.apply(node*2+1, st.lazy[node])
		st.lazy[node] = 0
	}
}

func (st *RangeMaxTree) update(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(node, val)
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.update(node*2, l, mid, ql, qr, val)
	st.update(node*2+1, mid+1, r, ql, qr, val)
	if st.tree[node*2] > st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *RangeMaxTree) Update(l, r, val int) {
	if l > r {
		return
	}
	if l < 1 {
		l = 1
	}
	if r > st.n {
		r = st.n
	}
	st.update(1, 1, st.n, l, r, val)
}

func (st *RangeMaxTree) query(node, l, r, pos int) int {
	if l == r {
		return st.tree[node]
	}
	st.push(node)
	mid := (l + r) >> 1
	if pos <= mid {
		return st.query(node*2, l, mid, pos)
	}
	return st.query(node*2+1, mid+1, r, pos)
}

func (st *RangeMaxTree) Query(pos int) int {
	if pos < 1 || pos > st.n {
		return 0
	}
	return st.query(1, 1, st.n, pos)
}

type MaxSegTree struct {
	size int
	tree []int
}

func NewMaxSegTree(arr []int) *MaxSegTree {
	n := len(arr) - 1
	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]int, size*2)
	for i := range tree {
		tree[i] = 0
	}
	for i := 1; i <= n; i++ {
		tree[size+i-1] = arr[i]
	}
	for i := size - 1; i >= 1; i-- {
		if tree[2*i] > tree[2*i+1] {
			tree[i] = tree[2*i]
		} else {
			tree[i] = tree[2*i+1]
		}
	}
	return &MaxSegTree{size: size, tree: tree}
}

func (st *MaxSegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.size - 1
	r += st.size - 1
	res := 0
	for l <= r {
		if l&1 == 1 {
			if st.tree[l] > res {
				res = st.tree[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.tree[r] > res {
				res = st.tree[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		last := make(map[int]int)
		prev := make([]int, n+1)
		for i := 1; i <= n; i++ {
			v := a[i]
			prev[i] = last[v]
			last[v] = i
		}

		cross := make([]int, n+1)
		rangeTree := NewRangeMaxTree(n)
		for i := 1; i <= n; i++ {
			if prev[i] > 0 {
				cross[i] = rangeTree.Query(prev[i])
			}
			if prev[i] > 0 && prev[i]+1 <= i-1 {
				rangeTree.Update(prev[i]+1, i-1, prev[i])
			}
		}

		maxTree := NewMaxSegTree(cross)

		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			if maxTree.Query(l, r) >= l {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		}
	}
}
