package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] } // max-heap by index
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	m := 2 * n
	val := make([]int, m+1)
	cnt := make([]int, n+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &val[i])
		cnt[val[i]]++
	}

	adj := make([][]int, m+1)
	deg := make([]int, m+1)
	for i := 0; i < m-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	for i := 1; i <= m; i++ {
		deg[i] = len(adj[i])
	}

	removed := make([]bool, m+1)
	h := &MaxHeap{}
	heap.Init(h)
	for i := 1; i <= m; i++ {
		if deg[i] <= 1 && cnt[val[i]] > 1 {
			heap.Push(h, i)
		}
	}

	for h.Len() > 0 {
		v := heap.Pop(h).(int)
		if removed[v] || deg[v] > 1 || cnt[val[v]] <= 1 {
			continue
		}
		removed[v] = true
		cnt[val[v]]--
		for _, u := range adj[v] {
			if removed[u] {
				continue
			}
			deg[u]--
			if deg[u] == 1 && cnt[val[u]] > 1 {
				heap.Push(h, u)
			}
		}
	}

	ans := make([]int, 0, m)
	for i := 1; i <= m; i++ {
		if !removed[i] {
			ans = append(ans, i)
		}
	}

	fmt.Println(len(ans))
	for i, v := range ans {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
