package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	cost int64
}

type Item struct {
	node int
	dist int64
}

// Priority queue implementation for Dijkstra
// Implements heap.Interface

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

const INF int64 = 1<<63 - 1

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	size := 2*n + 1
	g := make([][]Edge, size)

	// zero-cost edges to switch from forward to reverse layer
	for i := 1; i <= n; i++ {
		g[i] = append(g[i], Edge{i + n, 0})
	}

	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})         // forward edge
		g[v+n] = append(g[v+n], Edge{u + n, w}) // reverse edge in second layer
	}

	dist := make([]int64, size)
	for i := range dist {
		dist[i] = INF
	}

	pq := &PQ{}
	dist[1] = 0
	heap.Push(pq, Item{1, 0})
	visited := make([]bool, size)

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		v := cur.node
		if visited[v] {
			continue
		}
		visited[v] = true
		d := cur.dist
		for _, e := range g[v] {
			nd := d + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 2; i <= n; i++ {
		if dist[i+n] >= INF/2 {
			fmt.Fprint(out, -1)
		} else {
			fmt.Fprint(out, dist[i+n])
		}
		if i < n {
			fmt.Fprint(out, " ")
		}
	}
	fmt.Fprintln(out)
	out.Flush()
}
