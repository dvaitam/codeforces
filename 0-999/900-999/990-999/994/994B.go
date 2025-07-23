package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Knight struct {
	p   int
	c   int64
	idx int
}

type MinHeap []int64

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}
func (h *MinHeap) Pop() interface{} {
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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	knights := make([]Knight, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &knights[i].p)
		knights[i].idx = i
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &knights[i].c)
	}

	sort.Slice(knights, func(i, j int) bool {
		return knights[i].p < knights[j].p
	})

	ans := make([]int64, n)
	h := &MinHeap{}
	heap.Init(h)
	var sum int64
	for _, kn := range knights {
		ans[kn.idx] = sum + kn.c
		heap.Push(h, kn.c)
		sum += kn.c
		if h.Len() > k {
			v := heap.Pop(h).(int64)
			sum -= v
		}
	}

	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
