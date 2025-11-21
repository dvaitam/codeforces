package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l      int
	r      int
	center int
}

type candidate struct {
	r      int
	center int
}

type maxHeap []candidate

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].r > h[j].r }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(candidate)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
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
		e := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &e[i])
		}

		intervals := make([]interval, 0, n)
		for idx := 0; idx < n; idx++ {
			if e[idx] == 0 {
				continue
			}
			pos := idx + 1
			l := pos - e[idx] + 1
			if l < 1 {
				l = 1
			}
			r := pos + e[idx] - 1
			if r > n {
				r = n
			}
			if l <= r {
				intervals = append(intervals, interval{l: l, r: r, center: pos})
			}
		}

		sort.Slice(intervals, func(i, j int) bool {
			if intervals[i].l == intervals[j].l {
				return intervals[i].r > intervals[j].r
			}
			return intervals[i].l < intervals[j].l
		})

		ans := make([]int, 0)
		cur := 0
		idx := 0
		pq := &maxHeap{}
		heap.Init(pq)
		possible := true

		for cur < n {
			for idx < len(intervals) && intervals[idx].l <= cur+1 {
				heap.Push(pq, candidate{r: intervals[idx].r, center: intervals[idx].center})
				idx++
			}
			for pq.Len() > 0 && (*pq)[0].center <= cur {
				heap.Pop(pq)
			}
			if pq.Len() == 0 {
				possible = false
				break
			}
			best := heap.Pop(pq).(candidate)
			ans = append(ans, best.center)
			cur = best.r
		}

		if !possible {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprintln(out, len(ans))
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
