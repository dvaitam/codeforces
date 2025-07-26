package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	g  [][]int
	sz []int
	dp []int
)

func dfs(u int) {
	sz[u] = 1
	childSizes := make([]int, 0, len(g[u]))
	for _, v := range g[u] {
		dfs(v)
		sz[u] += sz[v]
		dp[u] += dp[v]
		childSizes = append(childSizes, sz[v])
	}
	if len(childSizes) == 0 {
		return
	}
	S := sz[u] - 1
	reachable := make([]bool, S+1)
	reachable[0] = true
	for _, s := range childSizes {
		for j := S; j >= s; j-- {
			if reachable[j-s] {
				reachable[j] = true
			}
		}
	}
	best := 0
	for t := 0; t <= S; t++ {
		if reachable[t] {
			prod := t * (S - t)
			if prod > best {
				best = prod
			}
		}
	}
	dp[u] += best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		g[p] = append(g[p], i)
	}

	sz = make([]int, n+1)
	dp = make([]int, n+1)
	dfs(1)
	fmt.Fprintln(out, dp[1])
}
