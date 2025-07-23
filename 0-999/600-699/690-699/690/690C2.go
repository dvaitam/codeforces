package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, adj [][]int) (int, int) {
	n := len(adj) - 1
	dist := make([]int, n+1)
	visited := make([]bool, n+1)
	queue := make([]int, 0, n)
	queue = append(queue, start)
	visited[start] = true
	far := start
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		for _, v := range adj[u] {
			if !visited[v] {
				visited[v] = true
				dist[v] = dist[u] + 1
				if dist[v] > dist[far] {
					far = v
				}
				queue = append(queue, v)
			}
		}
	}
	return far, dist[far]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}

	u, _ := bfs(1, adj)
	_, dist := bfs(u, adj)
	fmt.Fprintln(out, dist)
}
