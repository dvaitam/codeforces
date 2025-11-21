package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type segTree struct {
	n   int
	val []int64
}

func newSegTree(n int) *segTree {
	size := 4*n + 5
	val := make([]int64, size)
	for i := range val {
		val[i] = inf
	}
	return &segTree{n: n, val: val}
}

// rangeMinAssign sets each position in [l, r] to min(current, v).
func (st *segTree) rangeMinAssign(node, l, r, ql, qr int, v int64) {
	if ql <= l && r <= qr {
		if v < st.val[node] {
			st.val[node] = v
		}
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		st.rangeMinAssign(node<<1, l, mid, ql, qr, v)
	}
	if qr > mid {
		st.rangeMinAssign(node<<1|1, mid+1, r, ql, qr, v)
	}
}

// pointQuery returns the minimum value assigned that affects position idx.
func (st *segTree) pointQuery(node, l, r, idx int) int64 {
	res := st.val[node]
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if idx <= mid {
		v := st.pointQuery(node<<1, l, mid, idx)
		if v < res {
			res = v
		}
		return res
	}
	v := st.pointQuery(node<<1|1, mid+1, r, idx)
	if v < res {
		res = v
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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := range b {
			fmt.Fscan(in, &b[i])
		}

		pref := make([]int64, n+1)
		for i, v := range a {
			pref[i+1] = pref[i] + v
		}

		// Maximum index that can ever be reached (ignoring scores).
		R := 1
		for i := 1; i <= R; i++ {
			if b[i-1] > R {
				R = b[i-1]
			}
		}
		if R > n {
			R = n
		}

		if R == 1 {
			fmt.Fprintln(out, a[0])
			continue
		}

		size := R - 1 // boundaries to cover: 1..R-1
		seg := newSegTree(size)
		dp := make([]int64, R)
		for i := range dp {
			dp[i] = inf
		}
		dp[0] = 0

		for i := 1; i <= R; i++ {
			val := dp[i-1] + a[i-1]
			l := i
			r := b[i-1] - 1
			if r > R-1 {
				r = R - 1
			}
			if l <= r {
				seg.rangeMinAssign(1, 1, size, l, r, val)
			}
			if i <= R-1 {
				dp[i] = seg.pointQuery(1, 1, size, i)
			}
		}

		ans := int64(0)
		for k := 1; k <= R; k++ {
			if dp[k-1] >= inf {
				continue
			}
			cur := pref[k] - dp[k-1]
			if cur > ans {
				ans = cur
			}
		}

		fmt.Fprintln(out, ans)
	}
}
