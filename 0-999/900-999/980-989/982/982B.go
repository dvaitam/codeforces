package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type row struct {
	w   int
	idx int
}

type maxHeap []row

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].w > h[j].w }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(row)) }
func (h *maxHeap) Pop() interface{} {
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
	pairs := make([]row, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &pairs[i].w)
		pairs[i].idx = i + 1
	}
	var s string
	fmt.Fscan(reader, &s)

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].w < pairs[j].w
	})

	h := &maxHeap{}
	heap.Init(h)
	p := 0
	for _, ch := range s {
		if ch == '0' {
			row := pairs[p]
			p++
			heap.Push(h, row)
			fmt.Fprintf(writer, "%d ", row.idx)
		} else {
			row := heap.Pop(h).(row)
			fmt.Fprintf(writer, "%d ", row.idx)
		}
	}
	writer.WriteByte('\n')
}
