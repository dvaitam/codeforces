package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type MaxHeap struct{ sort.IntSlice }

func (h MaxHeap) Less(i, j int) bool  { return h.IntSlice[i] > h.IntSlice[j] }
func (h *MaxHeap) Push(x interface{}) { h.IntSlice = append(h.IntSlice, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := h.IntSlice
	x := old[len(old)-1]
	h.IntSlice = old[:len(old)-1]
	return x
}
func (h *MaxHeap) Peek() int { return h.IntSlice[0] }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		h := &MaxHeap{}
		heap.Init(h)
		var ans int64
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x == 0 {
				if h.Len() > 0 {
					ans += int64(h.Peek())
					heap.Pop(h)
				}
			} else {
				heap.Push(h, x)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
