package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	w  int64
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k, q int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &q); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	type OrigEdge struct {
		u, v int
		w    int64
	}
	orig := make([]OrigEdge, m)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(reader, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
		orig[i] = OrigEdge{u, v, w}
	}

	const INF int64 = 1 << 60
	dist := make([]int64, n+1)
	belong := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}

	pq := &PriorityQueue{}
	heap.Init(pq)
	for i := 1; i <= k; i++ {
		dist[i] = 0
		belong[i] = i
		heap.Push(pq, Item{i, 0})
	}

	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u := it.node
		d := it.dist
		if d != dist[u] {
			continue
		}
		for _, e := range g[u] {
			nd := d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				belong[e.to] = belong[u]
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}

	type E2 struct {
		u, v int
		w    int64
	}
	e2 := make([]E2, 0, m)
	for _, e := range orig {
		bu := belong[e.u]
		bv := belong[e.v]
		if bu != bv {
			cost := dist[e.u] + dist[e.v] + e.w
			e2 = append(e2, E2{bu, bv, cost})
		}
	}
	sort.Slice(e2, func(i, j int) bool { return e2[i].w < e2[j].w })

	parent := make([]int, k+1)
	for i := 1; i <= k; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	unite := func(a, b int) bool {
		fa := find(a)
		fb := find(b)
		if fa == fb {
			return false
		}
		parent[fb] = fa
		return true
	}

	adj := make([][]Edge, k+1)
	for _, e := range e2 {
		if unite(e.u, e.v) {
			adj[e.u] = append(adj[e.u], Edge{e.v, e.w})
			adj[e.v] = append(adj[e.v], Edge{e.u, e.w})
		}
	}

	LOG := 17
	up := make([][]int, LOG)
	mx := make([][]int64, LOG)
	for i := 0; i < LOG; i++ {
		up[i] = make([]int, k+1)
		mx[i] = make([]int64, k+1)
	}
	depth := make([]int, k+1)

	var dfs func(int, int)
	dfs = func(u, p int) {
		for _, e := range adj[u] {
			if e.to == p {
				continue
			}
			up[0][e.to] = u
			mx[0][e.to] = e.w
			depth[e.to] = depth[u] + 1
			dfs(e.to, u)
		}
	}

	dfs(1, 0)
	for i := 1; i < LOG; i++ {
		for v := 1; v <= k; v++ {
			up[i][v] = up[i-1][up[i-1][v]]
			mx[i][v] = max64(mx[i-1][v], mx[i-1][up[i-1][v]])
		}
	}

	query := func(a, b int) int64 {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		ans := int64(0)
		diff := depth[a] - depth[b]
		for i := LOG - 1; i >= 0; i-- {
			if diff&(1<<i) != 0 {
				ans = max64(ans, mx[i][a])
				a = up[i][a]
			}
		}
		if a == b {
			return ans
		}
		for i := LOG - 1; i >= 0; i-- {
			if up[i][a] != up[i][b] {
				ans = max64(ans, mx[i][a])
				ans = max64(ans, mx[i][b])
				a = up[i][a]
				b = up[i][b]
			}
		}
		ans = max64(ans, mx[0][a])
		ans = max64(ans, mx[0][b])
		return ans
	}

	for ; q > 0; q-- {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		fmt.Fprintln(writer, query(a, b))
	}
}
