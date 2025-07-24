package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	ch byte
}

// compute path string from node x to node y using BFS in the tree
func pathString(x, y int, g [][]Edge) string {
	n := len(g)
	parent := make([]int, n)
	edgeChar := make([]byte, n)
	for i := range parent {
		parent[i] = -1
	}
	queue := []int{x}
	parent[x] = x
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		if u == y {
			break
		}
		for _, e := range g[u] {
			if parent[e.to] == -1 {
				parent[e.to] = u
				edgeChar[e.to] = e.ch
				queue = append(queue, e.to)
			}
		}
	}
	var path []byte
	v := y
	for v != x {
		path = append(path, edgeChar[v])
		v = parent[v]
	}
	// reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return string(path)
}

func query(x, y int, g [][]Edge) int {
	target := pathString(x, y, g)
	n := len(g)
	visited := make([]bool, n)
	strs := make([]string, n)
	queue := []int{x}
	visited[x] = true
	strs[x] = ""
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, e := range g[u] {
			if !visited[e.to] {
				visited[e.to] = true
				strs[e.to] = strs[u] + string(e.ch)
				queue = append(queue, e.to)
			}
		}
	}
	cnt := 0
	for z := 0; z < n; z++ {
		if z == x || z == y {
			continue
		}
		if target > strs[z] {
			cnt++
		}
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var v, u int
		var s string
		fmt.Fscan(in, &v, &u, &s)
		v--
		u--
		ch := s[0]
		g[v] = append(g[v], Edge{u, ch})
		g[u] = append(g[u], Edge{v, ch})
	}

	for i := 0; i < q; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		ans := query(x, y, g)
		fmt.Fprintln(out, ans)
	}
}
