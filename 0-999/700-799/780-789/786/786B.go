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

var (
	n       int
	g       [][]edge
	offOut  int
	offIn   int
	sizeSeg int
)

func addEdge(u, v int, w int64) {
	g[u] = append(g[u], edge{v, w})
}

func buildOut(node, l, r int) {
	idx := offOut + node
	if l == r {
		addEdge(idx, l, 0)
		return
	}
	mid := (l + r) / 2
	buildOut(node*2, l, mid)
	buildOut(node*2+1, mid+1, r)
	addEdge(idx, offOut+node*2, 0)
	addEdge(idx, offOut+node*2+1, 0)
}

func buildIn(node, l, r int) {
	idx := offIn + node
	if l == r {
		addEdge(l, idx, 0)
		return
	}
	mid := (l + r) / 2
	buildIn(node*2, l, mid)
	buildIn(node*2+1, mid+1, r)
	addEdge(offIn+node*2, idx, 0)
	addEdge(offIn+node*2+1, idx, 0)
}

func addOutRange(node, l, r, ql, qr, v int, w int64) {
	if ql <= l && r <= qr {
		addEdge(v, offOut+node, w)
		return
	}
	mid := (l + r) / 2
	if ql <= mid {
		addOutRange(node*2, l, mid, ql, qr, v, w)
	}
	if qr > mid {
		addOutRange(node*2+1, mid+1, r, ql, qr, v, w)
	}
}

func addInRange(node, l, r, ql, qr, v int, w int64) {
	if ql <= l && r <= qr {
		addEdge(offIn+node, v, w)
		return
	}
	mid := (l + r) / 2
	if ql <= mid {
		addInRange(node*2, l, mid, ql, qr, v, w)
	}
	if qr > mid {
		addInRange(node*2+1, mid+1, r, ql, qr, v, w)
	}
}

func dijkstra(s int) []int64 {
	nNodes := offIn + sizeSeg
	dist := make([]int64, nNodes+1)
	for i := range dist {
		dist[i] = 1<<63 - 1
	}
	dist[s] = 0
	pq := &priorityQueue{}
	heap.Push(pq, &item{node: s, dist: 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		d := it.dist
		for _, e := range g[u] {
			nd := d + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, &item{node: e.to, dist: nd})
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q, s int
	if _, err := fmt.Fscan(in, &n, &q, &s); err != nil {
		return
	}

	sizeSeg = 1
	for sizeSeg < n {
		sizeSeg <<= 1
	}
	offOut = n
	offIn = n + sizeSeg*2
	g = make([][]edge, offIn+sizeSeg*2+1)

	buildOut(1, 1, sizeSeg)
	buildIn(1, 1, sizeSeg)

	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var v, u int
			var w int64
			fmt.Fscan(in, &v, &u, &w)
			addEdge(v, u, w)
		} else if t == 2 {
			var v, l, r int
			var w int64
			fmt.Fscan(in, &v, &l, &r, &w)
			addOutRange(1, 1, sizeSeg, l, r, v, w)
		} else if t == 3 {
			var v, l, r int
			var w int64
			fmt.Fscan(in, &v, &l, &r, &w)
			addInRange(1, 1, sizeSeg, l, r, v, w)
		}
	}

	dist := dijkstra(s)

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		if dist[i] == 1<<63-1 {
			fmt.Fprint(out, -1)
		} else {
			fmt.Fprint(out, dist[i])
		}
	}
	fmt.Fprintln(out)
}

// Priority queue implementation

type item struct {
	node int
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
	x := old[n-1]
	*pq = old[:n-1]
	return x
}
