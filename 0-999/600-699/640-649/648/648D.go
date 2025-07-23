package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// Bowl interval in terms of dog coordinates that can reach it.
type Bowl struct {
	l int
	r int
}

// MinHeap stores right endpoints of available bowls.
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	dogs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &dogs[i])
	}

	bowls := make([]Bowl, m)
	for i := 0; i < m; i++ {
		var u, t int
		fmt.Fscan(in, &u, &t)
		bowls[i] = Bowl{l: u - t, r: u + t}
	}

	sort.Ints(dogs)
	sort.Slice(bowls, func(i, j int) bool { return bowls[i].l < bowls[j].l })

	pq := &MinHeap{}
	heap.Init(pq)
	j := 0
	count := 0

	for _, x := range dogs {
		for j < m && bowls[j].l <= x {
			heap.Push(pq, bowls[j].r)
			j++
		}
		for pq.Len() > 0 && (*pq)[0] < x {
			heap.Pop(pq)
		}
		if pq.Len() > 0 {
			heap.Pop(pq)
			count++
		}
	}

	fmt.Fprintln(out, count)
}
