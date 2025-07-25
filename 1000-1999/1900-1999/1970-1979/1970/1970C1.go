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

	var n, t int
	if _, err := fmt.Fscan(in, &n, &t); err != nil {
		return
	}
	g := make([][]int, n+1)
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
		deg[u]++
		deg[v]++
	}
	starts := make([]int, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &starts[i])
	}
	start := starts[0]

	leaves := make([]int, 0, 2)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			leaves = append(leaves, i)
		}
	}

	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}

	d1 := dist[leaves[0]]
	d2 := dist[leaves[1]]
	if d1%2 == 1 || d2%2 == 1 {
		fmt.Fprintln(out, "Ron")
	} else {
		fmt.Fprintln(out, "Hermione")
	}
}
