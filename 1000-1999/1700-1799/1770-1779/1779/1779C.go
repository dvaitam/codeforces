package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solve(n, m int, a []int) int {
	ops := 0
	// process right side (indices m..n-1) => original positions m+1..n
	cur := 0
	right := &MaxHeap{}
	heap.Init(right)
	for i := m; i < n; i++ {
		v := a[i]
		if v < 0 {
			heap.Push(right, -v)
		}
		cur += v
		if cur < 0 {
			x := heap.Pop(right).(int)
			cur += 2 * x
			ops++
		}
	}
	// process left side (indices m-1 down to 1) => positions m..2
	cur = 0
	left := &MaxHeap{}
	heap.Init(left)
	for i := m - 1; i > 0; i-- {
		v := a[i]
		if v > 0 {
			heap.Push(left, v)
		}
		cur += v
		if cur > 0 {
			x := heap.Pop(left).(int)
			cur -= 2 * x
			ops++
		}
	}
	return ops
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, solve(n, m, arr))
	}
}
