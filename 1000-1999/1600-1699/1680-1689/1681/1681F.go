package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	val int
}

var (
	n       int
	adj     [][]Edge
	parent  []int
	depth   []int
	edgeVal []int
)

func dfs(v, p int, val int) {
	parent[v] = p
	edgeVal[v] = val
	for _, e := range adj[v] {
		if e.to != p {
			depth[e.to] = depth[v] + 1
			dfs(e.to, v, e.val)
		}
	}
}

func pathValues(u, v int) []int {
	var vals []int
	for depth[u] > depth[v] {
		vals = append(vals, edgeVal[u])
		u = parent[u]
	}
	for depth[v] > depth[u] {
		vals = append(vals, edgeVal[v])
		v = parent[v]
	}
	for u != v {
		vals = append(vals, edgeVal[u])
		vals = append(vals, edgeVal[v])
		u = parent[u]
		v = parent[v]
	}
	return vals
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	adj = make([][]Edge, n+1)
	for i := 0; i < n-1; i++ {
		var v, u, x int
		fmt.Fscan(reader, &v, &u, &x)
		adj[v] = append(adj[v], Edge{u, x})
		adj[u] = append(adj[u], Edge{v, x})
	}

	parent = make([]int, n+1)
	depth = make([]int, n+1)
	edgeVal = make([]int, n+1)
	dfs(1, 0, 0)

	var ans int64
	for v := 1; v <= n; v++ {
		for u := v + 1; u <= n; u++ {
			vals := pathValues(v, u)
			freq := make(map[int]int)
			for _, val := range vals {
				freq[val]++
			}
			for _, f := range freq {
				if f == 1 {
					ans++
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
