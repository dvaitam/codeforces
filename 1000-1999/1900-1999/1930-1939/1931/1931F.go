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
		var n, k int
		fmt.Fscan(in, &n, &k)
		adj := make([][]int, n+1)
		indeg := make([]int, n+1)
		for i := 0; i < k; i++ {
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &arr[j])
			}
			for j := 1; j < n-1; j++ {
				u, v := arr[j], arr[j+1]
				adj[u] = append(adj[u], v)
				indeg[v]++
			}
		}
		q := make([]int, 0, n)
		for i := 1; i <= n; i++ {
			if indeg[i] == 0 {
				q = append(q, i)
			}
		}
		cnt := 0
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			cnt++
			for _, v := range adj[u] {
				indeg[v]--
				if indeg[v] == 0 {
					q = append(q, v)
				}
			}
		}
		if cnt == n {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
