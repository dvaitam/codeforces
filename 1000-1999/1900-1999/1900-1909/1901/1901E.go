package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		parent := make([]int, n)
		order := make([]int, 0, n)
		parent[0] = -1
		order = append(order, 0)
		for i := 0; i < len(order); i++ {
			v := order[i]
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				order = append(order, to)
			}
		}

		dp := make([]int64, n)
		bestRoot := make([]int64, n)
		for idx := n - 1; idx >= 0; idx-- {
			v := order[idx]
			childs := make([]int64, 0, len(g[v]))
			for _, to := range g[v] {
				if to == parent[v] {
					continue
				}
				childs = append(childs, dp[to])
			}
			sort.Slice(childs, func(i, j int) bool { return childs[i] > childs[j] })

			pref := int64(0)
			best := a[v]
			rootAns := a[v]
			for i, val := range childs {
				pref += val
				m := i + 1
				cand := pref
				if m != 1 {
					cand += a[v]
				}
				if cand > best {
					best = cand
				}
				candRoot := pref
				if m != 2 {
					candRoot += a[v]
				}
				if candRoot > rootAns {
					rootAns = candRoot
				}
			}
			dp[v] = best
			bestRoot[v] = rootAns
		}

		ans := int64(0)
		for i := 0; i < n; i++ {
			if bestRoot[i] > ans {
				ans = bestRoot[i]
			}
		}
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}
