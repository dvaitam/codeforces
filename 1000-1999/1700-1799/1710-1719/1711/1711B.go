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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		edges := make([][2]int, m)
		deg := make([]int, n+1)
		for i := 0; i < m; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			edges[i] = [2]int{x, y}
			deg[x]++
			deg[y]++
		}

		if m%2 == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}

		ans := int(1e18)
		for i := 1; i <= n; i++ {
			if deg[i]%2 == 1 && a[i] < ans {
				ans = a[i]
			}
		}
		for _, e := range edges {
			u, v := e[0], e[1]
			if deg[u]%2 == 0 && deg[v]%2 == 0 {
				if a[u]+a[v] < ans {
					ans = a[u] + a[v]
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
