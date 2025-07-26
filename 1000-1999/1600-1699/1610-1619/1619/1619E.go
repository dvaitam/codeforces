package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

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
		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x <= n {
				freq[x]++
			}
		}

		fill := make([]int, n)
		for i := range fill {
			fill[i] = -1
		}
		h := &MaxHeap{}
		heap.Init(h)
		cost := 0
		for i := 0; i < n; i++ {
			for c := 0; c < freq[i]; c++ {
				heap.Push(h, i)
			}
			if h.Len() == 0 {
				break
			}
			val := heap.Pop(h).(int)
			cost += i - val
			fill[i] = cost
		}

		ans := make([]int, n+1)
		ans[0] = freq[0]
		impossible := false
		for i := 1; i <= n; i++ {
			if impossible || fill[i-1] == -1 {
				ans[i] = -1
				impossible = true
			} else {
				ans[i] = fill[i-1] + freq[i]
			}
		}

		for i := 0; i <= n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
