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

const INF int64 = 1 << 60

func dijkstra(n int, g [][]Edge, start int) []int64 {
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	dist[start] = 0
	pq := &PriorityQueue{}
	heap.Push(pq, Item{start, 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		v := it.node
		if it.dist != dist[v] {
			continue
		}
		for _, e := range g[v] {
			nd := it.dist + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	return dist
}

func bfsAll(n int, adj [][]int, dn []int64) ([][]int, [][]int64) {
	dist := make([][]int, n+1)
	best := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			dist[i][j] = -1
		}
		best[i] = make([]int64, n+1)
		for j := 0; j <= n; j++ {
			best[i][j] = INF
		}
		q := make([]int, n)
		head, tail := 0, 0
		dist[i][i] = 0
		q[tail] = i
		tail++
		for head < tail {
			v := q[head]
			head++
			d := dist[i][v]
			if dn[v] < best[i][d] {
				best[i][d] = dn[v]
			}
			for _, to := range adj[v] {
				if dist[i][to] == -1 {
					dist[i][to] = d + 1
					q[tail] = to
					tail++
				}
			}
		}
	}
	return dist, best
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		g := make([][]Edge, n+1)
		adj := make([][]int, n+1)
		edges := make([][3]int64, m)
		for i := 0; i < m; i++ {
			var u, v int
			var w int64
			fmt.Fscan(in, &u, &v, &w)
			edges[i] = [3]int64{int64(u), int64(v), w}
			g[u] = append(g[u], Edge{v, w})
			g[v] = append(g[v], Edge{u, w})
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		d1 := dijkstra(n, g, 1)
		dn := dijkstra(n, g, n)
		_, best := bfsAll(n, adj, dn)

		ans := d1[n]
		for _, e := range edges {
			u := int(e[0])
			v := int(e[1])
			w := e[2]
			// orientation u->v
			minCost := INF
			for d := 0; d <= n; d++ {
				if best[v][d] < INF {
					val := int64(d)*w + best[v][d]
					if val < minCost {
						minCost = val
					}
				}
			}
			cand := d1[u] + w + minCost
			if cand < ans {
				ans = cand
			}
			// orientation v->u
			minCost = INF
			for d := 0; d <= n; d++ {
				if best[u][d] < INF {
					val := int64(d)*w + best[u][d]
					if val < minCost {
						minCost = val
					}
				}
			}
			cand = d1[v] + w + minCost
			if cand < ans {
				ans = cand
			}
		}
		fmt.Fprintln(out, ans)
	}
}
