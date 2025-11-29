package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// IntHeap is a min-heap of ints.
type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
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
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	h := &IntHeap{}
	heap.Init(h)

	var s int64
	currN := n

	for i := n; i >= 1; i-- {
		if a[i] == -1 {
			break
		}
		heap.Push(h, a[i])
		if i == currN {
			currN >>= 1
			if h.Len() > 0 {
				val := heap.Pop(h).(int64)
				s += val
			}
		}
	}
	fmt.Fprintln(writer, s)
}