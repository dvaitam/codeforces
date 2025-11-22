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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		g := make([][]int, n)
		indeg := make([]int, n)

		for i := 0; i < n-1; i++ {
			var u, v int
			var x, y int
			fmt.Fscan(in, &u, &v, &x, &y)
			u--
			v--
			// We want the endpoint with larger assigned value to be the one that gives greater contribution.
			if x >= y {
				// prefer p[u] > p[v]
				g[u] = append(g[u], v)
				indeg[v]++
			} else {
				// prefer p[v] > p[u]
				g[v] = append(g[v], u)
				indeg[u]++
			}
		}

		// Topological sort (graph is a directed tree, hence acyclic).
		queue := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if indeg[i] == 0 {
				queue = append(queue, i)
			}
		}
		order := make([]int, 0, n)
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			order = append(order, u)
			for _, v := range g[u] {
				indeg[v]--
				if indeg[v] == 0 {
					queue = append(queue, v)
				}
			}
		}

		// Assign descending numbers along topological order so that all directed edges go from larger to smaller.
		perm := make([]int, n)
		val := n
		for _, v := range order {
			perm[v] = val
			val--
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, perm[i])
		}
		fmt.Fprintln(out)
	}
}

