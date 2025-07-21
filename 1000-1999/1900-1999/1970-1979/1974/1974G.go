package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// maxHeap for removing largest cost when over budget

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solve(m int, x int, costs []int) int {
	h := &intHeap{}
	heap.Init(h)
	total := 0
	for i := 1; i <= m; i++ {
		heap.Push(h, costs[i-1])
		total += costs[i-1]
		limit := x * (i - 1)
		for total > limit {
			removed := heap.Pop(h).(int)
			total -= removed
		}
	}
	return h.Len()
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var m, x int
		fmt.Fscan(in, &m, &x)
		costs := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &costs[i])
		}
		fmt.Fprintln(out, solve(m, x, costs))
	}
}
