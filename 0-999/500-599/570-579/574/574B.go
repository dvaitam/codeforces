package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	// adjacency matrix for quick edge lookup
	adj := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		adj[i] = make([]bool, n+1)
	}
	deg := make([]int, n+1)
	neigh := make([][]int, n+1)

	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if !adj[a][b] {
			adj[a][b] = true
			adj[b][a] = true
			deg[a]++
			deg[b]++
			neigh[a] = append(neigh[a], b)
			neigh[b] = append(neigh[b], a)
		}
	}

	inf := int(1<<31 - 1)
	ans := inf

	// For each vertex consider pairs of its neighbors
	for u := 1; u <= n; u++ {
		nu := neigh[u]
		l := len(nu)
		for i := 0; i < l; i++ {
			v := nu[i]
			for j := i + 1; j < l; j++ {
				w := nu[j]
				if adj[v][w] {
					sum := deg[u] + deg[v] + deg[w] - 6
					if sum < ans {
						ans = sum
					}
				}
			}
		}
	}

	if ans == inf {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
