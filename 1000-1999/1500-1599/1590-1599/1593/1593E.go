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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		g := make([][]int, n)
		deg := make([]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
			deg[u]++
			deg[v]++
		}
		if k == 0 {
			fmt.Fprintln(writer, n)
			continue
		}
		queue := make([]int, 0)
		time := make([]int, n)
		for i := 0; i < n; i++ {
			if deg[i] <= 1 {
				queue = append(queue, i)
				time[i] = 1
			}
		}
		for front := 0; front < len(queue); front++ {
			v := queue[front]
			for _, to := range g[v] {
				if time[to] > 0 {
					continue
				}
				deg[to]--
				if deg[to] == 1 {
					time[to] = time[v] + 1
					queue = append(queue, to)
				}
			}
		}
		ans := 0
		for i := 0; i < n; i++ {
			if time[i] == 0 || time[i] > k {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
