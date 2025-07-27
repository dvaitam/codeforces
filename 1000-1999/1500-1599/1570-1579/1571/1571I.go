package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type node struct {
	r   int
	idx int
}

type minHeap []node

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
	if h[i].r == h[j].r {
		return h[i].idx < h[j].idx
	}
	return h[i].r < h[j].r
}
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(node)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func findStart(L, R []int) (int, bool) {
	n := len(L)
	Ls := make([]int, n)
	copy(Ls, L)
	Rs := make([]int, n)
	copy(Rs, R)
	sort.Ints(Ls)
	sort.Ints(Rs)
	lo := -1 << 60
	hi := 1 << 60
	for i := 0; i < n; i++ {
		if Ls[i]-i > lo {
			lo = Ls[i] - i
		}
		if Rs[i]-i < hi {
			hi = Rs[i] - i
		}
	}
	if lo > hi {
		return 0, false
	}
	return lo, true
}

func scheduleWithX(L, R []int, x int) ([]int, bool) {
	n := len(L)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return L[idx[i]] < L[idx[j]] })
	pq := &minHeap{}
	heap.Init(pq)
	ptr := 0
	order := make([]int, 0, n)
	for t := x; t < x+n; t++ {
		for ptr < n && L[idx[ptr]] <= t {
			heap.Push(pq, node{R[idx[ptr]], idx[ptr]})
			ptr++
		}
		if pq.Len() == 0 {
			return nil, false
		}
		nd := heap.Pop(pq).(node)
		if nd.r < t {
			return nil, false
		}
		order = append(order, nd.idx)
	}
	return order, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		L := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &L[i])
		}
		R := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &R[i])
		}
		x, ok := findStart(L, R)
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		order, ok := scheduleWithX(L, R, x)
		if !ok {
			// try using upper bound maybe
			upper := x
			Ls := make([]int, len(L))
			copy(Ls, L)
			sort.Ints(Ls)
			Rs := make([]int, len(R))
			copy(Rs, R)
			sort.Ints(Rs)
			hi := 1 << 60
			for i := 0; i < n; i++ {
				if Rs[i]-i < hi {
					hi = Rs[i] - i
				}
			}
			if hi != upper {
				order, ok = scheduleWithX(L, R, hi)
				if !ok {
					fmt.Fprintln(out, -1)
					continue
				}
				x = hi
			} else {
				fmt.Fprintln(out, -1)
				continue
			}
		}
		fmt.Fprintln(out, x)
		for i, idx := range order {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, idx+1)
		}
		fmt.Fprintln(out)
	}
}
