package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type edge struct {
	to   int
	rev  int
	cap  int
	cost int64
}

type graph [][]edge

func (g *graph) addEdge(from, to, cap int, cost int64) {
	e := edge{to: to, rev: len((*g)[to]), cap: cap, cost: cost}
	r := edge{to: from, rev: len((*g)[from]), cap: 0, cost: -cost}
	(*g)[from] = append((*g)[from], e)
	(*g)[to] = append((*g)[to], r)
}

type item struct {
	d int64
	v int
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].d < p[j].d }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}

func minCostFlow(g graph, s, t, maxF int) int64 {
	const inf int64 = 1 << 60
	n := len(g)
	h := make([]int64, n)
	prevv := make([]int, n)
	preve := make([]int, n)
	res := int64(0)
	for f := 0; f < maxF; {
		dist := make([]int64, n)
		for i := range dist {
			dist[i] = inf
		}
		dist[s] = 0
		pq := &pq{}
		heap.Push(pq, item{0, s})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(item)
			v := it.v
			if dist[v] < it.d {
				continue
			}
			for i, e := range g[v] {
				if e.cap > 0 {
					nd := dist[v] + e.cost + h[v] - h[e.to]
					if nd < dist[e.to] {
						dist[e.to] = nd
						prevv[e.to] = v
						preve[e.to] = i
						heap.Push(pq, item{nd, e.to})
					}
				}
			}
		}
		if dist[t] == inf {
			return -1
		}
		for i := 0; i < n; i++ {
			if dist[i] < inf {
				h[i] += dist[i]
			}
		}
		d := maxF - f
		for v := t; v != s; v = prevv[v] {
			if g[prevv[v]][preve[v]].cap < d {
				d = g[prevv[v]][preve[v]].cap
			}
		}
		f += d
		res += int64(d) * h[t]
		for v := t; v != s; v = prevv[v] {
			e := &g[prevv[v]][preve[v]]
			e.cap -= d
			rev := &g[v][e.rev]
			rev.cap += d
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	s := n + 1
	t := n + 2
	g := make(graph, n+3)
	for i := 0; i < n; i++ {
		g.addEdge(i, i+1, k, 0)
		g.addEdge(s, i+1, 1, a[i])
		g.addEdge(i+1, t, 1, b[i])
	}
	ans := minCostFlow(g, s, t, k)
	fmt.Println(ans)
}
