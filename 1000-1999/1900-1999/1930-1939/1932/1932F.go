package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	pos int
	r   int
	id  int
}

type maxHeap []pair

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].r > h[j].r }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(pair))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
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
		diff := make([]int, n+3)
		startEvents := make([]pair, 0, m)
		endEvents := make([]pair, 0, m)
		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			diff[l]++
			if r+1 <= n+1 {
				diff[r+1]--
			}
			startEvents = append(startEvents, pair{pos: l, r: r, id: i})
			endEvents = append(endEvents, pair{pos: r + 1, id: i})
		}
		sort.Slice(startEvents, func(i, j int) bool { return startEvents[i].pos < startEvents[j].pos })
		sort.Slice(endEvents, func(i, j int) bool { return endEvents[i].pos < endEvents[j].pos })
		// prefix sum for g
		g := make([]int, n+2)
		cur := 0
		for i := 1; i <= n; i++ {
			cur += diff[i]
			g[i] = cur
		}
		// sweep line to compute R
		removed := make([]bool, m)
		var pq maxHeap
		heap.Init(&pq)
		R := make([]int, n+2)
		si, ei := 0, 0
		for i := 1; i <= n; i++ {
			for si < len(startEvents) && startEvents[si].pos == i {
				heap.Push(&pq, pair{r: startEvents[si].r, id: startEvents[si].id})
				si++
			}
			for ei < len(endEvents) && endEvents[ei].pos == i {
				removed[endEvents[ei].id] = true
				ei++
			}
			for pq.Len() > 0 {
				top := pq[0]
				if removed[top.id] {
					heap.Pop(&pq)
				} else {
					break
				}
			}
			if pq.Len() > 0 {
				R[i] = pq[0].r
			} else {
				R[i] = i - 1
			}
		}
		// DP from right
		dp := make([]int, n+2)
		for i := n; i >= 1; i-- {
			skip := dp[i+1]
			j := R[i] + 1
			if j > n+1 {
				j = n + 1
			}
			feed := g[i] + dp[j]
			if feed > skip {
				dp[i] = feed
			} else {
				dp[i] = skip
			}
		}
		fmt.Fprintln(out, dp[1])
	}
}
