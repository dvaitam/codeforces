package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents an edge in the flow network
type Edge struct {
	to, rev, cap, cost int
}

// Graph holds the flow network
type Graph struct {
	n   int
	adj [][]*Edge
}

// NewGraph creates a new graph with n nodes
func NewGraph(n int) *Graph {
	g := &Graph{n: n, adj: make([][]*Edge, n)}
	return g
}

// AddEdge adds an edge u->v with capacity cap and cost, and reverse edge
func (g *Graph) AddEdge(u, v, cap, cost int) {
	a := &Edge{to: v, rev: len(g.adj[v]), cap: cap, cost: cost}
	b := &Edge{to: u, rev: len(g.adj[u]), cap: 0, cost: -cost}
	g.adj[u] = append(g.adj[u], a)
	g.adj[v] = append(g.adj[v], b)
}

// Item for priority queue
type Item struct {
	v, dist int
}

// PQ is a min-heap of Items
type PQ struct {
	items []Item
}

// Len returns the number of items
func (pq *PQ) Len() int { return len(pq.items) }

// Push adds an item to the heap
func (pq *PQ) Push(it Item) {
	pq.items = append(pq.items, it)
	i := len(pq.items) - 1
	for i > 0 {
		p := (i - 1) / 2
		if pq.items[p].dist <= pq.items[i].dist {
			break
		}
		pq.items[p], pq.items[i] = pq.items[i], pq.items[p]
		i = p
	}
}

// Pop removes and returns the smallest item
func (pq *PQ) Pop() Item {
	min := pq.items[0]
	n := len(pq.items)
	pq.items[0] = pq.items[n-1]
	pq.items = pq.items[:n-1]
	i := 0
	for {
		left, right := 2*i+1, 2*i+2
		smallest := i
		if left < len(pq.items) && pq.items[left].dist < pq.items[smallest].dist {
			smallest = left
		}
		if right < len(pq.items) && pq.items[right].dist < pq.items[smallest].dist {
			smallest = right
		}
		if smallest == i {
			break
		}
		pq.items[i], pq.items[smallest] = pq.items[smallest], pq.items[i]
		i = smallest
	}
	return min
}

// MinCostFlow computes min cost flow from s to t, stopping when no negative cost path exists
// Returns total flow and total cost
func (g *Graph) MinCostFlow(s, t int) (int, int) {
	n := g.n
	const INF = int(1e9)
	flow, cost := 0, 0
	h := make([]int, n) // potentials
	prevv := make([]int, n)
	preve := make([]int, n)
	for {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = INF
		}
		dist[s] = 0
		pq := &PQ{}
		pq.Push(Item{v: s, dist: 0})
		for pq.Len() > 0 {
			it := pq.Pop()
			v := it.v
			if dist[v] < it.dist {
				continue
			}
			for i, e := range g.adj[v] {
				if e.cap > 0 && dist[e.to] > dist[v]+e.cost+h[v]-h[e.to] {
					dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
					prevv[e.to] = v
					preve[e.to] = i
					pq.Push(Item{v: e.to, dist: dist[e.to]})
				}
			}
		}
		if dist[t] == INF {
			break
		}
		// real distance
		if dist[t]+h[t] >= 0 {
			break
		}
		for v := 0; v < n; v++ {
			if dist[v] < INF {
				h[v] += dist[v]
			}
		}
		// add as much as possible (here 1)
		d := INF
		for v := t; v != s; v = prevv[v] {
			e := g.adj[prevv[v]][preve[v]]
			if e.cap < d {
				d = e.cap
			}
		}
		if d == INF {
			break
		}
		flow += d
		cost += d * h[t]
		for v := t; v != s; v = prevv[v] {
			e := g.adj[prevv[v]][preve[v]]
			e.cap -= d
			g.adj[v][e.rev].cap += d
		}
	}
	return flow, cost
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	sBuf := make([]byte, n)
	fmt.Fscan(in, &sBuf)
	s := string(sBuf)
	var m int
	fmt.Fscan(in, &m)
	words := make([]string, m)
	pts := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &words[i], &pts[i])
	}
	var x int
	fmt.Fscan(in, &x)
	// build graph nodes 0..n
	g := NewGraph(n + 1)
	// edges i->i+1 cap x cost 0
	for i := 0; i < n; i++ {
		g.AddEdge(i, i+1, x, 0)
	}
	// for each word, find occurrences
	for i := 0; i < m; i++ {
		w := words[i]
		L := len(w)
		for j := 0; j+L <= n; j++ {
			if s[j:j+L] == w {
				g.AddEdge(j, j+L, 1, -pts[i])
			}
		}
	}
	// compute min cost flow from 0 to n
	_, cost := g.MinCostFlow(0, n)
	// answer is -cost
	fmt.Println(-cost)
}
