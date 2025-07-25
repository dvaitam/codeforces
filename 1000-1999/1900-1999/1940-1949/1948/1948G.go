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

	var n int
	var c int
	if _, err := fmt.Fscan(in, &n, &c); err != nil {
		return
	}
	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &w[i][j])
		}
	}

	// Prim's algorithm to build MST
	const inf = int(1e18)
	inMST := make([]bool, n)
	minEdge := make([]int, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		minEdge[i] = inf
		parent[i] = -1
	}
	minEdge[0] = 0
	mstWeight := 0
	adj := make([][]int, n)

	for it := 0; it < n; it++ {
		u := -1
		for i := 0; i < n; i++ {
			if !inMST[i] && (u == -1 || minEdge[i] < minEdge[u]) {
				u = i
			}
		}
		if u == -1 {
			// should not happen, graph connected
			return
		}
		inMST[u] = true
		if parent[u] != -1 {
			mstWeight += w[u][parent[u]]
			adj[u] = append(adj[u], parent[u])
			adj[parent[u]] = append(adj[parent[u]], u)
		}
		for v := 0; v < n; v++ {
			if !inMST[v] && w[u][v] > 0 && w[u][v] < minEdge[v] {
				minEdge[v] = w[u][v]
				parent[v] = u
			}
		}
	}

	// compute maximum matching of the MST
	type pair struct{ a, b int }
	var dfs func(int, int) (int, int)
	dfs = func(v, p int) (int, int) {
		children := []pair{}
		for _, to := range adj[v] {
			if to == p {
				continue
			}
			c0, c1 := dfs(to, v)
			children = append(children, pair{c0, c1})
		}
		// dp1: v matched with parent
		dp1 := 0
		for _, ch := range children {
			if ch.a > ch.b {
				dp1 += ch.a
			} else {
				dp1 += ch.b
			}
		}
		// dp0: v not matched with parent
		dp0 := dp1
		for _, ch := range children {
			cand := dp1 - max(ch.a, ch.b) + ch.a + 1
			if cand > dp0 {
				dp0 = cand
			}
		}
		return dp0, dp1
	}

	root := 0
	m0, _ := dfs(root, -1)
	matchingSize := m0
	totalCost := mstWeight + matchingSize*c
	fmt.Fprintln(out, totalCost)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
