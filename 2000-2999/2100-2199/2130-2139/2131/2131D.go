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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]int, n)
		deg := make([]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			deg[u]++
			deg[v]++
		}
		if n <= 2 {
			fmt.Fprintln(out, 0)
			continue
		}
		totalLeaves := 0
		leafAdj := make([]int, n)
		for i := 0; i < n; i++ {
			if deg[i] == 1 {
				totalLeaves++
				neighbor := adj[i][0]
				leafAdj[neighbor]++
			}
		}
		ans := n
		for i := 0; i < n; i++ {
			cand := totalLeaves - leafAdj[i]
			if cand < ans {
				ans = cand
			}
		}
		fmt.Fprintln(out, ans)
	}
}
