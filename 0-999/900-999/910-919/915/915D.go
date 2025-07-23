package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct{ u, v int }

var (
	n, m int
	adj  [][]int
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &m)
	adj = make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
	}
	if check(-1, -1) {
		fmt.Println("YES")
		return
	}
	cyc, ok := findCycle()
	if !ok {
		fmt.Println("NO")
		return
	}
	for _, e := range cyc {
		if check(e.u, e.v) {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}

func check(skipU, skipV int) bool {
	indeg := make([]int, n+1)
	for u := 1; u <= n; u++ {
		for _, v := range adj[u] {
			if u == skipU && v == skipV {
				continue
			}
			indeg[v]++
		}
	}
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	count := 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		count++
		for _, to := range adj[v] {
			if v == skipU && to == skipV {
				continue
			}
			indeg[to]--
			if indeg[to] == 0 {
				q = append(q, to)
			}
		}
	}
	return count == n
}

func findCycle() ([]edge, bool) {
	color := make([]int, n+1)
	stack := make([]int, 0, n)
	pos := make([]int, n+1)
	var cycle []edge
	var found bool
	var dfs func(int)
	dfs = func(v int) {
		color[v] = 1
		pos[v] = len(stack)
		stack = append(stack, v)
		for _, to := range adj[v] {
			if found {
				return
			}
			if color[to] == 0 {
				dfs(to)
			} else if color[to] == 1 {
				found = true
				idx := pos[to]
				nodes := append([]int{}, stack[idx:]...)
				nodes = append(nodes, to)
				for i := 0; i < len(nodes)-1; i++ {
					cycle = append(cycle, edge{nodes[i], nodes[i+1]})
				}
				return
			}
		}
		stack = stack[:len(stack)-1]
		color[v] = 2
	}
	for i := 1; i <= n && !found; i++ {
		if color[i] == 0 {
			dfs(i)
		}
	}
	return cycle, found
}
