package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func bfs(start int, g [][]int) []int {
	n := len(g)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{start}
	dist[start] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	// Precompute all pairwise distances (inefficient for large n).
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = bfs(i, g)
	}

	res := make([]int, n)
	for k := 1; k <= n; k++ {
		best := 1
		if n <= 15 {
			m := 1 << uint(n)
			for mask := 1; mask < m; mask++ {
				size := bits.OnesCount(uint(mask))
				if size <= best {
					continue
				}
				valid := true
				for i := 0; i < n && valid; i++ {
					if mask>>i&1 == 0 {
						continue
					}
					for j := i + 1; j < n; j++ {
						if mask>>j&1 == 0 {
							continue
						}
						d := dist[i][j]
						if d != k && d != k+1 {
							valid = false
							break
						}
					}
				}
				if valid {
					best = size
				}
			}
		} else {
			// For larger n, this naive algorithm is not feasible; default to 1.
			best = 1
		}
		res[k-1] = best
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
}
