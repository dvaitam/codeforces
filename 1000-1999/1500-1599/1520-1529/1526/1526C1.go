package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	h := &IntHeap{}
	heap.Init(h)
	var sum int64
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		sum += int64(x)
		heap.Push(h, x)
		if sum < 0 {
			smallest := heap.Pop(h).(int)
			sum -= int64(smallest)
		}
	}
	fmt.Fprintln(writer, h.Len())
}
