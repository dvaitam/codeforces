package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		groups := make([][]int64, n+1)
		for i := 0; i < n; i++ {
			var m int
			var p int64
			fmt.Fscan(in, &m, &p)
			if m > n {
				m = n
			}
			groups[m] = append(groups[m], p)
		}
		for i := 0; i <= n; i++ {
			sort.Slice(groups[i], func(a, b int) bool { return groups[i][a] < groups[i][b] })
		}
		h := &IntHeap{}
		heap.Init(h)
		var cost int64
		for i := n; i >= 0; i-- {
			for _, v := range groups[i] {
				heap.Push(h, v)
			}
			for h.Len() > n-i {
				cost += heap.Pop(h).(int64)
			}
		}
		fmt.Fprintln(out, cost)
	}
}
