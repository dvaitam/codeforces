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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]int, n)
		pre := make([][]int, n)
		indeg := make([]int, n)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(in, &k)
			pre[i] = make([]int, k)
			for j := 0; j < k; j++ {
				var x int
				fmt.Fscan(in, &x)
				x--
				pre[i][j] = x
				adj[x] = append(adj[x], i)
				indeg[i]++
			}
		}

		queue := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		order := make([]int, 0, n)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			order = append(order, u)
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}
		if len(order) != n {
			fmt.Fprintln(out, -1)
			continue
		}

		dp := make([]int, n)
		ans := 0
		for _, u := range order {
			dp[u] = 1
			for _, v := range pre[u] {
				cand := dp[v]
				if v > u {
					cand++
				}
				if cand > dp[u] {
					dp[u] = cand
				}
			}
			if dp[u] > ans {
				ans = dp[u]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
