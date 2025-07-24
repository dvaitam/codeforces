package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	vis := make([]bool, n+1)
	queue := make([]int, 0)
	for i := 1; i <= n; i++ {
		if vis[i] {
			continue
		}
		queue = append(queue[:0], i)
		vis[i] = true
		nodes := 0
		edges := 0
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			nodes++
			edges += len(adj[v])
			for _, u := range adj[v] {
				if !vis[u] {
					vis[u] = true
					queue = append(queue, u)
				}
			}
		}
		if edges/2 != nodes*(nodes-1)/2 {
			fmt.Fprintln(out, "NO")
			return
		}
	}
	fmt.Fprintln(out, "YES")
}
