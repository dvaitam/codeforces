package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveCase(n int, arr []int) int64 {
	h := &MinHeap{}
	heap.Init(h)
	posOdd := 0
	var score int64

	for i, v := range arr {
		if i%2 == 0 { // 1-based odd index
			if v >= 0 {
				score += int64(v)
				posOdd++
			} else {
				heap.Push(h, -v)
			}
		} else {
			if v > 0 {
				if posOdd > 0 {
					score += int64(v)
					posOdd--
				} else if h.Len() > 0 && v > (*h)[0] {
					score += int64(v - (*h)[0])
					heap.Pop(h)
				}
			}
		}
	}
	return score
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
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		fmt.Fprintln(out, solveCase(n, arr))
	}
}
