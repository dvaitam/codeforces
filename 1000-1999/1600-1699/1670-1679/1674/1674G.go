package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adj := make([][]int, n+1)
	indeg := make([]int, n+1)
	outdeg := make([]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		indeg[v]++
		outdeg[u]++
	}

	indeg2 := make([]int, n+1)
	copy(indeg2, indeg)
	queue := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if indeg2[i] == 0 {
			queue = append(queue, i)
		}
	}
	order := make([]int, 0, n)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, to := range adj[v] {
			indeg2[to]--
			if indeg2[to] == 0 {
				queue = append(queue, to)
			}
		}
	}

	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = 1
	}
	ans := 1
	for _, v := range order {
		for _, to := range adj[v] {
			if outdeg[v] > 1 && indeg[to] > 1 {
				if dp[v]+1 > dp[to] {
					dp[to] = dp[v] + 1
					if dp[to] > ans {
						ans = dp[to]
					}
				}
			}
		}
	}
	fmt.Println(ans)
}
