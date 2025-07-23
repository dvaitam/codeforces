package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type edge struct {
	to   int
	cost int64
}

type item struct {
	node int
	dist int64
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

const INF int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([][]edge, n+1)
	for i := 0; i < m; i++ {
		var v, u int
		var w int64
		fmt.Fscan(in, &v, &u, &w)
		w *= 2
		adj[v] = append(adj[v], edge{u, w})
		adj[u] = append(adj[u], edge{v, w})
	}

	dist := make([]int64, n+1)
	h := &minHeap{}
	heap.Init(h)

	for i := 1; i <= n; i++ {
		var a int64
		fmt.Fscan(in, &a)
		dist[i] = a
		heap.Push(h, item{node: i, dist: a})
	}

	for h.Len() > 0 {
		it := heap.Pop(h).(item)
		u := it.node
		d := it.dist
		if d != dist[u] {
			continue
		}
		for _, e := range adj[u] {
			nd := d + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(h, item{node: e.to, dist: nd})
			}
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, dist[i])
	}
	out.WriteByte('\n')
}
