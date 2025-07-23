package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOGD = 17 // since n <= 1e5

var (
	adj   [][]int
	up    [][]int
	depth []int
)

func dfs(v, p int) {
	up[0][v] = p
	for i := 1; i < LOGD; i++ {
		up[i][v] = up[i-1][up[i-1][v]]
	}
	for _, to := range adj[v] {
		if to == p {
			continue
		}
		depth[to] = depth[v] + 1
		dfs(to, v)
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := 0; diff > 0; i++ {
		if diff&1 == 1 {
			u = up[i][u]
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for i := LOGD - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}

func dist(u, v int) int {
	w := lca(u, v)
	return depth[u] + depth[v] - 2*depth[w]
}

func interLen(s, f, t int) int {
	d1 := dist(s, f)
	d2 := dist(t, f)
	d3 := dist(s, t)
	return (d1+d2-d3)/2 + 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	adj = make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(in, &p)
		adj[p] = append(adj[p], i)
		adj[i] = append(adj[i], p)
	}
	up = make([][]int, LOGD)
	for i := 0; i < LOGD; i++ {
		up[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	dfs(1, 1)

	perms := [][]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for ; q > 0; q-- {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		arr := []int{a, b, c}
		best := 0
		for _, pm := range perms {
			s := arr[pm[0]]
			f := arr[pm[1]]
			t := arr[pm[2]]
			val := interLen(s, f, t)
			if val > best {
				best = val
			}
		}
		fmt.Fprintln(out, best)
	}
}
