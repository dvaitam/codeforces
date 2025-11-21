package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Edge struct {
	to int
	w  int64
}

type State struct {
	dist int64
	v    int
	used int
}

type PriorityQueue []State

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(start int, adj [][]Edge, horse []bool) [2][]int64 {
	n := len(adj) - 1
	const inf = math.MaxInt64 / 4
	dist := [2][]int64{make([]int64, n+1), make([]int64, n+1)}
	for s := 0; s < 2; s++ {
		for i := 1; i <= n; i++ {
			dist[s][i] = inf
		}
	}
	dist[0][start] = 0
	pq := PriorityQueue{{0, start, 0}}
	heap.Init(&pq)
	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(State)
		if cur.dist != dist[cur.used][cur.v] {
			continue
		}
		// transition to used horse state if horse present
		if cur.used == 0 && horse[cur.v] {
			if cur.dist < dist[1][cur.v] {
				dist[1][cur.v] = cur.dist
				heap.Push(&pq, State{cur.dist, cur.v, 1})
			}
		}
		for _, e := range adj[cur.v] {
			w := e.w
			if cur.used == 1 {
				w /= 2
			}
			nd := cur.dist + w
			if nd < dist[cur.used][e.to] {
				dist[cur.used][e.to] = nd
				heap.Push(&pq, State{nd, e.to, cur.used})
			}
		}
	}
	return dist
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, h int
		fmt.Fscan(in, &n, &m, &h)
		horse := make([]bool, n+1)
		for i := 0; i < h; i++ {
			var x int
			fmt.Fscan(in, &x)
			horse[x] = true
		}
		adj := make([][]Edge, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			var w int64
			fmt.Fscan(in, &u, &v, &w)
			adj[u] = append(adj[u], Edge{v, w})
			adj[v] = append(adj[v], Edge{u, w})
		}
		distA := dijkstra(1, adj, horse)
		distB := dijkstra(n, adj, horse)
		const inf = math.MaxInt64 / 4
		best := int64(inf)
		for v := 1; v <= n; v++ {
			da := distA[0][v]
			if distA[1][v] < da {
				da = distA[1][v]
			}
			db := distB[0][v]
			if distB[1][v] < db {
				db = distB[1][v]
			}
			if da == inf || db == inf {
				continue
			}
			maxv := da
			if db > maxv {
				maxv = db
			}
			if maxv < best {
				best = maxv
			}
		}
		if best == inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, best)
		}
	}
}
