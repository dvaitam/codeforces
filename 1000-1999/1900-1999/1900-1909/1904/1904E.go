package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	removed := make([]bool, n+1)
	visited := make([]bool, n+1)
	dist := make([]int, n+1)
	queue := make([]int, 0, n)

	for ; q > 0; q-- {
		var x, k int
		fmt.Fscan(reader, &x, &k)
		nodes := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &nodes[i])
			removed[nodes[i]] = true
		}

		queue = queue[:0]
		queue = append(queue, x)
		visited[x] = true
		dist[x] = 0
		best := 0
		for front := 0; front < len(queue); front++ {
			u := queue[front]
			if dist[u] > best {
				best = dist[u]
			}
			for _, v := range adj[u] {
				if visited[v] || removed[v] {
					continue
				}
				visited[v] = true
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}

		// reset visited and dist
		for _, u := range queue {
			visited[u] = false
			dist[u] = 0
		}
		for _, a := range nodes {
			removed[a] = false
		}

		fmt.Fprintln(writer, best)
	}
}
