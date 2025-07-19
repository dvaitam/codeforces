package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents a reversed graph edge with a weight (0 or 1).
type Edge struct {
	to, w int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	if n == 1 {
		fmt.Fprintln(writer, 0)
		fmt.Fprintln(writer, 0)
		return
	}
	// Build reversed graph
	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var x, y, z int
		fmt.Fscan(reader, &x, &y, &z)
		adj[y] = append(adj[y], Edge{to: x, w: z})
	}

	vis := make([]int, n+1)
	a := make([]int, n+1)
	d := make([]int, n+1)
	q := make([]int, n+1)
	H, T := 0, 1
	q[0] = n
	// Modified BFS: track visit flags and assign bits
	for H < T {
		k := q[H]
		H++
		for _, e := range adj[k] {
			to := e.to
			w := e.w
			if vis[to] == 3 {
				continue
			}
			vis[to] |= 1 << w
			if vis[to] == 3 {
				d[to] = d[k] + 1
				q[T] = to
				T++
			} else {
				a[to] = w ^ 1
			}
		}
	}

	// If node 1 not fully visited, no valid path
	if vis[1] != 3 {
		fmt.Fprintln(writer, -1)
		return
	}
	// Output distance and bitstring
	fmt.Fprintln(writer, d[1])
	for i := 1; i <= n; i++ {
		writer.WriteByte('0' + byte(a[i]))
	}
}
