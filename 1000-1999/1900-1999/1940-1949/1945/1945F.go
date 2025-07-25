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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		v := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &v[i])
		}
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
		}
		valByPos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			valByPos[i] = v[p[i]]
		}

		h := &IntHeap{}
		heap.Init(h)
		bestStrength := int64(0)
		bestCount := 0
		for r := n; r >= 1; r-- {
			heap.Push(h, valByPos[r])
			if h.Len() > r {
				heap.Pop(h)
			}
			if h.Len() == r {
				minVal := (*h)[0]
				strength := int64(r) * int64(minVal)
				if strength > bestStrength || (strength == bestStrength && r < bestCount) {
					bestStrength = strength
					bestCount = r
				}
			}
		}
		fmt.Fprintln(out, bestStrength, bestCount)
	}
}
