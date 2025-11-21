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

	const inf = int(1e9)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		adj := make([][]int, n+1)
		for i := 1; i < n; i++ {
			adj[i] = append(adj[i], i+1)
		}
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
		}

		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = inf
		}
		dist[1] = 0
		for i := 1; i <= n; i++ {
			if dist[i] == inf {
				continue
			}
			for _, v := range adj[i] {
				if dist[i]+1 < dist[v] {
					dist[v] = dist[i] + 1
				}
			}
		}

		d := dist[n]
		threshold := n - d - 1
		if threshold < 0 {
			threshold = 0
		}
		if threshold > n-1 {
			threshold = n - 1
		}
		res := make([]byte, n-1)
		for s := 1; s <= n-1; s++ {
			if s > threshold {
				res[s-1] = '1'
			} else {
				res[s-1] = '0'
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
