package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to, rev, cap int
	cost         int
}

type Graph [][]*Edge

func (g Graph) AddEdge(u, v, cap, cost int) {
	g[u] = append(g[u], &Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
	g[v] = append(g[v], &Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

type Item struct {
	v    int
	dist int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func minCostFlow(g Graph, s, t, maxf int) (int, int) {
	n := len(g)
	h := make([]int, n)
	prevv := make([]int, n)
	preve := make([]int, n)
	flow, cost := 0, 0
	const INF = int(1 << 60)
	for flow < maxf {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = INF
		}
		dist[s] = 0
		pq := &PriorityQueue{}
		heap.Init(pq)
		heap.Push(pq, Item{v: s, dist: 0})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			v := it.v
			if dist[v] < it.dist {
				continue
			}
			for i, e := range g[v] {
				if e.cap > 0 && dist[e.to] > dist[v]+e.cost+h[v]-h[e.to] {
					dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
					prevv[e.to] = v
					preve[e.to] = i
					heap.Push(pq, Item{v: e.to, dist: dist[e.to]})
				}
			}
		}
		if dist[t] == INF {
			break
		}
		for v := 0; v < n; v++ {
			if dist[v] < INF {
				h[v] += dist[v]
			}
		}
		d := maxf - flow
		for v := t; v != s; v = prevv[v] {
			if d > g[prevv[v]][preve[v]].cap {
				d = g[prevv[v]][preve[v]].cap
			}
		}
		flow += d
		cost += d * h[t]
		for v := t; v != s; v = prevv[v] {
			e := g[prevv[v]][preve[v]]
			e.cap -= d
			g[v][e.rev].cap += d
		}
	}
	return flow, cost
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	lb := make([]int, n)
	ub := make([]int, n)
	for i := 0; i < n; i++ {
		lb[i] = 1
		ub[i] = n
	}
	for i := 0; i < q; i++ {
		var t, l, r, v int
		fmt.Fscan(reader, &t, &l, &r, &v)
		l--
		r--
		if t == 1 {
			for j := l; j <= r; j++ {
				if lb[j] < v {
					lb[j] = v
				}
			}
		} else {
			for j := l; j <= r; j++ {
				if ub[j] > v {
					ub[j] = v
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		if lb[i] > ub[i] {
			fmt.Println(-1)
			return
		}
	}

	V := 2*n + 2
	s := 0
	posBase := 1
	valBase := posBase + n
	tIdx := valBase + n
	g := make(Graph, V)

	for i := 0; i < n; i++ {
		g.AddEdge(s, posBase+i, 1, 0)
		for v := lb[i]; v <= ub[i]; v++ {
			g.AddEdge(posBase+i, valBase+(v-1), 1, 0)
		}
	}
	for v := 0; v < n; v++ {
		for k := 1; k <= n; k++ {
			cost := 2*k - 1
			g.AddEdge(valBase+v, tIdx, 1, cost)
		}
	}

	flow, cost := minCostFlow(g, s, tIdx, n)
	if flow < n {
		fmt.Println(-1)
	} else {
		fmt.Println(cost)
	}
}
