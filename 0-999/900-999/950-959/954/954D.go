package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, adj [][]int, n int) []int {
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s, t int
	if _, err := fmt.Fscan(in, &n, &m, &s, &t); err != nil {
		return
	}
	adj := make([][]int, n+1)
	edges := make([][]bool, n+1)
	for i := 1; i <= n; i++ {
		edges[i] = make([]bool, n+1)
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		edges[u][v] = true
		edges[v][u] = true
	}

	ds := bfs(s, adj, n)
	dt := bfs(t, adj, n)
	base := ds[t]

	cnt := 0
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if edges[i][j] {
				continue
			}
			if ds[i]+1+dt[j] >= base && ds[j]+1+dt[i] >= base {
				cnt++
			}
		}
	}
	fmt.Fprintln(out, cnt)
}
