package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to  int
	rev int
	cap int64
}

var g [][]Edge

func addEdge(u, v int, c int64) {
	g[u] = append(g[u], Edge{v, len(g[v]), c})
	g[v] = append(g[v], Edge{u, len(g[u]) - 1, 0})
}

func bfs(s, t int, level []int) bool {
	for i := range level {
		level[i] = -1
	}
	q := []int{s}
	level[s] = 0
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if e.cap > 0 && level[e.to] < 0 {
				level[e.to] = level[v] + 1
				if e.to == t {
					return true
				}
				q = append(q, e.to)
			}
		}
	}
	return level[t] >= 0
}

func dfs(v, t int, f int64, level []int, it []int) int64 {
	if v == t {
		return f
	}
	for ; it[v] < len(g[v]); it[v]++ {
		e := &g[v][it[v]]
		if e.cap > 0 && level[v]+1 == level[e.to] {
			d := dfs(e.to, t, min(f, e.cap), level, it)
			if d > 0 {
				e.cap -= d
				rev := &g[e.to][e.rev]
				rev.cap += d
				return d
			}
		}
	}
	return 0
}

func maxFlow(s, t int) int64 {
	flow := int64(0)
	level := make([]int, len(g))
	it := make([]int, len(g))
	for bfs(s, t, level) {
		for i := range it {
			it[i] = 0
		}
		for {
			f := dfs(s, t, 1<<60, level, it)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	// Precompute divisors for 1..100
	divs := make([][]int, 101)
	for d := 1; d <= 100; d++ {
		for m := d; m <= 100; m += d {
			divs[m] = append(divs[m], d)
		}
	}

	s := 0
	t := n + 1
	g = make([][]Edge, n+2)
	const INF int64 = 1 << 60
	totalPos := int64(0)
	last := make([]int, 101)

	for i := 1; i <= n; i++ {
		if b[i] >= 0 {
			totalPos += int64(b[i])
			addEdge(s, i, int64(b[i]))
		} else {
			addEdge(i, t, int64(-b[i]))
		}
		for _, d := range divs[a[i]] {
			j := last[d]
			if j > 0 {
				addEdge(i, j, INF)
			}
		}
		last[a[i]] = i
	}

	flow := maxFlow(s, t)
	fmt.Println(totalPos - flow)
}
