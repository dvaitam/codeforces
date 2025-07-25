package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxV = 100000
const inf = int(1e9)

type SegTree struct {
	n int
	t []int
}

func NewSegTree(n int) *SegTree {
	N := 1
	for N < n {
		N <<= 1
	}
	t := make([]int, 2*N)
	for i := range t {
		t[i] = inf
	}
	return &SegTree{n: N, t: t}
}

func (st *SegTree) Update(pos, val int) {
	i := pos + st.n
	st.t[i] = val
	for i >>= 1; i > 0; i >>= 1 {
		a, b := st.t[i<<1], st.t[i<<1|1]
		if a < b {
			st.t[i] = a
		} else {
			st.t[i] = b
		}
	}
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return inf
	}
	l += st.n
	r += st.n
	res := inf
	for l <= r {
		if l&1 == 1 {
			if st.t[l] < res {
				res = st.t[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.t[r] < res {
				res = st.t[r]
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		seg := NewSegTree(maxV + 2)
		vals := make([]int, maxV+2)
		for i := range vals {
			vals[i] = inf
		}
		ans := int64(0)
		for i := 0; i < n; i++ {
			x := a[i]
			l := x - k
			if l < 1 {
				l = 1
			}
			r := x + k
			if r > maxV {
				r = maxV
			}
			q := seg.Query(l, r)
			dp := i + 1
			if q < dp {
				dp = q
			}
			ans += int64((i + 1) - dp + 1)
			if dp < vals[x] {
				vals[x] = dp
				seg.Update(x, dp)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
