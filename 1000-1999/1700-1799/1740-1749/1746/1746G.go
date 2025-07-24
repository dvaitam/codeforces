package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

type edge struct {
	to   int
	rev  int
	cap  int
	cost int64
}

type MCMF struct {
	n     int
	graph [][]edge
	dist  []int64
	prevv []int
	preve []int
}

func NewMCMF(n int) *MCMF {
	g := make([][]edge, n)
	return &MCMF{n: n, graph: g, dist: make([]int64, n), prevv: make([]int, n), preve: make([]int, n)}
}

func (f *MCMF) AddEdge(u, v, cap int, cost int64) {
	f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost})
	f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost})
}

func (f *MCMF) MinCostFlow(s, t, maxf int) (int64, int) {
	res := int64(0)
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
		res += int64(d) * f.dist[t]
	}
	return res, flow
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, a, b, c int
		fmt.Fscan(in, &n, &a, &b, &c)
		tasks := make([]struct{ r, typ, d int }, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &tasks[i].r, &tasks[i].typ, &tasks[i].d)
		}
		m := a + b + c
		if m > n {
			fmt.Fprintln(out, -1)
			continue
		}
		taskStart := 4
		dayStart := taskStart + n
		sink := dayStart + n + 1
		total := sink + 1

		flow := NewMCMF(total)
		flow.AddEdge(0, 1, a, 0)
		flow.AddEdge(0, 2, b, 0)
		flow.AddEdge(0, 3, c, 0)

		for i, t := range tasks {
			idx := taskStart + i
			switch t.typ {
			case 1:
				flow.AddEdge(1, idx, 1, 0)
			case 2:
				flow.AddEdge(2, idx, 1, 0)
			default:
				flow.AddEdge(3, idx, 1, 0)
			}
			if t.d > n {
				t.d = n
			}
			flow.AddEdge(idx, dayStart+t.d, 1, -int64(t.r))
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				flow.AddEdge(dayStart+i, dayStart+i-1, m, 0)
			}
			flow.AddEdge(dayStart+i, sink, 1, 0)
		}

		cost, got := flow.MinCostFlow(0, sink, m)
		if got < m {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, -cost)
		}
	}
}
