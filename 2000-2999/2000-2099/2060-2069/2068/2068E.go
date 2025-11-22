package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

// For every node u, blockCost[u][v] is the length of the shortest path
// from u to the target if the edge (u,v) gets blocked just before we
// traverse it from u. Since only that edge is removed, the best detour
// is to pick the best remaining neighbour; this depends only on the
// two smallest values of 1+dist[nei].

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	adj := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	target := n - 1

	// BFS from target to get dist
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[target] = 0
	q = append(q, target)
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if dist[v] == -1 {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	if dist[0] == -1 {
		fmt.Println(-1)
		return
	}

	// Precompute best1, best2, cnt1 for each node
	best1 := make([]int64, n)
	best2 := make([]int64, n)
	cnt1 := make([]int, n)
	for u := 0; u < n; u++ {
		b1, b2 := INF, INF
		c1 := 0
		for _, v := range adj[u] {
			val := int64(1 + dist[v])
			if val < b1 {
				b2 = b1
				b1 = val
				c1 = 1
			} else if val == b1 {
				c1++
			} else if val < b2 {
				b2 = val
			}
		}
		best1[u] = b1
		best2[u] = b2
		cnt1[u] = c1
	}

	// helper to get block cost
	blockCost := func(u, v int) int64 {
		val := int64(1 + dist[v])
		b1, b2, c1 := best1[u], best2[u], cnt1[u]
		if val > b1 {
			return b1
		}
		if val == b1 {
			if c1 > 1 {
				return b1
			}
			return b2
		}
		// val < b1 never happens because b1 is min
		return b1
	}

	// Dijkstra-like relaxation but edges used in reverse: knowing f[v]
	// we can try using v as next step from neighbour u.
	f := make([]int64, n)
	for i := range f {
		f[i] = INF
	}
	f[target] = 0
	pq := &MinPQ{}
	heap.Push(pq, Item{node: target, val: 0})

	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		u := it.node
		if it.val != f[u] {
			continue
		}
		for _, v := range adj[u] { // try to update neighbour v's parent (which is u in recurrence) ?
			// we want to update node v using neighbour u as next step
			bc := blockCost(v, u)
			cand := int64(1)
			if f[u] >= INF/2 {
				cand = INF
			} else {
				cand += f[u]
			}
			if bc > cand {
				cand = bc
			}
			if cand < f[v] {
				f[v] = cand
				heap.Push(pq, Item{node: v, val: cand})
			}
		}
	}

	ans := f[0]
	if ans >= INF/2 {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}

// Priority queue

type Item struct {
	node int
	val  int64
}

type MinPQ []Item

func (pq MinPQ) Len() int            { return len(pq) }
func (pq MinPQ) Less(i, j int) bool  { return pq[i].val < pq[j].val }
func (pq MinPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *MinPQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *MinPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}
