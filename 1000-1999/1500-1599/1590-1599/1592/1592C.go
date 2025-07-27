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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n+1)
		total := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			total ^= a[i]
		}
		g := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		if total == 0 {
			fmt.Fprintln(writer, "YES")
			continue
		}
		if k < 3 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		visited := make([]bool, n+1)
		count := 0
		var dfs func(int) int
		dfs = func(u int) int {
			visited[u] = true
			xor := a[u]
			for _, v := range g[u] {
				if !visited[v] {
					xor ^= dfs(v)
				}
			}
			if xor == total {
				count++
				return 0
			}
			return xor
		}
		dfs(1)

		if count >= 2 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
