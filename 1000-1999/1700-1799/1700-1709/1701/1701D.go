package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l, r int
	idx  int
}

type minHeap []interval

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].r < h[j].r }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(interval)) }
func (h *minHeap) Pop() interface{} {
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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		segs := make([]interval, 0, n)
		for i := 1; i <= n; i++ {
			bi := b[i-1]
			var L, R int
			if bi == 0 {
				L = i + 1
				R = n
			} else {
				L = i/(bi+1) + 1
				R = i / bi
			}
			if L < 1 {
				L = 1
			}
			if R > n {
				R = n
			}
			segs = append(segs, interval{l: L, r: R, idx: i - 1})
		}
		sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
		pq := &minHeap{}
		heap.Init(pq)
		ans := make([]int, n)
		p := 0
		for val := 1; val <= n; val++ {
			for p < n && segs[p].l <= val {
				heap.Push(pq, segs[p])
				p++
			}
			if pq.Len() == 0 {
				// should not happen due to problem constraints
				continue
			}
			it := heap.Pop(pq).(interval)
			ans[it.idx] = val
		}
		for i, v := range ans {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
