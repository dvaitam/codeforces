package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// MaxHeap implements a max-heap of ints.
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		h := &MaxHeap{}
		var ans int64
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x > 0 {
				heap.Push(h, x)
			} else {
				if h.Len() > 0 {
					ans += int64(heap.Pop(h).(int))
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
