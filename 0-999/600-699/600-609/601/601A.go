package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(rail [][]bool, n int, useRail bool) int {
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[1] = 0
	q = append(q, 1)
	for head := 0; head < len(q); head++ {
		u := q[head]
		for v := 1; v <= n; v++ {
			if v == u {
				continue
			}
			if rail[u][v] == useRail && dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	rail := make([][]bool, n+1)
	for i := range rail {
		rail[i] = make([]bool, n+1)
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		rail[u][v] = true
		rail[v][u] = true
	}

	if rail[1][n] {
		// Train can go directly via railway, need bus on roads
		d := bfs(rail, n, false)
		fmt.Println(d)
	} else {
		// Bus can go directly via road, need train on rails
		d := bfs(rail, n, true)
		fmt.Println(d)
	}
}
