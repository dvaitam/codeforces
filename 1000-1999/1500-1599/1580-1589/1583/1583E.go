package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

func lca(u, v int, depth []int, up [][]int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := LOG - 1; k >= 0; k-- {
		if diff&(1<<k) != 0 {
			u = up[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := LOG - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}

func getPath(a, b int, parent []int, depth []int, up [][]int) []int {
	l := lca(a, b, depth, up)
	var path []int
	u := a
	for u != l {
		path = append(path, u)
		u = parent[u]
	}
	path = append(path, l)
	var tail []int
	v := b
	for v != l {
		tail = append(tail, v)
		v = parent[v]
	}
	for i := len(tail) - 1; i >= 0; i-- {
		path = append(path, tail[i])
	}
	return path
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}

	parent := make([]int, n+1)
	depth := make([]int, n+1)
	visited := make([]bool, n+1)
	queue := make([]int, 0, n)
	root := 1
	queue = append(queue, root)
	visited[root] = true
	parent[root] = 0
	for idx := 0; idx < len(queue); idx++ {
		v := queue[idx]
		for _, to := range g[v] {
			if !visited[to] {
				visited[to] = true
				parent[to] = v
				depth[to] = depth[v] + 1
				queue = append(queue, to)
			}
		}
	}

	up := make([][]int, LOG)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	for v := 1; v <= n; v++ {
		up[0][v] = parent[v]
	}
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			if up[k-1][v] != 0 {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	queries := make([][2]int, q)
	parity := make([]int, n+1)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		queries[i] = [2]int{a, b}
		parity[a] ^= 1
		parity[b] ^= 1
	}

	odd := 0
	for i := 1; i <= n; i++ {
		if parity[i] == 1 {
			odd++
		}
	}
	if odd > 0 {
		fmt.Fprintln(writer, "NO")
		fmt.Fprintln(writer, odd/2)
		return
	}

	fmt.Fprintln(writer, "YES")
	for _, p := range queries {
		a, b := p[0], p[1]
		path := getPath(a, b, parent, depth, up)
		fmt.Fprintln(writer, len(path))
		for i, node := range path {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, node)
		}
		fmt.Fprintln(writer)
	}
}
