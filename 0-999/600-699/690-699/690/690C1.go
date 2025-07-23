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
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if a >= 1 && a <= n && b >= 1 && b <= n {
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
	}

	if m != n-1 {
		fmt.Fprintln(writer, "no")
		return
	}

	visited := make([]bool, n+1)
	queue := make([]int, 0, n)
	queue = append(queue, 1)
	visited[1] = true
	for len(queue) > 0 {
		x := queue[0]
		queue = queue[1:]
		for _, y := range adj[x] {
			if !visited[y] {
				visited[y] = true
				queue = append(queue, y)
			}
		}
	}

	for i := 1; i <= n; i++ {
		if !visited[i] {
			fmt.Fprintln(writer, "no")
			return
		}
	}

	fmt.Fprintln(writer, "yes")
}
