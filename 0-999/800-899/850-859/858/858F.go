package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
}

type Triple struct {
	x, y, z int
}

var (
	g       [][]Edge
	edges   [][2]int
	used    []bool
	visited []bool
	stack   [][]int
	res     []Triple
)

func other(id, v int) int {
	if edges[id][0] == v {
		return edges[id][1]
	}
	return edges[id][0]
}

func dfs(v int, pe int) {
	visited[v] = true
	for _, e := range g[v] {
		if e.id == pe {
			continue
		}
		if used[e.id] {
			continue
		}
		if !visited[e.to] {
			dfs(e.to, e.id)
		}
		if used[e.id] {
			continue
		}
		stack[v] = append(stack[v], e.id)
	}
	for len(stack[v]) >= 2 {
		e1 := stack[v][len(stack[v])-1]
		stack[v] = stack[v][:len(stack[v])-1]
		e2 := stack[v][len(stack[v])-1]
		stack[v] = stack[v][:len(stack[v])-1]
		res = append(res, Triple{other(e1, v), v, other(e2, v)})
		used[e1] = true
		used[e2] = true
	}
	if len(stack[v]) == 1 && pe != -1 {
		e1 := stack[v][len(stack[v])-1]
		stack[v] = stack[v][:len(stack[v])-1]
		res = append(res, Triple{other(e1, v), v, other(pe, v)})
		used[e1] = true
		used[pe] = true
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	edges = make([][2]int, m)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		edges[i] = [2]int{a, b}
		g[a] = append(g[a], Edge{b, i})
		g[b] = append(g[b], Edge{a, i})
	}
	used = make([]bool, m)
	visited = make([]bool, n+1)
	stack = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if !visited[i] {
			dfs(i, -1)
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(res))
	for _, t := range res {
		fmt.Fprintf(out, "%d %d %d\n", t.x, t.y, t.z)
	}
	out.Flush()
}
