package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	P   int64 = 10000000 + 19
	MOD int64 = 1000000000 + 7
)

type Item struct {
	val int64
	w   int
}

var (
	k       int
	tree    [][]Item
	evtType []int
	ans     []int64
	stack   [][]int64
)

func addItem(idx, l, r, ql, qr int, it Item) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		tree[idx] = append(tree[idx], it)
		return
	}
	mid := (l + r) >> 1
	addItem(idx<<1, l, mid, ql, qr, it)
	addItem(idx<<1|1, mid+1, r, ql, qr, it)
}

func apply(dp []int64, it Item) {
	w := it.w
	val := it.val
	for i := k; i >= w; i-- {
		if tmp := dp[i-w] + val; tmp > dp[i] {
			dp[i] = tmp
		}
	}
}

func dfs(idx, l, r, depth int, dp []int64) {
	copy(stack[depth], dp)
	for _, it := range tree[idx] {
		apply(dp, it)
	}
	if l == r {
		if evtType[l] == 3 {
			var res int64
			pow := int64(1)
			for i := 1; i <= k; i++ {
				res = (res + dp[i]%MOD*pow) % MOD
				pow = pow * P % MOD
			}
			ans[l] = res
		}
	} else {
		mid := (l + r) >> 1
		dfs(idx<<1, l, mid, depth+1, dp)
		dfs(idx<<1|1, mid+1, r, depth+1, dp)
	}
	copy(dp, stack[depth])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	type Info struct {
		val   int64
		w     int
		start int
		end   int
	}

	items := make([]Info, n+1)
	for i := 1; i <= n; i++ {
		var v, w int
		fmt.Fscan(in, &v, &w)
		items[i] = Info{val: int64(v), w: w, start: 1}
	}

	var q int
	fmt.Fscan(in, &q)

	evtType = make([]int, q+1)
	ans = make([]int64, q+1)
	nextID := n
	// we append new items as they appear
	for t := 1; t <= q; t++ {
		var tp int
		fmt.Fscan(in, &tp)
		evtType[t] = tp
		switch tp {
		case 1:
			var v, w int
			fmt.Fscan(in, &v, &w)
			nextID++
			items = append(items, Info{val: int64(v), w: w, start: t})
		case 2:
			var x int
			fmt.Fscan(in, &x)
			if x > 0 && x < len(items) {
				items[x].end = t
			}
		case 3:
			// nothing else
		}
	}

	for i := 1; i <= nextID; i++ {
		if items[i].start == 0 {
			continue
		}
		if items[i].end == 0 {
			items[i].end = q + 1
		}
	}

	tree = make([][]Item, 4*(q+2))
	for i := 1; i <= nextID; i++ {
		st := items[i].start
		ed := items[i].end - 1
		if st <= ed {
			addItem(1, 1, q, st, ed, Item{val: items[i].val, w: items[i].w})
		}
	}

	stack = make([][]int64, 20)
	for i := range stack {
		stack[i] = make([]int64, k+1)
	}
	dp := make([]int64, k+1)
	if q > 0 {
		dfs(1, 1, q, 0, dp)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= q; i++ {
		if evtType[i] == 3 {
			fmt.Fprintln(out, ans[i])
		}
	}
}
