package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

const INF int64 = 1 << 60

func dijkstra(n int, revAdj [][]Edge, highlighted []bool) []int64 {
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = INF
	}
	pq := &PriorityQueue{}
	heap.Init(pq)
	for i, hl := range highlighted {
		if hl {
			dist[i] = 0
			heap.Push(pq, Item{i, 0})
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		for _, e := range revAdj[u] {
			v := e.to
			nd := it.dist + int64(e.w)
			if nd < dist[v] {
				dist[v] = nd
				heap.Push(pq, Item{v, nd})
			}
		}
	}
	return dist
}

func computeCost(n int, adj, revAdj [][]Edge, highlighted []bool) int64 {
	// if no highlighted vertices, impossible
	anyHighlight := false
	for _, hl := range highlighted {
		if hl {
			anyHighlight = true
			break
		}
	}
	if !anyHighlight {
		return -1
	}
	dist := dijkstra(n, revAdj, highlighted)
	cost := int64(0)
	for i := 0; i < n; i++ {
		if highlighted[i] {
			continue
		}
		if dist[i] == INF {
			return -1
		}
		best := INF
		for _, e := range adj[i] {
			if dist[e.to]+int64(e.w) == dist[i] {
				if int64(e.w) < best {
					best = int64(e.w)
				}
			}
		}
		if best == INF {
			return -1
		}
		cost += best
	}
	return cost
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}
	adj := make([][]Edge, n)
	revAdj := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		u--
		v--
		adj[u] = append(adj[u], Edge{v, w})
		revAdj[v] = append(revAdj[v], Edge{u, w})
	}
	highlighted := make([]bool, n)
	queries := make([]struct {
		add bool
		v   int
	}, q)
	for i := 0; i < q; i++ {
		var typ string
		var v int
		fmt.Fscan(reader, &typ, &v)
		v--
		queries[i] = struct {
			add bool
			v   int
		}{typ == "+", v}
	}

	for i := 0; i < q; i++ {
		if queries[i].add {
			highlighted[queries[i].v] = true
		} else {
			highlighted[queries[i].v] = false
		}
		res := computeCost(n, adj, revAdj, highlighted)
		fmt.Fprintln(writer, res)
	}
}
