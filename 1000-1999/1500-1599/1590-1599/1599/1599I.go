package main

import (
	"bufio"
	"fmt"
	"os"
)

// get function used in cactus verification
func get(parent []int, v int) int {
	if parent[v] == v {
		return v
	}
	parent[v] = get(parent, parent[v])
	return parent[v]
}

// isDesert checks whether the subgraph induced by edges is a desert.
// edges is a slice of pairs {u,v} (1-based vertices).
func isDesert(n int, edges [][2]int) bool {
	g := make([][]int, n+1)
	for _, e := range edges {
		u := e[0]
		v := e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	visited := make([]bool, n+1)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	dsu := make([]int, n+1) // for skipping used edges
	used := make([]bool, n+1)

	var dfs func(int) bool
	dfs = func(v int) bool {
		visited[v] = true
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			if !visited[to] {
				parent[to] = v
				depth[to] = depth[v] + 1
				dsu[to] = to
				if !dfs(to) {
					return false
				}
			} else if depth[to] < depth[v] {
				x := get(dsu, v)
				y := get(dsu, to)
				for x != y {
					if depth[x] < depth[y] {
						x, y = y, x
					}
					if used[x] {
						return false
					}
					used[x] = true
					dsu[x] = get(dsu, parent[x])
					x = get(dsu, x)
				}
			}
		}
		return true
	}

	for i := 1; i <= n; i++ {
		if !visited[i] {
			parent[i] = 0
			depth[i] = 0
			dsu[i] = i
			if !dfs(i) {
				return false
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	u := make([]int, m)
	v := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &u[i], &v[i])
	}

	edges := make([][2]int, 0, m)
	ans := 0
	for l := 0; l < m; l++ {
		edges = edges[:0]
		for r := l; r < m; r++ {
			edges = append(edges, [2]int{u[r], v[r]})
			if isDesert(n, edges) {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
