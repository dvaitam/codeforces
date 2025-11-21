package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		count := make(map[int]int)
		sumA, sumB := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			count[x]++
			sumA += x
		}
		h := &MaxHeap{}
		heap.Init(h)
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(in, &x)
			heap.Push(h, x)
			sumB += x
		}
		if sumA != sumB {
			fmt.Fprintln(out, "No")
			continue
		}
		possible := true
		for h.Len() > 0 {
			val := heap.Pop(h).(int)
			if count[val] > 0 {
				count[val]--
				continue
			}
			if val == 1 {
				possible = false
				break
			}
			heap.Push(h, val/2)
			heap.Push(h, val-val/2)
		}
		if possible {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
