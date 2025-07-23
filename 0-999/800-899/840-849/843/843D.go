package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type Edge struct {
	to int
	w  int64
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func dijkstra(n int, start int, edges []Edge, adj [][]int) []int64 {
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	dist[start] = 0
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Item{start, 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		for _, idx := range adj[u] {
			e := edges[idx]
			v := e.to
			nd := it.dist + e.w
			if nd < dist[v] {
				dist[v] = nd
				heap.Push(pq, Item{v, nd})
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)
	edges := make([]Edge, m)
	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		var c int64
		fmt.Fscan(reader, &a, &b, &c)
		edges[i] = Edge{to: b, w: c}
		adj[a] = append(adj[a], i)
	}

	var dist []int64
	valid := false
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v int
			fmt.Fscan(reader, &v)
			if !valid {
				dist = dijkstra(n, 1, edges, adj)
				valid = true
			}
			if dist[v] >= inf/2 {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, dist[v])
			}
		} else if t == 2 {
			valid = false
			var c int
			fmt.Fscan(reader, &c)
			for i := 0; i < c; i++ {
				var idx int
				fmt.Fscan(reader, &idx)
				idx--
				edges[idx].w++
			}
		}
	}
}
