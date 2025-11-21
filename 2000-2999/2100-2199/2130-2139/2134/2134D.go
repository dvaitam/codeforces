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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		if n == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		found := false
		for i := 1; i <= n; i++ {
			if len(adj[i]) >= 3 {
				a := adj[i][0]
				c := adj[i][1]
				fmt.Fprintf(out, "%d %d %d\n", a, i, c)
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintln(out, -1)
		}
	}
}
