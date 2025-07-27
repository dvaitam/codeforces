package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	h := &IntHeap{}
	heap.Init(h)
	var sum int64
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		sum += x
		heap.Push(h, x)
		if sum < 0 {
			smallest := heap.Pop(h).(int64)
			sum -= smallest
		}
	}
	fmt.Fprintln(out, h.Len())
}
