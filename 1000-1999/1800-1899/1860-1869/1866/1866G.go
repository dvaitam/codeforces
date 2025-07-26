package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type Segment struct {
	L, R int
	A    int64
}

type Item struct {
	R      int
	remain int64
}

type MinHeap []Item

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].R < h[j].R }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

var (
	n     int
	segs  []Segment
	total int64
)

func feasible(z int64) bool {
	idx := 0
	pq := &MinHeap{}
	heap.Init(pq)
	for pos := 1; pos <= n; pos++ {
		for idx < n && segs[idx].L <= pos {
			if segs[idx].A > 0 {
				heap.Push(pq, Item{R: segs[idx].R, remain: segs[idx].A})
			}
			idx++
		}
		for pq.Len() > 0 && (*pq)[0].R < pos {
			return false
		}
		cap := z
		for cap > 0 && pq.Len() > 0 {
			top := &(*pq)[0]
			if top.R < pos {
				return false
			}
			if top.remain <= cap {
				cap -= top.remain
				heap.Pop(pq)
			} else {
				top.remain -= cap
				cap = 0
			}
		}
	}
	return pq.Len() == 0 && idx == n
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n)
	A := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
		total += A[i]
	}
	D := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &D[i])
	}
	segs = make([]Segment, n)
	for i := 0; i < n; i++ {
		L := i + 1 - D[i]
		if L < 1 {
			L = 1
		}
		R := i + 1 + D[i]
		if R > n {
			R = n
		}
		segs[i] = Segment{L: L, R: R, A: A[i]}
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].L < segs[j].L })

	low, high := int64(0), total
	for low < high {
		mid := (low + high) / 2
		if feasible(mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	fmt.Println(low)
}
