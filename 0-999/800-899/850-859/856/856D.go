package main

// Attempted solution for problem 856D - see problemD.txt for the statement.
// The goal is to choose a subset of additional edges so that each vertex
// belongs to at most one added cycle (the resulting graph is a cactus)
// and the sum of beauties of chosen edges is maximized.
//
// This implementation uses a greedy approach: edges are processed in
// descending order of beauty and added if their corresponding cycle does
// not share any vertex with previously selected cycles.  To speed up
// checking of vertex availability along paths we use a disjoint-set union
// structure that skips vertices already included in some cycle.
// The algorithm is not guaranteed to be optimal but demonstrates a
// reasonable attempt within given constraints.

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const logN = 20

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	parent := make([][]int, logN)
	for i := range parent {
		parent[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	g := make([][]int, n+1)

	for v := 2; v <= n; v++ {
		var p int
		fmt.Fscan(in, &p)
		parent[0][v] = p
		g[p] = append(g[p], v)
		g[v] = append(g[v], p)
	}

	var dfs func(int, int)
	dfs = func(v, p int) {
		for _, to := range g[v] {
			if to == p {
				continue
			}
			depth[to] = depth[v] + 1
			parent[0][to] = v
			dfs(to, v)
		}
	}
	dfs(1, 0)

	for k := 1; k < logN; k++ {
		for i := 1; i <= n; i++ {
			parent[k][i] = parent[k-1][parent[k-1][i]]
		}
	}

	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		for k := logN - 1; k >= 0; k-- {
			if depth[a]-depth[b] >= 1<<uint(k) {
				a = parent[k][a]
			}
		}
		if a == b {
			return a
		}
		for k := logN - 1; k >= 0; k-- {
			if parent[k][a] != parent[k][b] {
				a = parent[k][a]
				b = parent[k][b]
			}
		}
		return parent[0][a]
	}

	type edge struct{ u, v, w int }
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w > edges[j].w })

	// disjoint-set union used to skip already used vertices
	dsu := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dsu[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if dsu[x] != x {
			dsu[x] = find(dsu[x])
		}
		return dsu[x]
	}
	used := make([]bool, n+1)

	var pathCollect func(int, int) ([]int, bool)
	pathCollect = func(x, anc int) ([]int, bool) {
		res := []int{}
		for {
			x = find(x)
			if depth[x] <= depth[anc] {
				break
			}
			res = append(res, x)
			x = parent[0][x]
		}
		if x != anc {
			return nil, false
		}
		return res, true
	}

	var result int
	for _, e := range edges {
		l := lca(e.u, e.v)
		left, ok1 := pathCollect(e.u, l)
		if !ok1 {
			continue
		}
		right, ok2 := pathCollect(e.v, l)
		if !ok2 {
			continue
		}
		if used[l] {
			continue
		}
		// accept the edge
		result += e.w
		used[l] = true
		if l != 1 {
			dsu[l] = parent[0][l]
		}
		for _, x := range left {
			used[x] = true
			dsu[x] = parent[0][x]
		}
		for _, x := range right {
			used[x] = true
			dsu[x] = parent[0][x]
		}
	}

	fmt.Fprintln(out, result)
}
