package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to   int
	rev  int
	cap  int
	cost int
}

type MCMF struct {
	n     int
	graph [][]edge
	dist  []int
	prevv []int
	preve []int
}

func NewMCMF(n int) *MCMF {
	g := make([][]edge, n)
	return &MCMF{n: n, graph: g, dist: make([]int, n), prevv: make([]int, n), preve: make([]int, n)}
}

func (f *MCMF) AddEdge(u, v, cap, cost int) {
	f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost})
	f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost})
}

const INF = int(1e18)

func (f *MCMF) MinCostMaxFlow(s, t int) int {
	res := 0
	for {
		for i := 0; i < f.n; i++ {
			f.dist[i] = INF
		}
		inq := make([]bool, f.n)
		q := make([]int, 0)
		f.dist[s] = 0
		inq[s] = true
		q = append(q, s)
		for idx := 0; idx < len(q); idx++ {
			v := q[idx]
			inq[v] = false
			for i, e := range f.graph[v] {
				if e.cap > 0 && f.dist[e.to] > f.dist[v]+e.cost {
					f.dist[e.to] = f.dist[v] + e.cost
					f.prevv[e.to] = v
					f.preve[e.to] = i
					if !inq[e.to] {
						q = append(q, e.to)
						inq[e.to] = true
					}
				}
			}
		}
		if f.dist[t] >= 0 || f.dist[t] == INF {
			break
		}
		d := 1
		for v := t; v != s; v = f.prevv[v] {
			if f.graph[f.prevv[v]][f.preve[v]].cap < d {
				d = f.graph[f.prevv[v]][f.preve[v]].cap
			}
		}
		for v := t; v != s; v = f.prevv[v] {
			e := &f.graph[f.prevv[v]][f.preve[v]]
			e.cap -= d
			rev := &f.graph[v][e.rev]
			rev.cap += d
		}
		res += d * f.dist[t]
	}
	return -res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	tot := n + m + 2
	source := 0
	sink := n + m + 1
	flow := NewMCMF(tot)

	for i := 1; i <= m; i++ {
		flow.AddEdge(source, i, 1, 0)
	}
	for j := 1; j <= n; j++ {
		flow.AddEdge(m+j, sink, 1, 0)
	}

	for i := 1; i <= m; i++ {
		var a, b, w int
		fmt.Fscan(reader, &a, &b, &w)
		flow.AddEdge(i, m+a, 1, -w)
		flow.AddEdge(i, m+b, 1, -w)
	}

	ans := flow.MinCostMaxFlow(source, sink)
	fmt.Fprintln(writer, ans)
}
