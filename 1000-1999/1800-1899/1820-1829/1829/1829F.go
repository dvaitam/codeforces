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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		deg := make([]int, n+1)
		for i := 1; i <= n; i++ {
			deg[i] = len(adj[i])
		}
		leafNeighbors := make([]int, n+1)
		for i := 1; i <= n; i++ {
			cnt := 0
			for _, v := range adj[i] {
				if deg[v] == 1 {
					cnt++
				}
			}
			leafNeighbors[i] = cnt
		}
		x := 0
		y := 0
		for i := 1; i <= n; i++ {
			if leafNeighbors[i] > 0 {
				x++
				y = leafNeighbors[i]
			}
		}
		fmt.Fprintf(out, "%d %d\n", x, y)
	}
}
