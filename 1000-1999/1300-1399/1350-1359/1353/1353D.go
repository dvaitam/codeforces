package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type segment struct{ l, r int }

type segHeap []segment

func (h segHeap) Len() int { return len(h) }
func (h segHeap) Less(i, j int) bool {
	lenI := h[i].r - h[i].l + 1
	lenJ := h[j].r - h[j].l + 1
	if lenI == lenJ {
		return h[i].l < h[j].l
	}
	return lenI > lenJ
}
func (h segHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *segHeap) Push(x interface{}) { *h = append(*h, x.(segment)) }
func (h *segHeap) Pop() interface{} {
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
		ans := make([]int, n+1)
		h := &segHeap{{1, n}}
		heap.Init(h)
		for i := 1; i <= n; i++ {
			seg := heap.Pop(h).(segment)
			l, r := seg.l, seg.r
			var mid int
			length := r - l + 1
			if length%2 == 1 {
				mid = (l + r) / 2
			} else {
				mid = (l + r - 1) / 2
			}
			ans[mid] = i
			if l <= mid-1 {
				heap.Push(h, segment{l, mid - 1})
			}
			if mid+1 <= r {
				heap.Push(h, segment{mid + 1, r})
			}
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
