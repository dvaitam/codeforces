package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const inf = 1000000000

// Pair holds a pair of ints for heap
type Pair struct{ first, second int }

// MaxHeap implements a max-heap of Pair
type MaxHeap []Pair

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool {
	if h[i].first != h[j].first {
		return h[i].first > h[j].first
	}
	return h[i].second > h[j].second
}
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(Pair))
}
func (h *MaxHeap) Pop() interface{} {
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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	a[n+1] = inf
	dp := make([]int, n+2)
	prv := make([]int, n+2)
	ans := make([]int, n+2)
	// init dp to -1
	for i := 1; i <= n+1; i++ {
		dp[i] = -1
	}
	dp[1] = inf
	// heap of Pair
	h := &MaxHeap{}
	heap.Init(h)
	heap.Push(h, Pair{1, 0})
	l := 1
	for i := 2; i <= n+1; i++ {
		mxFirst, mxSecond := -1, 0
		if a[i] > a[i-1] {
			mxFirst, mxSecond = dp[i-1], i-1
		}
		// pop outdated
		for h.Len() > 0 && (*h)[0].second < l-1 {
			heap.Pop(h)
		}
		if h.Len() > 0 && -(*h)[0].first < a[i] {
			candFirst, candSecond := a[i-1], (*h)[0].second
			if candFirst > mxFirst || (candFirst == mxFirst && candSecond > mxSecond) {
				mxFirst, mxSecond = candFirst, candSecond
			}
		}
		dp[i] = mxFirst
		prv[i] = mxSecond
		if dp[i-1] > a[i] {
			heap.Push(h, Pair{-a[i-1], i - 1})
		}
		if a[i] >= a[i-1] {
			l = i
		}
	}
	cur := 0
	if dp[n] != -1 {
		fmt.Fprintln(out, "YES")
		cur = n
	} else if dp[n+1] != -1 {
		fmt.Fprintln(out, "YES")
		cur = n + 1
	} else {
		fmt.Fprintln(out, "NO")
		return
	}
	for cur > 0 {
		ans[cur] = 1
		cur = prv[cur]
	}
	// output sequence: 1-ans[i]
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprint(1 - ans[i]))
	}
	out.WriteByte('\n')
}
