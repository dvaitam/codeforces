package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)

	// A hypercube of dimension n has 2^n vertices.
	m := 1 << n
	// It has n * 2^(n-1) edges.
	numEdges := n * (1 << (n - 1))
	if n == 0 { // Special case for n=0
		numEdges = 0
		m = 1
	}

	adj := make([][]int, m)
	for i := 0; i < numEdges; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// Step 1: BFS from vertex 0 to determine layers by distance.
	dist := make([]int, m)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	q := []int{0}
	head := 0
	nodesByDist := make([][]int, n+1)
	if m > 0 {
		nodesByDist[0] = append(nodesByDist[0], 0)
	}

	for head < len(q) {
		u := q[head]
		head++
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				if dist[v] <= n {
					nodesByDist[dist[v]] = append(nodesByDist[dist[v]], v)
				}
				q = append(q, v)
			}
		}
	}

	// Step 2: Assign hypercube coordinates (values) based on layers.
	// val[graph_vertex] = hypercube_vertex
	val := make([]int, m)
	if m > 0 {
		val[0] = 0
	}

	// Layer 1 (neighbors of vertex 0)
	if n > 0 && len(nodesByDist) > 1 {
		for i, v := range nodesByDist[1] {
			if i < n {
				val[v] = 1 << i
			}
		}
	}

	// Subsequent layers
	for d := 2; d <= n; d++ {
		for _, u := range nodesByDist[d] {
			// The value of a node is the OR of values of its neighbors in the previous layer.
			for _, v := range adj[u] {
				if dist[v] == d-1 {
					val[u] |= val[v]
				}
			}
		}
	}

	// Step 3: Construct the permutation array p.
	// p[hypercube_vertex] = graph_vertex
	p := make([]int, m)
	for i := 0; i < m; i++ {
		if val[i] < m {
			p[val[i]] = i
		}
	}

	for i := 0; i < m; i++ {
		fmt.Fprintf(out, "%d ", p[i])
	}
	fmt.Fprintln(out)

	// Step 4: Determine the coloring.
	// A solution for coloring exists only if n is a power of 2.
	if n > 0 && (n&(n-1)) != 0 {
		fmt.Fprintln(out, "-1")
		return
	}

	// color[graph_vertex] = color
	color := make([]int, m)
	for i := 0; i < m; i++ { // i is a hypercube vertex
		c := 0
		for j := 0; j < n; j++ {
			if (i>>j)&1 > 0 {
				c ^= j
			}
		}
		if m > 0 && i < len(p) {
			color[p[i]] = c // p[i] is the graph vertex for hypercube vertex i
		}
	}
	
	for i := 0; i < m; i++ {
		fmt.Fprintf(out, "%d ", color[i])
	}
	fmt.Fprintln(out)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for i := 0; i < t; i++ {
		solve(in, out)
	}
}
