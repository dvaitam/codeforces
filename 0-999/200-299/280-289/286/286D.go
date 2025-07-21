package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// MinHeap is a min-heap of int64
type MinHeap []int64

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	walls := make([]struct{ l, r, t int64 }, m)
	ys := make([]int64, 0, 2*m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &walls[i].l, &walls[i].r, &walls[i].t)
		ys = append(ys, walls[i].l, walls[i].r)
	}
	qs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &qs[i])
	}
	// coordinate compress y
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	// remove duplicates
	uniq := ys[:0]
	for i, v := range ys {
		if i == 0 || v != ys[i-1] {
			uniq = append(uniq, v)
		}
	}
	ys = uniq
	// map y to index
	idx := func(y int64) int {
		return sort.Search(len(ys), func(i int) bool { return ys[i] >= y })
	}
	// events at each y index: remove first then add
	adds := make([][]int64, len(ys))
	rems := make([][]int64, len(ys))
	for _, w := range walls {
		li := idx(w.l)
		ri := idx(w.r)
		adds[li] = append(adds[li], w.t)
		rems[ri] = append(rems[ri], w.t)
	}
	// sweep
	active := &MinHeap{}
	heap.Init(active)
	removeCnt := make(map[int64]int)
	type seg struct{ A, B, L int64 }
	segs := make([]seg, 0, len(ys))
	for k := 0; k < len(ys); k++ {
		// process removals
		for _, t := range rems[k] {
			removeCnt[t]++
		}
		// process adds
		for _, t := range adds[k] {
			heap.Push(active, t)
		}
		if k+1 < len(ys) {
			// get current min ti
			for active.Len() > 0 {
				top := (*active)[0]
				if removeCnt[top] > 0 {
					heap.Pop(active)
					removeCnt[top]--
				} else {
					break
				}
			}
			if active.Len() > 0 {
				M := (*active)[0]
				y0 := ys[k]
				y1 := ys[k+1]
				L := y1 - y0
				A := M - y1
				// B = A + L
				B := A + L
				segs = append(segs, seg{A, B, L})
			}
		}
	}
	// prepare events for queries
	// sort A's and B's
	As := make([]int64, len(segs))
	Bs := make([]struct{ B, A, L int64 }, len(segs))
	for i, s := range segs {
		As[i] = s.A
		Bs[i] = struct{ B, A, L int64 }{s.B, s.A, s.L}
	}
	sort.Slice(As, func(i, j int) bool { return As[i] < As[j] })
	sort.Slice(Bs, func(i, j int) bool { return Bs[i].B < Bs[j].B })
	// sort queries
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return qs[order[i]] < qs[order[j]] })
	ans := make([]int64, n)
	pa, pb := 0, 0
	var cntA, cntB int64
	var sumA, sumAB, sumFullL int64
	for _, oi := range order {
		q := qs[oi]
		for pa < len(As) && As[pa] <= q {
			sumA += As[pa]
			cntA++
			pa++
		}
		for pb < len(Bs) && Bs[pb].B <= q {
			sumAB += Bs[pb].A
			sumFullL += Bs[pb].L
			cntB++
			pb++
		}
		// total = full + partial
		// partial count = cntA - cntB, sum of A in partial = sumA - sumAB
		ans[oi] = sumFullL + (cntA-cntB)*q - (sumA - sumAB)
	}
	// output
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
