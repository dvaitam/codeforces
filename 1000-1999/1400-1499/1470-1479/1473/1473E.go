package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to   int
	cost int64
}

type State struct {
	node int
	mask int
	dist int64
}

type PQ []State

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

const INF int64 = 1<<63 - 1

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	dist := make([][4]int64, n+1)
	for i := 1; i <= n; i++ {
		for j := 0; j < 4; j++ {
			dist[i][j] = INF
		}
	}
	dist[1][0] = 0

	pq := &PQ{}
	heap.Push(pq, State{1, 0, 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(State)
		if cur.dist != dist[cur.node][cur.mask] {
			continue
		}
		for _, e := range g[cur.node] {
			v := e.to
			w := e.cost
			if cur.mask == 0 {
				if nd := cur.dist + w; nd < dist[v][0] {
					dist[v][0] = nd
					heap.Push(pq, State{v, 0, nd})
				}
				if nd := cur.dist + 2*w; nd < dist[v][1] {
					dist[v][1] = nd
					heap.Push(pq, State{v, 1, nd})
				}
				if nd := cur.dist; nd < dist[v][2] {
					dist[v][2] = nd
					heap.Push(pq, State{v, 2, nd})
				}
				if nd := cur.dist + w; nd < dist[v][3] {
					dist[v][3] = nd
					heap.Push(pq, State{v, 3, nd})
				}
			} else if cur.mask == 1 {
				if nd := cur.dist + w; nd < dist[v][1] {
					dist[v][1] = nd
					heap.Push(pq, State{v, 1, nd})
				}
				if nd := cur.dist; nd < dist[v][3] {
					dist[v][3] = nd
					heap.Push(pq, State{v, 3, nd})
				}
			} else if cur.mask == 2 {
				if nd := cur.dist + w; nd < dist[v][2] {
					dist[v][2] = nd
					heap.Push(pq, State{v, 2, nd})
				}
				if nd := cur.dist + 2*w; nd < dist[v][3] {
					dist[v][3] = nd
					heap.Push(pq, State{v, 3, nd})
				}
			} else if cur.mask == 3 {
				if nd := cur.dist + w; nd < dist[v][3] {
					dist[v][3] = nd
					heap.Push(pq, State{v, 3, nd})
				}
			}
		}
	}

	for i := 2; i <= n; i++ {
		fmt.Fprint(out, dist[i][3])
		if i != n {
			fmt.Fprint(out, " ")
		}
	}
	fmt.Fprintln(out)
}
