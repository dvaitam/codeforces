package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

const inf int64 = 1 << 60

type Edge struct {
	to int
	w  int64
}

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
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

type Line struct {
	m, b int64
}

func (ln Line) value(x int64) int64 { return ln.m*x + ln.b }

type Node struct {
	ln          Line
	left, right *Node
}

func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	leftBetter := ln.value(l) < node.ln.value(l)
	midBetter := ln.value(mid) < node.ln.value(mid)
	if midBetter {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if leftBetter != midBetter {
		node.left = insert(node.left, l, mid, ln)
	} else {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return math.MaxInt64
	}
	res := node.ln.value(x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		val := query(node.left, l, mid, x)
		if val < res {
			res = val
		}
	} else {
		val := query(node.right, mid+1, r, x)
		if val < res {
			res = val
		}
	}
	return res
}

type LiChao struct {
	root        *Node
	left, right int64
}

func NewLiChao(l, r int64) *LiChao { return &LiChao{left: l, right: r} }
func (lc *LiChao) Insert(ln Line)  { lc.root = insert(lc.root, lc.left, lc.right, ln) }
func (lc *LiChao) Query(x int64) int64 {
	return query(lc.root, lc.left, lc.right, x)
}

func dijkstra(start []int64, g [][]Edge) []int64 {
	n := len(g) - 1
	dist := make([]int64, n+1)
	copy(dist, start)
	pq := &priorityQueue{}
	heap.Init(pq)
	for i := 1; i <= n; i++ {
		if dist[i] < inf {
			heap.Push(pq, &item{node: i, dist: dist[i]})
		}
	}
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		d := it.dist
		for _, e := range g[u] {
			nd := d + e.w
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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	dist[1] = 0
	dist = dijkstra(dist, g)

	for iter := 0; iter < k; iter++ {
		tree := NewLiChao(1, int64(n))
		for i := 1; i <= n; i++ {
			if dist[i] < inf {
				x := int64(i)
				tree.Insert(Line{m: -2 * x, b: dist[i] + x*x})
			}
		}
		newDist := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			x := int64(i)
			newDist[i] = x*x + tree.Query(x)
		}
		dist = dijkstra(newDist, g)
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, dist[i])
	}
	out.WriteByte('\n')
}
