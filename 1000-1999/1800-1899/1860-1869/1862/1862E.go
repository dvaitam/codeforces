package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// minHeap implements a min-heap for ints.
type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
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
		var d int
		fmt.Fscan(in, &n, &m, &d)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		h := &minHeap{}
		heap.Init(h)
		sumPrev := 0
		best := 0
		for i := 0; i < n; i++ {
			// candidate where day i+1 is the last movie we watch
			cand := a[i] + sumPrev - d*(i+1)
			if cand > best {
				best = cand
			}
			if a[i] > 0 {
				heap.Push(h, a[i])
				sumPrev += a[i]
				if h.Len() > m-1 {
					sumPrev -= heap.Pop(h).(int)
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}
