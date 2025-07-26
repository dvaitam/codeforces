package main

import (
	"bufio"
	"fmt"
	"os"
)

func bfs(start int, adj [][]int) ([]int, int) {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	far := start
	for head := 0; head < len(q); head++ {
		v := q[head]
		if dist[v] > dist[far] {
			far = v
		}
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist, far
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var k, c int64
		fmt.Fscan(reader, &n, &k, &c)
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		dist1, far1 := bfs(1, adj)
		distA, far2 := bfs(far1, adj)
		distB, _ := bfs(far2, adj)
		var best int64 = -1 << 63
		for v := 1; v <= n; v++ {
			ecc := distA[v]
			if distB[v] > ecc {
				ecc = distB[v]
			}
			profit := int64(ecc)*k - int64(dist1[v])*c
			if profit > best {
				best = profit
			}
		}
		fmt.Fprintln(writer, best)
	}
}
