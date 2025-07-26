package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Edge represents a road between cities with a hardness weight.
type Edge struct {
	to int
	w  int
}

// Node is a state for Dijkstra consisting of a city and skill level exponent.
type Node struct {
	d int64 // accumulated time
	v int
	e int
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].d < pq[j].d }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Node)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func dijkstra(g [][]Edge, courses string, T int, start, end int) int64 {
	const maxExp = 20
	n := len(g) - 1
	INF := int64(1<<63 - 1)
	dist := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int64, maxExp+1)
		for j := 0; j <= maxExp; j++ {
			dist[i][j] = INF
		}
	}
	pq := &PriorityQueue{}
	dist[start][0] = 0
	heap.Push(pq, Node{0, start, 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Node)
		if cur.d != dist[cur.v][cur.e] {
			continue
		}
		if cur.v == end {
			return cur.d
		}
		c := 1 << cur.e
		// take course if possible
		if courses[cur.v-1] == '1' && cur.e < maxExp {
			nd := cur.d + int64(T)
			if nd < dist[cur.v][cur.e+1] {
				dist[cur.v][cur.e+1] = nd
				heap.Push(pq, Node{nd, cur.v, cur.e + 1})
			}
		}
		// traverse edges
		for _, e := range g[cur.v] {
			cost := (e.w + c - 1) / c
			nd := cur.d + int64(cost)
			if nd < dist[e.to][cur.e] {
				dist[e.to][cur.e] = nd
				heap.Push(pq, Node{nd, e.to, cur.e})
			}
		}
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, T int
		fmt.Fscan(reader, &n, &T)
		g := make([][]Edge, n+1)
		for i := 0; i < n-1; i++ {
			var u, v, w int
			fmt.Fscan(reader, &u, &v, &w)
			g[u] = append(g[u], Edge{v, w})
			g[v] = append(g[v], Edge{u, w})
		}
		var s string
		fmt.Fscan(reader, &s)
		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var a, b int
			fmt.Fscan(reader, &a, &b)
			ans := dijkstra(g, s, T, a, b)
			fmt.Fprintln(writer, ans)
		}
	}
}
