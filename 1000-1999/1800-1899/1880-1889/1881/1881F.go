package main

import (
	"bufio"
	"fmt"
	"os"
)

// bfs returns distances from start node in an unweighted tree
func bfs(n int, g [][]int, start int) []int {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	q = append(q, start)
	dist[start] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		marks := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &marks[i])
			marks[i]-- // convert to 0-based
		}
		g := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}

		// find one endpoint A of the diameter among marked nodes
		d0 := bfs(n, g, marks[0])
		A := marks[0]
		for _, x := range marks {
			if d0[x] > d0[A] {
				A = x
			}
		}
		// distances from A
		distA := bfs(n, g, A)
		// find other endpoint B
		B := marks[0]
		for _, x := range marks {
			if distA[x] > distA[B] {
				B = x
			}
		}
		distB := bfs(n, g, B)

		// compute minimal maximum distance
		best := n // since n-1 is max possible distance
		for i := 0; i < n; i++ {
			d := distA[i]
			if distB[i] > d {
				d = distB[i]
			}
			if d < best {
				best = d
			}
		}
		fmt.Fprintln(out, best)
	}
}
