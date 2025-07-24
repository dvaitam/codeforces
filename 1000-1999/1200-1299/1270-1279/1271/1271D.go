package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i], &b[i], &c[i])
	}
	last := make([]int, n+1)
	for i := 1; i <= n; i++ {
		last[i] = i
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if u > v {
			if u > last[v] {
				last[v] = u
			}
		}
	}

	base := make([]int, n+1)
	base[0] = k
	for i := 1; i <= n; i++ {
		if base[i-1] < a[i] {
			fmt.Println(-1)
			return
		}
		base[i] = base[i-1] + b[i]
	}

	cap := make([]int, n+1)
	cap[n] = base[n]
	for i := 1; i < n; i++ {
		cap[i] = base[i] - a[i+1]
		if cap[i] < 0 {
			fmt.Println(-1)
			return
		}
	}

	buckets := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		d := last[i]
		buckets[d] = append(buckets[d], c[i])
	}

	h := &IntHeap{}
	heap.Init(h)
	total := 0
	size := 0
	for t := 1; t <= n; t++ {
		for _, val := range buckets[t] {
			heap.Push(h, val)
			total += val
			size++
		}
		for size > cap[t] {
			removed := heap.Pop(h).(int)
			total -= removed
			size--
		}
	}
	fmt.Println(total)
}
