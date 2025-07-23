package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]Edge, n+1)
	total := 0
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
		total += w
	}
	var s int
	fmt.Fscan(in, &s)
	var m int
	fmt.Fscan(in, &m)
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(in, &x)
	}

	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{s}
	dist[s] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if dist[e.to] == -1 {
				dist[e.to] = dist[v] + e.w
				q = append(q, e.to)
			}
		}
	}
	maxd := 0
	for i := 1; i <= n; i++ {
		if dist[i] > maxd {
			maxd = dist[i]
		}
	}
	ans := 2*total - maxd
	fmt.Println(ans)
}
