package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	cost int
}

func dfs(v, parent int, dist int, adj [][]Edge, max *int) {
	if dist > *max {
		*max = dist
	}
	for _, e := range adj[v] {
		if e.to != parent {
			dfs(e.to, v, dist+e.cost, adj, max)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		adj[u] = append(adj[u], Edge{v, c})
		adj[v] = append(adj[v], Edge{u, c})
	}
	maxDist := 0
	dfs(0, -1, 0, adj, &maxDist)
	fmt.Println(maxDist)
}
