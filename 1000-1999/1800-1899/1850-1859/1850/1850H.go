package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to int
	w  int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		adj := make([][]edge, n+1)
		for i := 0; i < m; i++ {
			var a, b int
			var d int64
			fmt.Fscan(reader, &a, &b, &d)
			adj[b] = append(adj[b], edge{to: a, w: d})
			adj[a] = append(adj[a], edge{to: b, w: -d})
		}
		pos := make([]int64, n+1)
		vis := make([]bool, n+1)
		queue := make([]int, 0)
		ok := true
		for i := 1; i <= n && ok; i++ {
			if !vis[i] {
				vis[i] = true
				pos[i] = 0
				queue = append(queue, i)
				for len(queue) > 0 && ok {
					u := queue[0]
					queue = queue[1:]
					for _, e := range adj[u] {
						v := e.to
						val := pos[u] + e.w
						if !vis[v] {
							vis[v] = true
							pos[v] = val
							queue = append(queue, v)
						} else if pos[v] != val {
							ok = false
							break
						}
					}
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
