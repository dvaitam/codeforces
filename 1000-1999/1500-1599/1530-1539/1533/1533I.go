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

const INF = int(1e9)

func (f *MCMF) MinCostFlow(s, t, maxf int) (int, int) {
	res := 0
	flow := 0
	for maxf > 0 {
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
		if f.dist[t] == INF {
			break
		}
		d := maxf
		for v := t; v != s; v = f.prevv[v] {
			if d > f.graph[f.prevv[v]][f.preve[v]].cap {
				d = f.graph[f.prevv[v]][f.preve[v]].cap
			}
		}
		maxf -= d
		flow += d
		for v := t; v != s; v = f.prevv[v] {
			e := &f.graph[f.prevv[v]][f.preve[v]]
			e.cap -= d
			rev := &f.graph[v][e.rev]
			rev.cap += d
		}
		res += d * f.dist[t]
	}
	return res, flow
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n1, n2, m int
	if _, err := fmt.Fscan(in, &n1, &n2, &m); err != nil {
		return
	}
	k := make([]int, n1)
	for i := 0; i < n1; i++ {
		fmt.Fscan(in, &k[i])
	}
	deg := make([]int, n1)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		edges[i] = [2]int{x, y}
		deg[x]++
	}

	total := 1 + n2 + n1 + 1
	source := 0
	sink := total - 1
	offsetL := 1 + n2

	flow := NewMCMF(total)
	for j := 0; j < n2; j++ {
		flow.AddEdge(source, 1+j, 1, 0)
	}
	for _, e := range edges {
		i := e[0]
		j := e[1]
		flow.AddEdge(1+j, offsetL+i, 1, 0)
	}
	for i := 0; i < n1; i++ {
		if deg[i] > 1 {
			flow.AddEdge(offsetL+i, sink, deg[i]-1, 0)
		}
		flow.AddEdge(offsetL+i, sink, 1, k[i])
	}

	cost, got := flow.MinCostFlow(source, sink, n2)
	if got < n2 {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, cost)
}
