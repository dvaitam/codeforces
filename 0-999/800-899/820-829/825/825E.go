package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	rev := make([][]int, n+1)
	outdeg := make([]int, n+1)

	for i := 0; i < m; i++ {
		var v, u int
		fmt.Fscan(reader, &v, &u)
		rev[u] = append(rev[u], v)
		outdeg[v]++
	}

	pq := &IntHeap{}
	heap.Init(pq)
	for i := 1; i <= n; i++ {
		if outdeg[i] == 0 {
			heap.Push(pq, i)
		}
	}

	labels := make([]int, n+1)
	curr := n
	for pq.Len() > 0 {
		v := heap.Pop(pq).(int)
		labels[v] = curr
		curr--
		for _, x := range rev[v] {
			outdeg[x]--
			if outdeg[x] == 0 {
				heap.Push(pq, x)
			}
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, labels[i])
	}
	fmt.Fprintln(writer)
}
