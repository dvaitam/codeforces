package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		var k int
		fmt.Fscan(reader, &k)
		chips := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &chips[i])
		}
		// multi-source BFS assigning vertices to nearest chip (tie break by index order)
		dist := make([]int, n+1)
		owner := make([]int, n+1)
		for i := 1; i <= n; i++ {
			dist[i] = -1
		}
		q := make([]int, 0, n)
		for i, v := range chips {
			idx := i + 1
			dist[v] = 0
			owner[v] = idx
			q = append(q, v)
		}
		head := 0
		for head < len(q) {
			v := q[head]
			head++
			for _, to := range adj[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					owner[to] = owner[v]
					q = append(q, to)
				} else if dist[to] == dist[v]+1 && owner[to] > owner[v] {
					owner[to] = owner[v]
				}
			}
		}
		maxDist := make([]int, k+1)
		for i := 1; i <= n; i++ {
			id := owner[i]
			if id >= 1 {
				if dist[i] > maxDist[id] {
					maxDist[id] = dist[i]
				}
			}
		}
		r := maxDist[1]
		j := 1
		for i := 1; i <= k; i++ {
			if maxDist[i] < r {
				r = maxDist[i]
				j = i
			}
		}
		moves := r*k + j - 1
		if moves > n-k {
			moves = n - k
		}
		fmt.Fprintln(writer, moves)
	}
}
