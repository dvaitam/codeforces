package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents a directed edge with weight
type Edge struct {
	v, w int
}

var (
	n, m, T int
	adj     [][]Edge
	dp      []map[int]int
	visited []bool
)

// dfs computes dp[u]: map[k] = min weight to reach n using k nodes from u
func dfs(u int) {
	if visited[u] {
		return
	}
	visited[u] = true
	dp[u] = make(map[int]int)
	if u == n {
		dp[u][1] = 0
		return
	}
	for _, e := range adj[u] {
		dfs(e.v)
		for k, val := range dp[e.v] {
			wsum := val + e.w
			if wsum > T {
				continue
			}
			nk := k + 1
			if prev, ok := dp[u][nk]; !ok || wsum < prev {
				dp[u][nk] = wsum
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for {
		if _, err := fmt.Fscan(reader, &n, &m, &T); err != nil {
			break
		}
		adj = make([][]Edge, n+1)
		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Fscan(reader, &u, &v, &w)
			adj[u] = append(adj[u], Edge{v: v, w: w})
		}
		dp = make([]map[int]int, n+1)
		visited = make([]bool, n+1)
		dfs(1)
		// find maximum k reachable from 1
		best := 0
		for k := range dp[1] {
			if k > best {
				best = k
			}
		}
		fmt.Fprintln(writer, best)
		// reconstruct path
		path := make([]int, 0, best)
		cur, k := 1, best
		for cur != n {
			path = append(path, cur)
			for _, e := range adj[cur] {
				if nextVal, ok := dp[e.v][k-1]; ok && nextVal+e.w == dp[cur][k] {
					cur = e.v
					k--
					break
				}
			}
		}
		path = append(path, n)
		for i, x := range path {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, x)
		}
		fmt.Fprintln(writer)
	}
}
