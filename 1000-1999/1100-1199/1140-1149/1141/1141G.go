package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	g := make([][]int, n+1)
	edges := make([][2]int, 0, n-1)
	deg := make([]int, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &v, &u)
		g[v] = append(g[v], u)
		g[u] = append(g[u], v)
		deg[v]++
		deg[u]++
		edges = append(edges, [2]int{u, v})
	}
	// compute x
	a := make([]int, n)
	for i := 1; i <= n; i++ {
		a[i-1] = deg[i]
	}
	sort.Ints(a)
	x := 1
	pIdx := n - 1
	for i := n - 1; i >= 1; i-- {
		for pIdx >= 0 && a[pIdx] > i {
			pIdx--
		}
		if (n - 1 - pIdx) <= k {
			x = i
		}
	}
	fmt.Fprintln(out, x)
	// prepare for dfs
	dis := make([]int, n+1)
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		c[i] = -1
	}
	// iterative dfs
	type frame struct{ v, parent, idx, kVal int }
	stack := make([]frame, 0, n)
	stack = append(stack, frame{v: 1, parent: 0, idx: 0, kVal: 0})
	for len(stack) > 0 {
		f := &stack[len(stack)-1]
		v := f.v
		if f.idx >= len(g[v]) {
			stack = stack[:len(stack)-1]
			continue
		}
		u := g[v][f.idx]
		f.idx++
		if u == f.parent {
			continue
		}
		dis[u] = dis[v] + 1
		if len(g[v]) <= x && f.kVal == c[v] {
			f.kVal = (f.kVal + 1) % x
		}
		c[u] = f.kVal
		f.kVal = (f.kVal + 1) % x
		stack = append(stack, frame{v: u, parent: v, idx: 0, kVal: 0})
	}
	// output colors for edges in input order
	for _, e := range edges {
		u, v := e[0], e[1]
		if dis[u] > dis[v] {
			u, v = v, u
		}
		// color of deeper node v
		fmt.Fprintf(out, "%d ", c[v]+1)
	}
	fmt.Fprintln(out)
}
