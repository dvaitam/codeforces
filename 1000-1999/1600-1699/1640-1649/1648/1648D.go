package main

import (
	"bufio"
	"fmt"
	"os"
)

type segTree struct {
	n    int
	tree []int64
	lazy []int64
}

func newSegTree(arr []int64) *segTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &segTree{n: n, tree: make([]int64, 2*n), lazy: make([]int64, 2*n)}
	for i := 0; i < len(arr); i++ {
		st.tree[n+i] = arr[i]
	}
	for i := n - 1; i > 0; i-- {
		if st.tree[2*i] > st.tree[2*i+1] {
			st.tree[i] = st.tree[2*i]
		} else {
			st.tree[i] = st.tree[2*i+1]
		}
	}
	return st
}

func (st *segTree) apply(node int, val int64) {
	st.tree[node] += val
	st.lazy[node] += val
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		st.apply(node<<1, st.lazy[node])
		st.apply(node<<1|1, st.lazy[node])
		st.lazy[node] = 0
	}
}

func (st *segTree) rangeAdd(node, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(node, val)
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.rangeAdd(node<<1, l, mid, ql, qr, val)
	st.rangeAdd(node<<1|1, mid+1, r, ql, qr, val)
	if st.tree[node<<1] > st.tree[node<<1|1] {
		st.tree[node] = st.tree[node<<1]
	} else {
		st.tree[node] = st.tree[node<<1|1]
	}
}

func (st *segTree) RangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	st.rangeAdd(1, 0, st.n-1, l, r, val)
}

func (st *segTree) query(node, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return -1 << 60
	}
	if ql <= l && r <= qr {
		return st.tree[node]
	}
	st.push(node)
	mid := (l + r) >> 1
	left := st.query(node<<1, l, mid, ql, qr)
	right := st.query(node<<1|1, mid+1, r, ql, qr)
	if left > right {
		return left
	}
	return right
}

func (st *segTree) Query(l, r int) int64 {
	if l > r {
		return -1 << 60
	}
	return st.query(1, 0, st.n-1, l, r)
}

func (st *segTree) pointUpdate(idx int, val int64) {
	idx += st.n
	if st.tree[idx] < val {
		st.tree[idx] = val
		idx >>= 1
		for idx > 0 {
			if st.tree[idx<<1] > st.tree[idx<<1|1] {
				st.tree[idx] = st.tree[idx<<1]
			} else {
				st.tree[idx] = st.tree[idx<<1|1]
			}
			idx >>= 1
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	a := make([][]int64, 3)
	for i := 0; i < 3; i++ {
		a[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}

	type offer struct {
		l, r int
		k    int64
	}
	offers := make([][]offer, n+1)
	for i := 0; i < q; i++ {
		var l, r int
		var k int64
		fmt.Fscan(in, &l, &r, &k)
		if r <= n {
			offers[r] = append(offers[r], offer{l - 1, r - 1, k})
		}
	}

	p1 := make([]int64, n+1)
	p2 := make([]int64, n+1)
	p3 := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		p1[i] = p1[i-1] + a[0][i-1]
		p2[i] = p2[i-1] + a[1][i-1]
		p3[i] = p3[i-1] + a[2][i-1]
	}
	left := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		left[i] = p1[i] - p2[i-1]
	}
	suf := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		suf[i] = suf[i+1] + a[2][i-1]
	}
	right := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		right[i] = p2[i] + suf[i]
	}

	st := newSegTree(left)

	ans := int64(-1 << 60)
	for i := 1; i <= n; i++ {
		for _, of := range offers[i] {
			st.RangeAdd(0, of.l, -of.k)
		}
		cur := st.Query(0, i)
		if val := cur + right[i]; val > ans {
			ans = val
		}
		st.pointUpdate(i+1, cur)
	}
	fmt.Fprintln(out, ans)
}
