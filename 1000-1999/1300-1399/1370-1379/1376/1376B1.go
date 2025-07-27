package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
		deg[a]++
		deg[b]++
	}

	type pair struct{ v, d int }
	vertices := make([]pair, n)
	for i := 1; i <= n; i++ {
		vertices[i-1] = pair{i, deg[i]}
	}
	sort.Slice(vertices, func(i, j int) bool {
		if vertices[i].d == vertices[j].d {
			return vertices[i].v < vertices[j].v
		}
		return vertices[i].d < vertices[j].d
	})

	used := make([]bool, n+1)
	selected := make([]int, n+1)
	count := 0
	for _, p := range vertices {
		v := p.v
		if !used[v] {
			selected[v] = 1
			count++
			used[v] = true
			for _, u := range adj[v] {
				used[u] = true
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, count)
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, selected[i])
	}
	fmt.Fprintln(out)
}
