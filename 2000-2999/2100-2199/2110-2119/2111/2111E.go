package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		var s string
		fmt.Fscan(in, &s)
		letters := []byte(s)

		heaps := [3]*IntHeap{}
		for i := 0; i < 3; i++ {
			h := &IntHeap{}
			heap.Init(h)
			heaps[i] = h
		}
		for i := 0; i < n; i++ {
			if letters[i] == 'a' {
				heap.Push(heaps[0], i)
			} else if letters[i] == 'b' {
				heap.Push(heaps[1], i)
			} else {
				heap.Push(heaps[2], i)
			}
		}

		for ; q > 0; q-- {
			var xStr, yStr string
			fmt.Fscan(in, &xStr, &yStr)
			x := int(xStr[0] - 'a')
			y := int(yStr[0] - 'a')
			if y >= x {
				continue
			}
			h := heaps[x]
			if h.Len() == 0 {
				continue
			}
			idx := heap.Pop(h).(int)
			letters[idx] = byte('a' + y)
			heap.Push(heaps[y], idx)
		}

		fmt.Fprintln(out, string(letters))
	}
}
