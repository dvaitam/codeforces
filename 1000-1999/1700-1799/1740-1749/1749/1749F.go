package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

var (
	n     int
	adj   [][]int
	up    [LOG + 1][]int
	depth []int
	value []int64
)

func dfs(root int) {
	stack := []int{root}
	parent := make([]int, n+1)
	parent[root] = 0
	depth[root] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	up[0] = parent
	for k := 1; k <= LOG; k++ {
		up[k] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			up[k][i] = up[k-1][up[k-1][i]]
		}
	}
}

func lca(a, b int) int {
	if depth[a] < depth[b] {
		a, b = b, a
	}
	diff := depth[a] - depth[b]
	for k := LOG; k >= 0; k-- {
		if diff&(1<<uint(k)) != 0 {
			a = up[k][a]
		}
	}
	if a == b {
		return a
	}
	for k := LOG; k >= 0; k-- {
		if up[k][a] != up[k][b] {
			a = up[k][a]
			b = up[k][b]
		}
	}
	return up[0][a]
}

func getPath(u, v int) []int {
	w := lca(u, v)
	var path []int
	x := u
	for x != w {
		path = append(path, x)
		x = up[0][x]
	}
	path = append(path, w)
	var temp []int
	x = v
	for x != w {
		temp = append(temp, x)
		x = up[0][x]
	}
	for i := len(temp) - 1; i >= 0; i-- {
		path = append(path, temp[i])
	}
	return path
}

func update(u, v, k, d int) {
	path := getPath(u, v)
	visited := make([]bool, n+1)
	q := make([]int, 0, len(path))
	dist := make([]int, 0, len(path))
	for _, x := range path {
		if !visited[x] {
			visited[x] = true
			q = append(q, x)
			dist = append(dist, 0)
		}
	}
	head := 0
	for head < len(q) {
		node := q[head]
		dd := dist[head]
		head++
		value[node] += int64(k)
		if dd == d {
			continue
		}
		for _, to := range adj[node] {
			if !visited[to] {
				visited[to] = true
				q = append(q, to)
				dist = append(dist, dd+1)
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	depth = make([]int, n+1)
	for i := 0; i <= LOG; i++ {
		up[i] = make([]int, n+1)
	}
	value = make([]int64, n+1)
	dfs(1)

	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v int
			fmt.Fscan(reader, &v)
			fmt.Fprintln(writer, value[v])
		} else if t == 2 {
			var u, v, k, d int
			fmt.Fscan(reader, &u, &v, &k, &d)
			update(u, v, k, d)
		}
	}
}
