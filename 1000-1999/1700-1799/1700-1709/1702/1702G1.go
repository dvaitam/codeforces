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
	up    [][]int
	depth []int
)

func build(root int) {
	depth = make([]int, n+1)
	parent := make([]int, n+1)
	order := make([]int, 0, n)
	queue := make([]int, 0, n)
	queue = append(queue, root)
	parent[root] = root
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}
	up = make([][]int, LOG)
	up[0] = make([]int, n+1)
	for i := 1; i <= n; i++ {
		up[0][i] = parent[i]
	}
	for j := 1; j < LOG; j++ {
		up[j] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			up[j][i] = up[j-1][up[j-1][i]]
		}
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<uint(i)) != 0 {
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

func dist(u, v int) int {
	l := lca(u, v)
	return depth[u] + depth[v] - 2*depth[l]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	adj = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	build(1)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var k int
		fmt.Fscan(reader, &k)
		vertices := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &vertices[i])
		}
		if k <= 2 {
			fmt.Fprintln(writer, "YES")
			continue
		}
		a := vertices[0]
		b := a
		maxd := -1
		for _, v := range vertices {
			d := dist(a, v)
			if d > maxd {
				maxd = d
				b = v
			}
		}
		c := b
		maxd = -1
		for _, v := range vertices {
			d := dist(b, v)
			if d > maxd {
				maxd = d
				c = v
			}
		}
		good := true
		dTotal := dist(b, c)
		for _, v := range vertices {
			if dist(b, v)+dist(v, c) != dTotal {
				good = false
				break
			}
		}
		if good {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
