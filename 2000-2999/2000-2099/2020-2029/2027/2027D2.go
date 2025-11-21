package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	MOD = 1_000_000_007
	INF = int64(1 << 60)
)

type segTree struct {
	n   int
	val []int64
	cnt []int64
	dp  []int64
	wt  []int64
}

func newSegTree(dp, wt []int64) *segTree {
	n := len(dp) - 1
	size := 4 * (n + 5)
	return &segTree{
		n:   n,
		val: make([]int64, size),
		cnt: make([]int64, size),
		dp:  dp,
		wt:  wt,
	}
}

func (st *segTree) build(node, l, r int) {
	if l == r {
		st.val[node] = st.dp[l]
		if st.dp[l] >= INF {
			st.cnt[node] = 0
		} else {
			st.cnt[node] = st.wt[l] % MOD
		}
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid)
	st.build(node<<1|1, mid+1, r)
	st.pull(node)
}

func (st *segTree) pull(node int) {
	lc := node << 1
	rc := lc | 1
	if st.val[lc] < st.val[rc] {
		st.val[node] = st.val[lc]
		st.cnt[node] = st.cnt[lc]
	} else if st.val[rc] < st.val[lc] {
		st.val[node] = st.val[rc]
		st.cnt[node] = st.cnt[rc]
	} else {
		st.val[node] = st.val[lc]
		if st.val[node] >= INF {
			st.cnt[node] = 0
		} else {
			st.cnt[node] = (st.cnt[lc] + st.cnt[rc]) % MOD
		}
	}
}

func (st *segTree) update(node, l, r, idx int, value, count int64) {
	if l == r {
		st.val[node] = value
		if value >= INF {
			st.cnt[node] = 0
		} else {
			st.cnt[node] = count % MOD
		}
		return
	}
	mid := (l + r) >> 1
	if idx <= mid {
		st.update(node<<1, l, mid, idx, value, count)
	} else {
		st.update(node<<1|1, mid+1, r, idx, value, count)
	}
	st.pull(node)
}

func (st *segTree) query(node, l, r, ql, qr int) (int64, int64) {
	if ql > r || qr < l {
		return INF, 0
	}
	if ql <= l && r <= qr {
		return st.val[node], st.cnt[node]
	}
	mid := (l + r) >> 1
	v1, c1 := st.query(node<<1, l, mid, ql, qr)
	v2, c2 := st.query(node<<1|1, mid+1, r, ql, qr)
	if v1 < v2 {
		return v1, c1
	} else if v2 < v1 {
		return v2, c2
	}
	if v1 >= INF {
		return INF, 0
	}
	return v1, (c1 + c2) % MOD
}

func solveCase(n, m int, a []int64, b []int64) (int64, int64, bool) {
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}

	dp := make([]int64, n+1)
	ways := make([]int64, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	ways[0] = 1

	st := newSegTree(dp, ways)
	for k := 1; k <= m; k++ {
		st.build(1, 0, n)
		for j := 1; j <= n; j++ {
			target := pref[j] - b[k-1]
			L := sort.Search(n+1, func(i int) bool { return pref[i] >= target })
			if L > j-1 {
				continue
			}
			minVal, cnt := st.query(1, 0, n, L, j-1)
			if minVal >= INF || cnt == 0 {
				continue
			}
			newCost := minVal + int64(m-k)
			if newCost < dp[j] {
				dp[j] = newCost
				ways[j] = cnt % MOD
				st.update(1, 0, n, j, dp[j], ways[j])
			} else if newCost == dp[j] {
				ways[j] = (ways[j] + cnt) % MOD
				st.update(1, 0, n, j, dp[j], ways[j])
			}
		}
	}

	if dp[n] >= INF {
		return 0, 0, false
	}
	return dp[n], ways[n] % MOD, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}
		cost, cnt, ok := solveCase(n, m, a, b)
		if !ok {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintf(out, "%d %d\n", cost, cnt)
		}
	}
}
