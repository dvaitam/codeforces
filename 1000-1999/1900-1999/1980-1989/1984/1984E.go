package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ to int }

var (
	g      [][]int
	dp0    []int
	dp1    []int
	parent []int
)

func dfs(u, p int) {
	parent[u] = p
	for _, v := range g[u] {
		if v == p {
			continue
		}
		dfs(v, u)
	}
	base := 0
	for _, v := range g[u] {
		if v == p {
			continue
		}
		if dp0[v] > dp1[v] {
			base += dp0[v]
		} else {
			base += dp1[v]
		}
	}
	best := base
	for _, v := range g[u] {
		if v == p {
			continue
		}
		cand := 1 + dp1[v] + base
		if dp0[v] > dp1[v] {
			cand -= dp0[v]
		} else {
			cand -= dp1[v]
		}
		if cand > best {
			best = cand
		}
	}
	dp0[u] = best
	sum := 0
	for _, v := range g[u] {
		if v == p {
			continue
		}
		sum += dp0[v]
	}
	dp1[u] = sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		g = make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		dp0 = make([]int, n)
		dp1 = make([]int, n)
		parent = make([]int, n)
		dfs(0, -1)
		match := dp0[0]
		if dp1[0] > match {
			match = dp1[0]
		}
		ans := n - match
		fmt.Fprintln(out, ans)
	}
}
