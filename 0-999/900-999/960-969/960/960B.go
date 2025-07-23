package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// An intHeap is a max-heap of ints.
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k1, k2 int
	if _, err := fmt.Fscan(reader, &n, &k1, &k2); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	k := k1 + k2
	h := &intHeap{}
	for i := 0; i < n; i++ {
		diff := a[i] - b[i]
		if diff < 0 {
			diff = -diff
		}
		*h = append(*h, diff)
	}
	heap.Init(h)
	for k > 0 {
		x := heap.Pop(h).(int)
		if x > 0 {
			x--
		} else {
			x = 1
		}
		heap.Push(h, x)
		k--
	}

	var ans int64
	for h.Len() > 0 {
		x := heap.Pop(h).(int)
		ans += int64(x) * int64(x)
	}
	fmt.Fprintln(writer, ans)
}
