package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adj := make([][]int, n+1)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		edges[i] = [2]int{u, v}
	}

	// Step 1: peel leaves to identify cycle vertices
	deg := make([]int, n+1)
	for i := 1; i <= n; i++ {
		deg[i] = len(adj[i])
	}
	removed := make([]bool, n+1)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			q = append(q, i)
			removed[i] = true
		}
	}
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range adj[u] {
			if removed[v] {
				continue
			}
			deg[v]--
			if deg[v] == 1 {
				removed[v] = true
				q = append(q, v)
			}
		}
	}

	compID := make([]int, n+1)
	for i := range compID {
		compID[i] = -1
	}
	weight := make([]int, 1) // index 0 unused
	comp := 0

	visited := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if !removed[i] && !visited[i] {
			comp++
			stack := []int{i}
			visited[i] = true
			compID[i] = comp
			cnt := 0
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				cnt++
				for _, v := range adj[u] {
					if removed[v] {
						continue
					}
					if !visited[v] {
						visited[v] = true
						compID[v] = comp
						stack = append(stack, v)
					}
				}
			}
			weight = append(weight, cnt)
		}
	}

	for i := 1; i <= n; i++ {
		if compID[i] == -1 {
			comp++
			compID[i] = comp
			weight = append(weight, 0)
		}
	}

	tree := make([][]int, comp+1)
	for _, e := range edges {
		u := compID[e[0]]
		v := compID[e[1]]
		if u == v {
			continue
		}
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	// prepare LCA
	LOG := 0
	for (1 << LOG) <= comp {
		LOG++
	}
	up := make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, comp+1)
	}
	depth := make([]int, comp+1)
	prefix := make([]int, comp+1)
	bfs := []int{1}
	visited2 := make([]bool, comp+1)
	visited2[1] = true
	prefix[1] = weight[1]
	for len(bfs) > 0 {
		u := bfs[0]
		bfs = bfs[1:]
		for _, v := range tree[u] {
			if !visited2[v] {
				visited2[v] = true
				depth[v] = depth[u] + 1
				up[0][v] = u
				prefix[v] = prefix[u] + weight[v]
				bfs = append(bfs, v)
			}
		}
	}
	for k := 1; k < LOG; k++ {
		for v := 1; v <= comp; v++ {
			p := up[k-1][v]
			if p != 0 {
				up[k][v] = up[k-1][p]
			}
		}
	}

	var qn int
	fmt.Fscan(in, &qn)
	out := bufio.NewWriter(os.Stdout)
	for ; qn > 0; qn-- {
		var a, b int
		fmt.Fscan(in, &a, &b)
		u := compID[a]
		v := compID[b]
		l := lca(u, v, up, depth, LOG)
		ans := prefix[u] + prefix[v] - 2*prefix[l] + weight[l]
		fmt.Fprintln(out, ans)
	}
	out.Flush()
}

func lca(u, v int, up [][]int, depth []int, LOG int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			u = up[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}
