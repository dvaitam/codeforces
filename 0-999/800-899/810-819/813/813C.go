package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, adj [][]int) []int {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0, n)
	queue = append(queue, start)
	dist[start] = 0
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, x int
	if _, err := fmt.Fscan(reader, &n, &x); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}

	distAlice := bfs(1, adj)
	distBob := bfs(x, adj)

	ans := 0
	for i := 1; i <= n; i++ {
		if distAlice[i] > distBob[i] {
			if distAlice[i]*2 > ans {
				ans = distAlice[i] * 2
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
