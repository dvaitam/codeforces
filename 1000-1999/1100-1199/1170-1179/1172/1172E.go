package main

import (
	"bufio"
	"fmt"
	"os"
)

func getPath(u, v int, adj [][]int) []int {
	if u == v {
		return []int{u}
	}
	n := len(adj) - 1
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -1
	}
	q := []int{u}
	parent[u] = 0
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			break
		}
		for _, nb := range adj[x] {
			if parent[nb] == -1 {
				parent[nb] = x
				q = append(q, nb)
			}
		}
	}
	path := []int{}
	cur := v
	for cur != u {
		path = append(path, cur)
		cur = parent[cur]
	}
	path = append(path, u)
	return path
}

func sumDistinct(n int, adj [][]int, color []int) int64 {
	total := int64(0)
	seen := make(map[int]bool)
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			p := getPath(i, j, adj)
			for k := range seen {
				delete(seen, k)
			}
			for _, node := range p {
				seen[color[node]] = true
			}
			total += int64(len(seen))
		}
	}
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &color[i])
	}
	ans := make([]int64, 0, m+1)
	ans = append(ans, sumDistinct(n, adj, color))
	for k := 0; k < m; k++ {
		var u, x int
		fmt.Fscan(in, &u, &x)
		color[u] = x
		ans = append(ans, sumDistinct(n, adj, color))
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
}
