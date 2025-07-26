package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	idx int
}

type Node struct {
	id    int
	idx   int
	level int
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
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]Edge, n+1)
		for i := 1; i <= n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], Edge{v, i})
			adj[v] = append(adj[v], Edge{u, i})
		}

		visited := make([]bool, n+1)
		queue := []Node{{1, 0, 1}}
		visited[1] = true
		ans := 1

		for head := 0; head < len(queue); head++ {
			cur := queue[head]
			for _, e := range adj[cur.id] {
				if visited[e.to] {
					continue
				}
				lvl := cur.level
				if e.idx < cur.idx {
					lvl++
				}
				if lvl > ans {
					ans = lvl
				}
				visited[e.to] = true
				queue = append(queue, Node{e.to, e.idx, lvl})
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
