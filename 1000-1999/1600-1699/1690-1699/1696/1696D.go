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
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}
		g := make([][]int, n)
		// build edges naively
		for i := 0; i < n; i++ {
			mn, mx := a[i], a[i]
			for j := i + 1; j < n; j++ {
				if a[j] < mn {
					mn = a[j]
				}
				if a[j] > mx {
					mx = a[j]
				}
				if (mn == a[i] && mx == a[j]) || (mn == a[j] && mx == a[i]) {
					g[i] = append(g[i], j)
					g[j] = append(g[j], i)
				}
			}
		}
		// BFS
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := make([]int, 0, n)
		dist[0] = 0
		q = append(q, 0)
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for _, v := range g[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
		fmt.Fprintln(writer, dist[n-1])
	}
}
