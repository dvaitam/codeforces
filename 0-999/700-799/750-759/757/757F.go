package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

const INF int64 = 1 << 60

type Edge struct {
	to int
	w  int64
}

type Item struct {
	node int
	dist int64
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

const LOG = 20

var (
	up    [LOG][]int
	depth []int
	idom  []int
	tree  [][]int
)

func lca(u, v int) int {
	if u == 0 || v == 0 {
		if u == 0 {
			return v
		}
		return u
	}
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff>>i&1 == 1 {
			u = up[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}

func dfs(u int) int {
	sz := 1
	for _, v := range tree[u] {
		sz += dfs(v)
	}
	if u != 0 {
		if sz > best {
			best = sz
		}
	}
	return sz
}

var best int

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, s int
	if _, err := fmt.Fscan(in, &n, &m, &s); err != nil {
		return
	}
	adj := make([][]Edge, n+1)
	edges := make([][3]int64, m)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
		edges[i] = [3]int64{int64(u), int64(v), w}
	}
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	dist[s] = 0
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, Item{s, 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		u := cur.node
		d := cur.dist
		if d != dist[u] {
			continue
		}
		for _, e := range adj[u] {
			nd := d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}

	rev := make([][]int, n+1)
	for _, e := range edges {
		u := int(e[0])
		v := int(e[1])
		w := e[2]
		if dist[u]+w == dist[v] {
			rev[v] = append(rev[v], u)
		}
		if dist[v]+w == dist[u] {
			rev[u] = append(rev[u], v)
		}
	}

	order := make([]int, 0)
	for i := 1; i <= n; i++ {
		if dist[i] < INF {
			order = append(order, i)
		}
	}
	sort.Slice(order, func(i, j int) bool { return dist[order[i]] < dist[order[j]] })

	depth = make([]int, n+1)
	idom = make([]int, n+1)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, n+1)
	}
	idom[s] = s
	for i := 0; i < LOG; i++ {
		up[i][s] = s
	}
	depth[s] = 0

	for _, v := range order {
		if v == s {
			continue
		}
		preds := rev[v]
		if len(preds) == 0 {
			continue
		}
		id := preds[0]
		for _, p := range preds {
			if p == id {
				continue
			}
			id = lca(id, p)
		}
		idom[v] = id
		depth[v] = depth[id] + 1
		up[0][v] = id
		for k := 1; k < LOG; k++ {
			up[k][v] = up[k-1][up[k-1][v]]
		}
	}

	tree = make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if i == s || idom[i] == 0 {
			continue
		}
		tree[idom[i]] = append(tree[idom[i]], i)
	}

	best = 0
	dfs(s)
	fmt.Println(best)
}
