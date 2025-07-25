package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int = -1 << 60

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// segTree supports point updates and range maximum queries.
type segTree struct {
	n    int
	tree []int
}

func newSegTree(size int) *segTree {
	n := 1
	for n < size {
		n <<= 1
	}
	tree := make([]int, 2*n)
	for i := range tree {
		tree[i] = negInf
	}
	return &segTree{n, tree}
}

func (st *segTree) update(pos, val int) {
	pos += st.n
	st.tree[pos] = val
	for pos > 1 {
		pos >>= 1
		st.tree[pos] = max(st.tree[pos<<1], st.tree[pos<<1|1])
	}
}

func (st *segTree) query(l, r int) int {
	if l > r {
		return negInf
	}
	l += st.n
	r += st.n
	res := negInf
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
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	first := make([]int, n+1)
	for i := 0; i <= n; i++ {
		first[i] = n
	}
	for i := n - 1; i >= 0; i-- {
		first[a[i]] = i
	}

	st := newSegTree(n)
	dp := make([]int, n)
	st.update(0, 0) // dp[0]-0

	for i := 1; i < n; i++ {
		l := first[a[i]]
		best := negInf
		if l < i {
			q := st.query(l, i-1)
			if q != negInf {
				best = q + i - 1
			}
		}
		dp[i] = dp[i-1]
		if best > dp[i] {
			dp[i] = best
		}
		st.update(i, dp[i]-i)
	}

	if n == 1 {
		fmt.Println(0)
		return
	}
	fmt.Println(dp[n-1])
}
