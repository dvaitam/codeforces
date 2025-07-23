package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	c  int
}

type State struct {
	v    int
	used int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		g[x] = append(g[x], Edge{y, z})
		g[y] = append(g[y], Edge{x, z})
	}
	var a, b int
	fmt.Fscan(in, &a, &b)

	visited := make([][2]bool, n+1)
	queue := make([]State, 0, n*2)
	queue = append(queue, State{a, 0})
	visited[a][0] = true
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		if cur.v == b && cur.used == 1 {
			fmt.Println("YES")
			return
		}
		for _, e := range g[cur.v] {
			nextUsed := cur.used | e.c
			if !visited[e.to][nextUsed] {
				visited[e.to][nextUsed] = true
				queue = append(queue, State{e.to, nextUsed})
			}
		}
	}
	fmt.Println("NO")
}
