package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Minimum cost flow with potentials (successive shortest paths)
type edge struct {
	to, rev int
	cap     int64
	cost    int64
}

type mcmf struct {
	n     int
	graph [][]edge
	h     []int64
	dist  []int64
	prevV []int
	prevE []int
}

func newMCMF(n int) *mcmf {
	return &mcmf{
		n:     n,
		graph: make([][]edge, n),
		h:     make([]int64, n),
		dist:  make([]int64, n),
		prevV: make([]int, n),
		prevE: make([]int, n),
	}
}

func (m *mcmf) addEdge(from, to int, cap, cost int64) {
	m.graph[from] = append(m.graph[from], edge{to, len(m.graph[to]), cap, cost})
	m.graph[to] = append(m.graph[to], edge{from, len(m.graph[from]) - 1, 0, -cost})
}

// minCostFlow returns (flow, cost)
func (m *mcmf) minCostFlow(s, t int, maxf int64) (int64, int64) {
	n := m.n
	const inf = int64(4e18)
	// initialize potentials
	for i := 0; i < n; i++ {
		m.h[i] = 0
	}
	var flow, cost int64
	for flow < maxf {
		// dijkstra
		for i := 0; i < n; i++ {
			m.dist[i] = inf
		}
		m.dist[s] = 0
		// priority queue
		pq := &priorityQueue{}
		heap.Init(pq)
		heap.Push(pq, &item{v: s, dist: 0})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(*item)
			v, d := it.v, it.dist
			if d > m.dist[v] {
				continue
			}
			for ei, e := range m.graph[v] {
				if e.cap > 0 {
					nd := d + e.cost + m.h[v] - m.h[e.to]
					if nd < m.dist[e.to] {
						m.dist[e.to] = nd
						m.prevV[e.to] = v
						m.prevE[e.to] = ei
						heap.Push(pq, &item{e.to, nd})
					}
				}
			}
		}
		if m.dist[t] == inf {
			break
		}
		for v := 0; v < n; v++ {
			if m.dist[v] < inf {
				m.h[v] += m.dist[v]
			}
		}
		// add as much as possible
		d := maxf - flow
		// find minimum residual capacity
		for v := t; v != s; v = m.prevV[v] {
			e := m.graph[m.prevV[v]][m.prevE[v]]
			if e.cap < d {
				d = e.cap
			}
		}
		flow += d
		cost += d * m.h[t]
		// update residual capacities
		for v := t; v != s; v = m.prevV[v] {
			e := &m.graph[m.prevV[v]][m.prevE[v]]
			e.cap -= d
			m.graph[v][e.rev].cap += d
		}
	}
	return flow, cost
}

// priority queue for Dijkstra
type item struct {
	v    int
	dist int64
}

type priorityQueue []*item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[0 : n-1]
	return it
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	u := make([]int, m)
	v := make([]int, m)
	c := make([]int64, m)
	f := make([]int64, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &u[i], &v[i], &c[i], &f[i])
		u[i]--
		v[i]--
	}
	var res int64
	// adjust capacities for f > c
	for i := 0; i < m; i++ {
		if f[i] > c[i] {
			res += f[i] - c[i]
			c[i] = f[i]
		}
	}
	// compute imbalance b[v] = sum in f - sum out f
	b := make([]int64, n)
	for i := 0; i < m; i++ {
		b[u[i]] -= f[i]
		b[v[i]] += f[i]
	}
	// total demand for internal nodes
	var totalDemand int64
	for i := 1; i < n-1; i++ {
		if b[i] > 0 {
			totalDemand += b[i]
		}
	}
	if totalDemand > 0 {
		// build mcmf
		SS := n
		TT := n + 1
		mvc := newMCMF(n + 2)
		// edges for flow change y
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if f[i] > 0 {
				mvc.addEdge(vi, ui, f[i], 1)
			}
			capInc := c[i] - f[i]
			if capInc > 0 {
				mvc.addEdge(ui, vi, capInc, 1)
			}
			// allow extra increase at cost 2
			mvc.addEdge(ui, vi, totalDemand, 2)
		}
		// demands for internal nodes
		for i := 1; i < n-1; i++ {
			if b[i] > 0 {
				mvc.addEdge(i, TT, b[i], 0)
			} else if b[i] < 0 {
				mvc.addEdge(SS, i, -b[i], 0)
			}
		}
		// compute flow
		_, cost := mvc.minCostFlow(SS, TT, totalDemand)
		// ideally flow == totalDemand
		res += cost
	}
	// output
	fmt.Println(res)
}
