package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		groups := make([][]int, n+1)
		for i := 0; i < n; i++ {
			var m, p int
			fmt.Fscan(reader, &m, &p)
			if m > n {
				m = n
			}
			groups[m] = append(groups[m], p)
		}
		for i := 0; i <= n; i++ {
			sort.Ints(groups[i])
		}
		pq := &IntHeap{}
		heap.Init(pq)
		bought := 0
		var cost int64
		for need := n; need >= 0; need-- {
			for _, v := range groups[need] {
				heap.Push(pq, v)
			}
			for bought < n-need && pq.Len() > 0 {
				cost += int64(heap.Pop(pq).(int))
				bought++
			}
		}
		fmt.Fprintln(writer, cost)
	}
}
