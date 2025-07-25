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

type Item struct {
	dist int64
	node int
	slow int
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		g := make([][]Edge, n+1)
		for i := 0; i < m; i++ {
			var u, v int
			var w int64
			fmt.Fscan(in, &u, &v, &w)
			g[u] = append(g[u], Edge{to: v, cost: w})
			g[v] = append(g[v], Edge{to: u, cost: w})
		}

		s := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &s[i])
		}

		maxS := 1000
		dist := make([][]int64, n+1)
		for i := 0; i <= n; i++ {
			dist[i] = make([]int64, maxS+1)
			for j := 0; j <= maxS; j++ {
				dist[i][j] = INF
			}
		}

		pq := &PriorityQueue{}
		startSlow := s[1]
		dist[1][startSlow] = 0
		heap.Push(pq, Item{dist: 0, node: 1, slow: startSlow})

		for pq.Len() > 0 {
			cur := heap.Pop(pq).(Item)
			if cur.dist != dist[cur.node][cur.slow] {
				continue
			}
			u := cur.node
			d := cur.dist
			slow := cur.slow

			for _, e := range g[u] {
				v := e.to
				ns := slow
				if s[v] < ns {
					ns = s[v]
				}
				nd := d + e.cost*int64(slow)
				if nd < dist[v][ns] {
					dist[v][ns] = nd
					heap.Push(pq, Item{dist: nd, node: v, slow: ns})
				}
			}
		}

		ans := INF
		for i := 1; i <= maxS; i++ {
			if dist[n][i] < ans {
				ans = dist[n][i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
