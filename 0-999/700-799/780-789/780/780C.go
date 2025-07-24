package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
		deg[x]++
		deg[y]++
	}
	maxDeg := 0
	for i := 1; i <= n; i++ {
		if deg[i] > maxDeg {
			maxDeg = deg[i]
		}
	}
	k := maxDeg + 1
	colors := make([]int, n+1)
	// BFS
	type pair struct{ v, p int }
	q := make([]pair, 0, n)
	colors[1] = 1
	q = append(q, pair{1, 0})
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		c := 1
		for _, to := range adj[cur.v] {
			if to == cur.p {
				continue
			}
			for c == colors[cur.v] || c == colors[cur.p] {
				c++
			}
			colors[to] = c
			c++
			q = append(q, pair{to, cur.v})
		}
	}
	fmt.Println(k)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Print(" ")
		}
		fmt.Print(colors[i])
	}
	fmt.Println()
}
