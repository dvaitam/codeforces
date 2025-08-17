package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type maxHeap []int

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solve(a []int) int {
	n := len(a)
	freq := make([]int, n+1)
	for _, v := range a {
		freq[v]++
	}
	h := &maxHeap{}
	heap.Init(h)
	s := 0
	sum := 0
	for v := 1; v <= n; v++ {
		f := freq[v]
		if f > 0 {
			heap.Push(h, f)
			sum += f
			for sum > s {
				big := heap.Pop(h).(int)
				sum -= big
				s++
			}
		}
	}
	return s
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
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(reader, &a[i])
		}
		fmt.Fprintln(writer, solve(a))
	}
}
