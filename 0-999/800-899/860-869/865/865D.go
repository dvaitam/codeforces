package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
func (h maxHeap) Peek() int { return h[0] }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	prices := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &prices[i])
	}
	h := &maxHeap{}
	var profit int64
	for i := n - 1; i >= 0; i-- {
		if i < n-1 {
			heap.Push(h, prices[i+1])
		}
		if h.Len() > 0 && h.Peek() > prices[i] {
			sell := heap.Pop(h).(int)
			profit += int64(sell - prices[i])
			heap.Push(h, prices[i])
		}
	}
	fmt.Println(profit)
}
