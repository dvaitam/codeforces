package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	to   int
	cost int
}

type Item struct {
	node int
	dist int
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

const inf = int(1e9 + 7)

func dijkstra(start, v int, g [][]Edge) []int {
	dist := make([]int, v+1)
	for i := range dist {
		dist[i] = math.MaxInt64
	}
	dist[start] = 0
	pq := &PQ{Item{start, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.dist != dist[it.node] {
			continue
		}
		for _, e := range g[it.node] {
			nd := it.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	return dist
}

func maxMatching(adj [][]int, m, need int) bool {
	match := make([]int, m+1)
	for i := range match {
		match[i] = -1
	}
	var dfs func(int, []bool) bool
	dfs = func(u int, vis []bool) bool {
		for _, v := range adj[u] {
			if vis[v] {
				continue
			}
			vis[v] = true
			if match[v] == -1 || dfs(match[v], vis) {
				match[v] = u
				return true
			}
		}
		return false
	}
	cnt := 0
	for u := 0; u < len(adj); u++ {
		vis := make([]bool, m+1)
		if dfs(u, vis) {
			cnt++
			if cnt >= need {
				return true
			}
		}
	}
	return cnt >= need
}

func feasible(T int, starts []int, distMap map[int][]int, V int, K int) bool {
	adj := make([][]int, len(starts))
	for i, s := range starts {
		d := distMap[s]
		list := make([]int, 0)
		for v := 1; v <= V; v++ {
			if d[v] <= T {
				list = append(list, v)
			}
		}
		adj[i] = list
	}
	return maxMatching(adj, V, K)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var V, E, N, K int
	if _, err := fmt.Fscan(in, &V, &E, &N, &K); err != nil {
		return
	}

	starts := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &starts[i])
	}

	g := make([][]Edge, V+1)
	for i := 0; i < E; i++ {
		var a, b, t int
		fmt.Fscan(in, &a, &b, &t)
		g[a] = append(g[a], Edge{b, t})
		g[b] = append(g[b], Edge{a, t})
	}

	distMap := make(map[int][]int)
	for _, s := range starts {
		if _, ok := distMap[s]; !ok {
			distMap[s] = dijkstra(s, V, g)
		}
	}

	lo, hi := 0, 1731311
	ans := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if feasible(mid, starts, distMap, V, K) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	fmt.Fprintln(out, ans)
}
