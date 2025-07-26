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
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		required := make([]bool, n)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &a[i])
			a[i]--
			required[a[i]] = true
		}
		required[x] = true
		required[y] = true

		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		deg := make([]int, n)
		for i := 0; i < n; i++ {
			deg[i] = len(adj[i])
		}
		removed := make([]bool, n)
		queue := make([]int, 0)
		for i := 0; i < n; i++ {
			if deg[i] == 1 && !required[i] {
				removed[i] = true
				queue = append(queue, i)
			}
		}
		for h := 0; h < len(queue); h++ {
			v := queue[h]
			for _, to := range adj[v] {
				if removed[to] {
					continue
				}
				deg[to]--
				if deg[to] == 1 && !required[to] {
					removed[to] = true
					queue = append(queue, to)
				}
			}
		}

		edges := 0
		for i := 0; i < n; i++ {
			if removed[i] {
				continue
			}
			edges += deg[i]
		}
		edges /= 2

		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{x}
		dist[x] = 0
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			if v == y {
				break
			}
			for _, to := range adj[v] {
				if removed[to] || dist[to] != -1 {
					continue
				}
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}

		ans := 2*edges - dist[y]
		fmt.Fprintln(out, ans)
	}
}
