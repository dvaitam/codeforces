package main

import (
	"bufio"
	"fmt"
	"os"
)

func dfs(v, p int, adj [][]int, parent, sub []int) {
	parent[v] = p
	sub[v] = 1
	for _, u := range adj[v] {
		if u == p {
			continue
		}
		dfs(u, v, adj, parent, sub)
		sub[v] += sub[u]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	parent := make([]int, n)
	sub := make([]int, n)
	dfs(0, -1, adj, parent, sub)

	possible := make([]bool, n)

	for v := 0; v < n; v++ {
		if len(adj[v]) < 2 {
			continue
		}
		comps := make([]int, 0, len(adj[v]))
		for _, u := range adj[v] {
			if u == parent[v] {
				comps = append(comps, n-sub[v])
			} else {
				comps = append(comps, sub[u])
			}
		}
		reachable := make([]bool, n)
		reachable[0] = true
		for _, size := range comps {
			for s := n - 1 - size; s >= 0; s-- {
				if reachable[s] {
					reachable[s+size] = true
				}
			}
		}
		for s := 1; s <= n-2; s++ {
			if reachable[s] {
				possible[s] = true
			}
		}
	}

	pairs := make([]int, 0)
	for s := 1; s <= n-2; s++ {
		if possible[s] {
			pairs = append(pairs, s)
		}
	}

	fmt.Fprintln(out, len(pairs))
	for _, a := range pairs {
		fmt.Fprintf(out, "%d %d\n", a, n-1-a)
	}
}
