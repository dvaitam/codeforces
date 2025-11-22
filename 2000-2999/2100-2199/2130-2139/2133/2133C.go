package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		g := make([][]int, n)
		indeg := make([]int, n)
		for u := 0; u < n; u++ {
			var k int
			fmt.Fscan(in, &k)
			for i := 0; i < k; i++ {
				var v int
				fmt.Fscan(in, &v)
				v--
				g[u] = append(g[u], v)
				indeg[v]++
			}
		}

		// topological order using Kahn
		q := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				q = append(q, i)
			}
		}
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range g[u] {
				indeg[v]--
				if indeg[v] == 0 {
					q = append(q, v)
				}
			}
		}

		// dp for longest path starting at node
		dp := make([]int, n)
		nxt := make([]int, n)
		for i := range nxt {
			nxt[i] = -1
		}
		for i := n - 1; i >= 0; i-- {
			u := q[i]
			bestLen := 1
			bestTo := -1
			for _, v := range g[u] {
				if 1+dp[v] > bestLen {
					bestLen = 1 + dp[v]
					bestTo = v
				}
			}
			dp[u] = bestLen
			nxt[u] = bestTo
		}

		// find start of overall longest path
		start := 0
		for i := 1; i < n; i++ {
			if dp[i] > dp[start] {
				start = i
			}
		}

		// reconstruct path
		path := make([]int, 0, dp[start])
		for cur := start; cur != -1; cur = nxt[cur] {
			path = append(path, cur+1) // convert back to 1-indexed
		}

		fmt.Fprint(out, len(path))
		for _, v := range path {
			fmt.Fprint(out, " ", v)
		}
		if T > 1 {
			fmt.Fprintln(out)
		}
	}
}
