package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = int(1e9)

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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, m)
	for j := 0; j < m; j++ {
		fmt.Fscan(in, &B[j])
	}
	sumA, sumB := 0, 0
	for _, v := range A {
		sumA += v
	}
	for _, v := range B {
		sumB += v
	}
	if sumA != sumB {
		fmt.Fprintln(out, -1)
		return
	}
	total1 := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if a[i][j] == 1 {
				total1++
			}
		}
	}
	s := n + m
	t := s + 1
	flow := NewMCMF(t + 1)
	for i := 0; i < n; i++ {
		flow.AddEdge(s, i, A[i], 0)
	}
	for j := 0; j < m; j++ {
		flow.AddEdge(n+j, t, B[j], 0)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cost := 1
			if a[i][j] == 1 {
				cost = -1
			}
			flow.AddEdge(i, n+j, 1, cost)
		}
	}
	cost, got := flow.MinCostFlow(s, t, sumA)
	if got < sumA {
		fmt.Fprintln(out, -1)
		return
	}
	ans := total1 + cost
	fmt.Fprintln(out, ans)
}
