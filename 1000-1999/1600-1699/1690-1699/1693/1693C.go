package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

// This program solves the problem described in problemC.txt.
// We need the minimum number of days to guarantee reaching city n
// from city 1 when we may either block one road or order a move each day.
// The optimal strategy can be computed using a Dijkstra-like process
// on the reverse graph. For each city we track how many outgoing
// edges have been processed. When exploring an edge v->u from the
// reverse graph, the potential cost to move from v through this edge is
// dp[u] + 1 + (outDeg[v] - processed[v]), since remaining outgoing roads
// must be blocked before moving. Processing nodes in order of dp gives
// the minimal guaranteed days.

type item struct {
	node int
	dist int
}

type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	rev := make([][]int, n+1)
	outDeg := make([]int, n+1)
	for i := 0; i < m; i++ {
		var v, u int
		fmt.Fscan(in, &v, &u)
		rev[u] = append(rev[u], v)
		outDeg[v]++
	}

	const inf = math.MaxInt32
	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = inf
	}
	processed := make([]int, n+1)

	pq := &priorityQueue{}
	heap.Init(pq)
	dist[n] = 0
	heap.Push(pq, item{n, 0})

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		u := cur.node
		if cur.dist != dist[u] {
			continue
		}
		for _, v := range rev[u] {
			processed[v]++
			cand := dist[u] + 1 + outDeg[v] - processed[v]
			if cand < dist[v] {
				dist[v] = cand
				heap.Push(pq, item{v, cand})
			}
		}
	}

	fmt.Fprintln(out, dist[1])
}
