package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxN = 100000

var (
	n   int
	g   [MaxN + 1][]int
	dp  [MaxN + 1][2]int
	vis [MaxN + 1]bool
)

func dfs(u, parent int) {
	vis[u] = true
	dp[u][1] = 1 // include u
	for _, v := range g[u] {
		if v == parent {
			continue
		}
		dfs(v, u)
		dp[u][0] += max(dp[v][0], dp[v][1])
		dp[u][1] += dp[v][0]
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
	}
	dfs(1, 0)
	ans := max(dp[1][0], dp[1][1])
	fmt.Println(ans)
}
