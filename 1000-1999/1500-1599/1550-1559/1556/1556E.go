package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

type segtree struct {
	n  int
	mn []int64
	mx []int64
}

func newSegTree(arr []int64) *segtree {
	n := len(arr)
	st := &segtree{
		n:  n,
		mn: make([]int64, 4*n),
		mx: make([]int64, 4*n),
	}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segtree) build(node, l, r int, arr []int64) {
	if l == r {
		st.mn[node] = arr[l]
		st.mx[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, arr)
	st.build(node*2+1, mid+1, r, arr)
	if st.mn[node*2] < st.mn[node*2+1] {
		st.mn[node] = st.mn[node*2]
	} else {
		st.mn[node] = st.mn[node*2+1]
	}
	if st.mx[node*2] > st.mx[node*2+1] {
		st.mx[node] = st.mx[node*2]
	} else {
		st.mx[node] = st.mx[node*2+1]
	}
}

func (st *segtree) query(node, l, r, L, R int) (int64, int64) {
	if R < l || r < L {
		return INF, -INF
	}
	if L <= l && r <= R {
		return st.mn[node], st.mx[node]
	}
	mid := (l + r) / 2
	mn1, mx1 := st.query(node*2, l, mid, L, R)
	mn2, mx2 := st.query(node*2+1, mid+1, r, L, R)
	mn := mn1
	if mn2 < mn {
		mn = mn2
	}
	mx := mx1
	if mx2 > mx {
		mx = mx2
	}
	return mn, mx
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + a[i] - b[i]
	}

	// build segment tree over prefix[1..n]
	st := newSegTree(pref[1:])

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		if pref[l-1] != pref[r] {
			fmt.Fprintln(out, -1)
			continue
		}
		mn, mx := st.query(1, 0, n-1, l-1, r-1)
		if mx > pref[l-1] {
			fmt.Fprintln(out, -1)
			continue
		}
		if mn > pref[l-1] {
			// Shouldn't happen because mx <= pref[l-1], but just in case.
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, pref[l-1]-mn)
	}
}
