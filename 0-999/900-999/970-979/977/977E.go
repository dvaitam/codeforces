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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	visited := make([]bool, n+1)
	ans := 0
	for start := 1; start <= n; start++ {
		if visited[start] {
			continue
		}
		queue := []int{start}
		visited[start] = true
		good := true
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			if len(adj[v]) != 2 {
				good = false
			}
			for _, to := range adj[v] {
				if !visited[to] {
					visited[to] = true
					queue = append(queue, to)
				}
			}
		}
		if good {
			ans++
		}
	}

	fmt.Fprintln(writer, ans)
}
