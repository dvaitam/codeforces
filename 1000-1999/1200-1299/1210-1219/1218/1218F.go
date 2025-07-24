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
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	x := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}
	var a int64
	fmt.Fscan(in, &a)
	c := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}

	req := make([]int, n)
	var maxX int64
	for i := 0; i < n; i++ {
		if x[i] > maxX {
			maxX = x[i]
		}
		needed := int64(0)
		if maxX > k {
			needed = (maxX - k + a - 1) / a
		}
		req[i] = int(needed)
		if req[i] > i+1 {
			fmt.Println(-1)
			return
		}
	}

	h := &IntHeap{}
	heap.Init(h)
	selected := 0
	var total int64
	for i := 0; i < n; i++ {
		heap.Push(h, c[i])
		for selected < req[i] {
			if h.Len() == 0 {
				fmt.Println(-1)
				return
			}
			total += heap.Pop(h).(int64)
			selected++
		}
	}

	fmt.Println(total)
}
