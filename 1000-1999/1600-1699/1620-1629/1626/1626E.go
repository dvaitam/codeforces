package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, g [][]int) []int {
	n := len(g)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	col := make([]int, n)
	blacks := make([]int, 0)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &col[i])
		if col[i] == 1 {
			blacks = append(blacks, i)
		}
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	if len(blacks) >= 3 {
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(1)
		}
		fmt.Println()
		return
	}

	if len(blacks) != 2 {
		// Problem guarantees at least two black vertices
		return
	}

	b1 := blacks[0]
	b2 := blacks[1]
	d1 := bfs(b1, g)
	d2 := bfs(b2, g)
	dist := d1[b2]
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		t := (d1[i] + d2[i] - dist) / 2
		w := d1[i] - t
		if w <= 1 || dist-w <= 1 {
			ans[i] = 1
		} else {
			ans[i] = 0
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ans[i])
	}
	fmt.Println()
}
