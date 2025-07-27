package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
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

func (f *MCMF) MinCostFlow(s, t, maxf int) int {
	res := 0
	flow := 0
	for flow < maxf {
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
		d := maxf - flow
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
		flow += d
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	m := 1 << n
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &a[i])
	}

	type node struct {
		idx int
		val int
	}
	nodes := make([]node, m)
	for i := 0; i < m; i++ {
		nodes[i] = node{i, a[i]}
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].val > nodes[j].val })
	tsize := 4 * k
	if tsize > m {
		tsize = m
	}
	selected := make([]bool, m)
	for i := 0; i < tsize; i++ {
		selected[nodes[i].idx] = true
	}
	for i := 0; i < tsize; i++ {
		v := nodes[i].idx
		for b := 0; b < n; b++ {
			selected[v^(1<<b)] = true
		}
	}

	leftNodes := make([]int, 0)
	rightNodes := make([]int, 0)
	idxL := make(map[int]int)
	idxR := make(map[int]int)
	for i := 0; i < m; i++ {
		if selected[i] {
			if bits.OnesCount(uint(i))%2 == 0 {
				idxL[i] = len(leftNodes)
				leftNodes = append(leftNodes, i)
			} else {
				idxR[i] = len(rightNodes)
				rightNodes = append(rightNodes, i)
			}
		}
	}

	total := 1 + len(leftNodes) + len(rightNodes) + 1
	source := 0
	sink := total - 1
	offsetR := 1 + len(leftNodes)
	flow := NewMCMF(total)

	for i := 0; i < len(leftNodes); i++ {
		flow.AddEdge(source, 1+i, 1, 0)
	}
	for i := 0; i < len(rightNodes); i++ {
		flow.AddEdge(offsetR+i, sink, 1, 0)
	}

	for _, u := range leftNodes {
		lu := 1 + idxL[u]
		for b := 0; b < n; b++ {
			v := u ^ (1 << b)
			ri, ok := idxR[v]
			if ok {
				w := -(a[u] + a[v])
				flow.AddEdge(lu, offsetR+ri, 1, w)
			}
		}
	}

	cost := flow.MinCostFlow(source, sink, k)
	fmt.Fprintln(out, -cost)
}
