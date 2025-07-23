package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	to int
	w  int64
}

func score(w int64) int64 {
	if w <= 0 {
		return 0
	}
	k := int64(math.Floor((math.Sqrt(float64(8*w+1)) - 1) / 2))
	return (k+1)*w - k*(k+1)*(k+2)/6
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	g := make([][]Edge, n)
	rev := make([][]int, n)
	type E struct {
		u, v int
		w    int64
	}
	edges := make([]E, m)
	for i := 0; i < m; i++ {
		var x, y int
		var w int64
		fmt.Fscan(in, &x, &y, &w)
		x--
		y--
		g[x] = append(g[x], Edge{y, w})
		rev[y] = append(rev[y], x)
		edges[i] = E{x, y, w}
	}
	var s int
	fmt.Fscan(in, &s)
	s--
	// Kosaraju
	visited := make([]bool, n)
	order := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		stack := []struct{ v, idx int }{{i, 0}}
		visited[i] = true
		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			v := top.v
			if top.idx < len(g[v]) {
				e := g[v][top.idx]
				top.idx++
				if !visited[e.to] {
					visited[e.to] = true
					stack = append(stack, struct{ v, idx int }{e.to, 0})
				}
			} else {
				order = append(order, v)
				stack = stack[:len(stack)-1]
			}
		}
	}
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	cid := 0
	for oi := len(order) - 1; oi >= 0; oi-- {
		v := order[oi]
		if comp[v] != -1 {
			continue
		}
		stack := []int{v}
		comp[v] = cid
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, w := range rev[u] {
				if comp[w] == -1 {
					comp[w] = cid
					stack = append(stack, w)
				}
			}
		}
		cid++
	}
	C := cid
	internal := make([]int64, C)
	dag := make([][]Edge, C)
	indeg := make([]int, C)
	for _, e := range edges {
		cu := comp[e.u]
		cv := comp[e.v]
		if cu == cv {
			internal[cu] += score(e.w)
		} else {
			dag[cu] = append(dag[cu], Edge{cv, e.w})
			indeg[cv]++
		}
	}
	// topological order
	queue := make([]int, 0, C)
	for i := 0; i < C; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	order2 := make([]int, 0, C)
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		order2 = append(order2, u)
		for _, e := range dag[u] {
			indeg[e.to]--
			if indeg[e.to] == 0 {
				queue = append(queue, e.to)
			}
		}
	}
	dp := make([]int64, C)
	copy(dp, internal)
	for _, u := range order2 {
		for _, e := range dag[u] {
			cand := dp[u] + e.w + internal[e.to]
			if cand > dp[e.to] {
				dp[e.to] = cand
			}
		}
	}
	fmt.Println(dp[comp[s]])
}
