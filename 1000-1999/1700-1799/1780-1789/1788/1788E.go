package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n, tree: make([]int, n+2)}
	for i := range b.tree {
		b.tree[i] = -1 << 60
	}
	return b
}

func (b *BIT) update(idx int, val int) {
	for idx <= b.n {
		if val > b.tree[idx] {
			b.tree[idx] = val
		}
		idx += idx & -idx
	}
}

func (b *BIT) query(idx int) int {
	res := -1 << 60
	for idx > 0 {
		if b.tree[idx] > res {
			res = b.tree[idx]
		}
		idx &= idx - 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + arr[i-1]
	}

	comp := make([]int64, n+1)
	copy(comp, pref)
	sort.Slice(comp, func(i, j int) bool { return comp[i] < comp[j] })
	m := 0
	for i := 0; i <= n; i++ {
		if m == 0 || comp[i] != comp[m-1] {
			comp[m] = comp[i]
			m++
		}
	}

	idx := func(x int64) int {
		return sort.Search(m, func(i int) bool { return comp[i] >= x }) + 1
	}

	bit := NewBIT(m)
	dp := make([]int, n+1)

	bit.update(idx(0), 0) // dp[0]-0

	for i := 1; i <= n; i++ {
		id := idx(pref[i])
		best := bit.query(id)
		dp[i] = dp[i-1]
		if best != -1<<60 {
			cand := best + i
			if cand > dp[i] {
				dp[i] = cand
			}
		}
		bit.update(id, dp[i]-i)
	}

	fmt.Fprintln(out, dp[n])
}
