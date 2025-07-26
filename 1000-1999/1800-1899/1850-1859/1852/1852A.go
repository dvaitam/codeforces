package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Int64Heap is a min-heap of int64.
type Int64Heap []int64

func (h Int64Heap) Len() int           { return len(h) }
func (h Int64Heap) Less(i, j int) bool { return h[i] < h[j] }
func (h Int64Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *Int64Heap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}
func (h *Int64Heap) Pop() interface{} {
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if a[0] != 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		if n == 1 {
			fmt.Fprintln(out, int64(k+1))
			continue
		}
		h := &Int64Heap{}
		for i := 1; i < n; i++ {
			*h = append(*h, a[i])
		}
		heap.Init(h)
		x := int64(1)
		step := int64(n)
		for i := 0; i < k; i++ {
			cand := x + 1
			for h.Len() > 0 && cand == (*h)[0] {
				v := heap.Pop(h).(int64)
				heap.Push(h, v+step)
				cand++
			}
			x = cand
		}
		fmt.Fprintln(out, x)
	}
}
