package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/big"
	"os"
	"sort"
)

type Edge struct {
	to  int
	w   int64
	idx int
}

type Item struct {
	v    int
	dist int64
}

type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

const INF int64 = 1 << 60

func dijkstra(start int, adj [][]Edge) []int64 {
	n := len(adj) - 1
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	pq := &MinHeap{}
	heap.Init(pq)
	dist[start] = 0
	heap.Push(pq, Item{v: start, dist: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.dist != dist[cur.v] {
			continue
		}
		for _, e := range adj[cur.v] {
			nd := cur.dist + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{v: e.to, dist: nd})
			}
		}
	}
	return dist
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, s, t int
	if _, err := fmt.Fscan(reader, &n, &m, &s, &t); err != nil {
		return
	}
	edges := make([][3]int64, m)
	adj := make([][]Edge, n+1)
	radj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		var l int64
		fmt.Fscan(reader, &a, &b, &l)
		edges[i] = [3]int64{int64(a), int64(b), l}
		adj[a] = append(adj[a], Edge{to: b, w: l, idx: i})
		radj[b] = append(radj[b], Edge{to: a, w: l, idx: i})
	}

	distS := dijkstra(s, adj)
	distT := dijkstra(t, radj)
	L := distS[t]

	// build DAG of edges on some shortest path
	dag := make([][]Edge, n+1)
	revDag := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := int(edges[i][0])
		v := int(edges[i][1])
		w := edges[i][2]
		if distS[u] != INF && distT[v] != INF && distS[u]+w+distT[v] == L {
			dag[u] = append(dag[u], Edge{to: v, w: w})
			revDag[v] = append(revDag[v], Edge{to: u, w: w})
		}
	}

	// sort nodes by distS for topological order
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i + 1
	}
	sort.Slice(order, func(i, j int) bool { return distS[order[i]] < distS[order[j]] })

	f := make([]*big.Int, n+1)
	for i := 1; i <= n; i++ {
		f[i] = big.NewInt(0)
	}
	f[s].SetInt64(1)
	for _, u := range order {
		for _, e := range dag[u] {
			f[e.to].Add(f[e.to], f[u])
		}
	}

	// reverse order for g
	revOrder := make([]int, n)
	copy(revOrder, order)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		revOrder[i], revOrder[j] = revOrder[j], revOrder[i]
	}
	g := make([]*big.Int, n+1)
	for i := 1; i <= n; i++ {
		g[i] = big.NewInt(0)
	}
	g[t].SetInt64(1)
	for _, u := range revOrder {
		for _, e := range revDag[u] {
			g[e.to].Add(g[e.to], g[u])
		}
	}

	total := new(big.Int).Set(f[t])

	for i := 0; i < m; i++ {
		u := int(edges[i][0])
		v := int(edges[i][1])
		w := edges[i][2]
		if distS[u] == INF || distT[v] == INF {
			fmt.Fprintln(writer, "NO")
			continue
		}
		pathLen := distS[u] + w + distT[v]
		if pathLen == L {
			prod := new(big.Int).Mul(f[u], g[v])
			if prod.Cmp(total) == 0 {
				fmt.Fprintln(writer, "YES")
				continue
			}
		}
		T := L - distS[u] - distT[v]
		if T > 1 {
			newW := T - 1
			if newW >= 1 && newW < w {
				fmt.Fprintf(writer, "CAN %d\n", w-newW)
				continue
			}
		}
		fmt.Fprintln(writer, "NO")
	}
}
