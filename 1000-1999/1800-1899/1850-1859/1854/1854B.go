package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// IntHeap is a min-heap for ints.
type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	h := &IntHeap{}
	heap.Init(h)
	var invest, sumTotal, ans int64

	for i := 1; i <= n; i++ {
		for invest < int64(i-1) {
			if h.Len() == 0 {
				fmt.Fprintln(writer, ans)
				return
			}
			invest += int64(heap.Pop(h).(int))
		}
		if invest < int64(i-1) {
			fmt.Fprintln(writer, ans)
			return
		}
		heap.Push(h, a[i-1])
		sumTotal += int64(a[i-1])
		if cur := sumTotal - invest; cur > ans {
			ans = cur
		}
	}

	fmt.Fprintln(writer, ans)
}
