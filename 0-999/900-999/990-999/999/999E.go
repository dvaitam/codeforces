package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s int
	if _, err := fmt.Fscan(in, &n, &m, &s); err != nil {
		return
	}
	g := make([][]int, n)
	rg := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		rg[v] = append(rg[v], u)
	}

	visited := make([]bool, n)
	order := make([]int, 0, n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		visited[v] = true
		for _, to := range g[v] {
			if !visited[to] {
				dfs1(to)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}

	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	var dfs2 func(int, int)
	dfs2 = func(v, c int) {
		comp[v] = c
		for _, to := range rg[v] {
			if comp[to] == -1 {
				dfs2(to, c)
			}
		}
	}
	cid := 0
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] == -1 {
			dfs2(v, cid)
			cid++
		}
	}

	indeg := make([]int, cid)
	for u := 0; u < n; u++ {
		for _, v := range g[u] {
			cu, cv := comp[u], comp[v]
			if cu != cv {
				indeg[cv]++
			}
		}
	}

	sComp := comp[s-1]
	ans := 0
	for c := 0; c < cid; c++ {
		if indeg[c] == 0 && c != sComp {
			ans++
		}
	}

	fmt.Fprintln(out, ans)
}
