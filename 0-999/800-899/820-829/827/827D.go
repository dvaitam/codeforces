package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
	id int
}

type Item struct {
	node int
	val  int64
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].val < pq[j].val }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

const INF int64 = 1 << 62

func minMaxPath(adj [][]Edge, banned int, start, end int) int64 {
	n := len(adj)
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	pq := &PQ{}
	heap.Push(pq, Item{start, 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.val > dist[cur.node] {
			continue
		}
		if cur.node == end {
			return cur.val
		}
		for _, e := range adj[cur.node] {
			if e.id == banned {
				continue
			}
			nv := cur.val
			if e.w > nv {
				nv = e.w
			}
			if nv < dist[e.to] {
				dist[e.to] = nv
				heap.Push(pq, Item{e.to, nv})
			}
		}
	}
	return INF
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	edges := make([]struct {
		u, v int
		w    int64
	}, m)
	adj := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		var c int64
		fmt.Fscan(reader, &u, &v, &c)
		u--
		v--
		edges[i] = struct {
			u, v int
			w    int64
		}{u, v, c}
		adj[u] = append(adj[u], Edge{v, c, i})
		adj[v] = append(adj[v], Edge{u, c, i})
	}

	for i := 0; i < m; i++ {
		u := edges[i].u
		v := edges[i].v
		mm := minMaxPath(adj, i, u, v)
		if mm == INF {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, mm-1)
		}
	}
}
