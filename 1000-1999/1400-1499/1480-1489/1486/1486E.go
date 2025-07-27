package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int
}

type Item struct {
	d    int64
	node int
	last int
	idx  int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].idx = i; pq[j].idx = j }
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.idx = len(*pq)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	return item
}

const INF int64 = 1 << 60

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	dist := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int64, 51)
		for j := 0; j <= 50; j++ {
			dist[i][j] = INF
		}
	}
	dist[1][0] = 0
	pq := &PriorityQueue{}
	heap.Push(pq, &Item{d: 0, node: 1, last: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*Item)
		if cur.d != dist[cur.node][cur.last] {
			continue
		}
		if cur.last == 0 {
			for _, e := range g[cur.node] {
				if dist[e.to][e.w] > cur.d {
					dist[e.to][e.w] = cur.d
					heap.Push(pq, &Item{d: cur.d, node: e.to, last: e.w})
				}
			}
		} else {
			for _, e := range g[cur.node] {
				nd := cur.d + int64(cur.last+e.w)*int64(cur.last+e.w)
				if dist[e.to][0] > nd {
					dist[e.to][0] = nd
					heap.Push(pq, &Item{d: nd, node: e.to, last: 0})
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		if dist[i][0] == INF {
			fmt.Fprint(writer, "-1")
		} else {
			fmt.Fprint(writer, dist[i][0])
		}
		if i == n {
			fmt.Fprintln(writer)
		} else {
			fmt.Fprint(writer, " ")
		}
	}
}
