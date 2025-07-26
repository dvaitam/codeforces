package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to, moment int
}

type Item struct {
	step, city int
}

type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].step < h[j].step }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	adj := make([][]Edge, n+1)
	for t := 1; t <= m; t++ {
		var cnt int
		fmt.Fscan(reader, &cnt)
		for i := 0; i < cnt; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			adj[u] = append(adj[u], Edge{v, t})
			adj[v] = append(adj[v], Edge{u, t})
		}
	}

	var k int
	fmt.Fscan(reader, &k)
	times := make([]int, k+1)
	for i := 1; i <= k; i++ {
		fmt.Fscan(reader, &times[i])
	}

	visits := make([][]int, m+1)
	for i := 1; i <= k; i++ {
		t := times[i]
		visits[t] = append(visits[t], i)
	}

	const INF = int(1e9)
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	dist[1] = 0
	pq := &MinHeap{{0, 1}}
	heap.Init(pq)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(Item)
		if item.step != dist[item.city] {
			continue
		}
		if item.city == n {
			fmt.Fprintln(writer, item.step)
			return
		}
		if item.step >= k {
			continue
		}
		for _, e := range adj[item.city] {
			arr := visits[e.moment]
			idx := sort.Search(len(arr), func(i int) bool { return arr[i] > item.step })
			if idx < len(arr) {
				step := arr[idx]
				if step < dist[e.to] {
					dist[e.to] = step
					heap.Push(pq, Item{step, e.to})
				}
			}
		}
	}

	fmt.Fprintln(writer, -1)
}
