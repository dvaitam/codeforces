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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	color := make([]int, n+1)
	for i := range color {
		color[i] = -1
	}

	// BFS to color the tree
	queue := make([]int, 0, n)
	queue = append(queue, 1)
	color[1] = 0
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range adj[u] {
			if color[v] == -1 {
				color[v] = color[u] ^ 1
				queue = append(queue, v)
			}
		}
	}

	var c0, c1 int64
	for i := 1; i <= n; i++ {
		if color[i] == 0 {
			c0++
		} else {
			c1++
		}
	}

	res := c0*c1 - int64(n-1)
	fmt.Fprintln(out, res)
}
