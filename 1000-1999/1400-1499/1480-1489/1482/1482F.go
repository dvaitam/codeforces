package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type Edge struct {
	u, v int
	w    int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = inf
			}
		}
	}
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		if w < dist[u][v] {
			dist[u][v] = w
			dist[v][u] = w
		}
		edges[i] = Edge{u, v, w}
	}
	// Floyd-Warshall for all-pairs shortest paths
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			dik := dist[i][k]
			if dik == inf {
				continue
			}
			for j := 0; j < n; j++ {
				val := dik + dist[k][j]
				if val < dist[i][j] {
					dist[i][j] = val
				}
			}
		}
	}

	var q int
	fmt.Fscan(in, &q)
	best := make([][]int64, n)
	for i := 0; i < n; i++ {
		best[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			best[i][j] = -inf
		}
	}
	used := make([]bool, n)
	for ; q > 0; q-- {
		var u, v int
		var l int64
		fmt.Fscan(in, &u, &v, &l)
		u--
		v--
		used[u] = true
		used[v] = true
		for t := 0; t < n; t++ {
			if l-dist[v][t] > best[u][t] {
				best[u][t] = l - dist[v][t]
			}
			if l-dist[u][t] > best[v][t] {
				best[v][t] = l - dist[u][t]
			}
		}
	}

	ans := 0
	for _, e := range edges {
		useful := false
		for u := 0; u < n; u++ {
			if !used[u] {
				continue
			}
			if dist[u][e.u]+e.w <= best[u][e.v] || dist[u][e.v]+e.w <= best[u][e.u] {
				useful = true
				break
			}
		}
		if useful {
			ans++
		}
	}
	fmt.Println(ans)
}
