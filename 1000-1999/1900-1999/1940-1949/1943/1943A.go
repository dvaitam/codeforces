package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Node struct {
	cnt int
	val int
}

type MinHeap []Node

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	if h[i].cnt == h[j].cnt {
		return h[i].val < h[j].val
	}
	return h[i].cnt < h[j].cnt
}
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Node)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	freq := make(map[int]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		freq[x]++
	}
	h := &MinHeap{}
	heap.Init(h)
	for v, c := range freq {
		heap.Push(h, Node{cnt: c, val: v})
	}
	collected := make(map[int]bool)
	turn := 0
	for h.Len() > 0 {
		node := heap.Pop(h).(Node)
		if turn%2 == 0 { // Alice picks
			collected[node.val] = true
		} else { // Bob removes one occurrence
			node.cnt--
			if node.cnt > 0 {
				heap.Push(h, node)
			}
		}
		turn++
	}
	mex := 0
	for collected[mex] {
		mex++
	}
	fmt.Fprintln(writer, mex)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t++ {
		solve(reader, writer)
	}
}
