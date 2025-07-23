package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	to   int
	cost int
}

var (
	n, k int
	g    [][]edge
	dp   []int64
)

func dfs(u, p int) {
	gains := make([]int64, 0, len(g[u]))
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		dfs(e.to, u)
		gains = append(gains, int64(e.cost)+dp[e.to])
	}
	sort.Slice(gains, func(i, j int) bool { return gains[i] > gains[j] })
	limit := k
	if len(gains) < limit {
		limit = len(gains)
	}
	var sum int64
	for i := 0; i < limit; i++ {
		sum += gains[i]
	}
	dp[u] = sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	g = make([][]edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		g[u] = append(g[u], edge{v, c})
		g[v] = append(g[v], edge{u, c})
	}
	dp = make([]int64, n)
	dfs(0, -1)
	fmt.Println(dp[0])
}
