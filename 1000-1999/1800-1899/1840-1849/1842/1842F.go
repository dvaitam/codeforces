package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = int(-1e18)

var (
	n  int
	g  [][]int
	sz []int
)

func dfsSize(v, p int) {
	sz[v] = 1
	for _, to := range g[v] {
		if to == p {
			continue
		}
		dfsSize(to, v)
		sz[v] += sz[to]
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(k int) int {
	dp := make([][]int, n)
	var dfs func(v, p int)
	dfs = func(v, p int) {
		dp[v] = make([]int, sz[v]+1)
		for i := range dp[v] {
			dp[v][i] = negInf
		}
		dp[v][0] = 0
		dp[v][1] = 0
		cur := 1
		for _, to := range g[v] {
			if to == p {
				continue
			}
			dfs(to, v)
			tmp := make([]int, cur+sz[to]+1)
			for i := range tmp {
				tmp[i] = negInf
			}
			for i := 0; i <= cur; i++ {
				if dp[v][i] == negInf {
					continue
				}
				for j := 0; j <= sz[to]; j++ {
					if dp[to][j] == negInf {
						continue
					}
					val := dp[v][i] + dp[to][j] + abs(2*j-k)
					if val > tmp[i+j] {
						tmp[i+j] = val
					}
				}
			}
			dp[v] = tmp
			cur += sz[to]
		}
	}
	dfs(0, -1)
	if k >= 0 && k < len(dp[0]) {
		return dp[0][k]
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	g = make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	sz = make([]int, n)
	dfsSize(0, -1)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for k := 0; k <= n; k++ {
		if k > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, solve(k))
	}
	fmt.Fprintln(writer)
}
