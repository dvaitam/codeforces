package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int64
}

type Item struct {
	v int
	d int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].d < pq[j].d }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(n int, start int, adj [][]Item) []int64 {
	const INF int64 = 1 << 62
	dist := make([]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
	}
	dist[start] = 0
	pq := &PriorityQueue{Item{start, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.d != dist[cur.v] {
			continue
		}
		for _, e := range adj[cur.v] {
			nd := cur.d + e.d
			if nd < dist[e.v] {
				dist[e.v] = nd
				heap.Push(pq, Item{e.v, nd})
			}
		}
	}
	return dist
}

type pair struct {
	d    int64
	a, b int
}

func computeAnswer(n int, edges []Edge) int64 {
	adj := make([][]Item, n)
	for _, e := range edges {
		u := e.u - 1
		v := e.v - 1
		adj[u] = append(adj[u], Item{v, e.w})
		adj[v] = append(adj[v], Item{u, e.w})
	}
	const INF int64 = 1 << 62
	pairs := make([]pair, 0)
	for i := 0; i < n; i++ {
		dist := dijkstra(n, i, adj)
		for j := i + 1; j < n; j++ {
			if dist[j] < INF {
				pairs = append(pairs, pair{dist[j], i, j})
			}
		}
	}
	// sort by distance
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].d < pairs[j].d })
	ans := INF
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			a1 := pairs[i]
			a2 := pairs[j]
			if a1.a != a2.a && a1.a != a2.b && a1.b != a2.a && a1.b != a2.b {
				if a1.d+a2.d < ans {
					ans = a1.d + a2.d
				}
			}
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
	}
	var q int
	fmt.Fscan(reader, &q)
	ans := computeAnswer(n, edges)
	fmt.Fprintln(writer, ans)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 0 {
			var v, u int
			fmt.Fscan(reader, &v, &u)
			for i := 0; i < len(edges); i++ {
				if (edges[i].u == v && edges[i].v == u) || (edges[i].u == u && edges[i].v == v) {
					edges = append(edges[:i], edges[i+1:]...)
					break
				}
			}
		} else {
			var v, u int
			var w int64
			fmt.Fscan(reader, &v, &u, &w)
			edges = append(edges, Edge{v, u, w})
		}
		ans = computeAnswer(n, edges)
		fmt.Fprintln(writer, ans)
	}
}
