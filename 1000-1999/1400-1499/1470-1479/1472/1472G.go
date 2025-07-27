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
		var n, m int
		fmt.Fscan(in, &n, &m)
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
		}

		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := make([]int, 0, n)
		q = append(q, 0)
		dist[0] = 0
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range g[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}

		ord := make([]int, n)
		for i := range ord {
			ord[i] = i
		}
		sort.Slice(ord, func(i, j int) bool {
			return dist[ord[i]] > dist[ord[j]]
		})

		dp := make([]int, n)
		copy(dp, dist)
		for _, u := range ord {
			for _, v := range g[u] {
				if dist[u] < dist[v] {
					if dp[v] < dp[u] {
						dp[u] = dp[v]
					}
				} else {
					if dist[v] < dp[u] {
						dp[u] = dist[v]
					}
				}
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, dp[i])
		}
		fmt.Fprintln(out)
	}
}
