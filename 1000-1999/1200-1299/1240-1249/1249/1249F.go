package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = int(1e9)

var (
	n, k int
	w    []int
	g    [][]int
	dp   [][]int
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dfs(u, p int) {
	dp[u] = make([]int, k+2)
	for i := range dp[u] {
		dp[u][i] = -INF
	}
	dp[u][0] = w[u]
	dp[u][k+1] = 0

	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u)
		ndp := make([]int, k+2)
		for i := range ndp {
			ndp[i] = -INF
		}
		for du := 0; du <= k+1; du++ {
			if dp[u][du] < -INF/2 {
				continue
			}
			for dv := 0; dv <= k+1; dv++ {
				if dp[v][dv] < -INF/2 {
					continue
				}
				distU := du
				distV := dv
				if distU > k {
					distU = INF
				}
				if distV > k {
					distV = INF
				}
				if distU+distV+1 <= k {
					continue
				}
				t := dv + 1
				if t > k+1 {
					t = k + 1
				}
				newDist := du
				if t < newDist {
					newDist = t
				}
				val := dp[u][du] + dp[v][dv]
				if val > ndp[newDist] {
					ndp[newDist] = val
				}
			}
		}
		dp[u] = ndp
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &k)
	w = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &w[i])
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dp = make([][]int, n+1)
	dfs(1, 0)
	ans := 0
	for i := 0; i <= k+1; i++ {
		if dp[1][i] > ans {
			ans = dp[1][i]
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
