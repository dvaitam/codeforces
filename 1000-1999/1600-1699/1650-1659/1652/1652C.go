package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// IntHeap implements a max heap for int64 numbers.
type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

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
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

		h := &IntHeap{}
		heap.Push(h, sum)
		possible := true
		for i := 0; i < n && possible; i++ {
			w := a[i]
			for {
				if h.Len() == 0 {
					possible = false
					break
				}
				x := (*h)[0]
				if x == w {
					heap.Pop(h)
					break
				}
				if x < w || x == 1 {
					possible = false
					break
				}
				heap.Pop(h)
				left := x / 2
				right := x - left
				heap.Push(h, left)
				heap.Push(h, right)
			}
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
