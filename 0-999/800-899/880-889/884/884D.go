package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int64

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	h := &IntHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		heap.Push(h, x)
	}

	if n == 1 {
		fmt.Fprintln(out, 0)
		return
	}
	if n%2 == 0 {
		heap.Push(h, 0)
	}

	var cost int64
	for h.Len() > 1 {
		a := heap.Pop(h).(int64)
		b := heap.Pop(h).(int64)
		c := heap.Pop(h).(int64)
		sum := a + b + c
		cost += sum
		heap.Push(h, sum)
	}

	fmt.Fprintln(out, cost)
}
