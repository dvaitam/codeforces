package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type node struct {
	cost int
	rem  int
}

type minHeap []node

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].cost < h[j].cost }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(node)) }
func (h *minHeap) Pop() interface{} {
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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, p int
		fmt.Fscan(in, &n, &p)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		// sort residents by b ascending
		type pair struct{ a, b int }
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{a: a[i], b: b[i]}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].b < arr[j].b })

		idx := 0
		var total int64
		pq := &minHeap{}
		heap.Init(pq)
		for informed := 0; informed < n; informed++ {
			if pq.Len() > 0 && (*pq)[0].cost < p {
				x := heap.Pop(pq).(node)
				total += int64(x.cost)
				if x.rem > 1 {
					x.rem--
					heap.Push(pq, x)
				}
				// inform next resident with smallest b
				cur := arr[idx]
				idx++
				if cur.a > 0 {
					heap.Push(pq, node{cost: cur.b, rem: cur.a})
				}
			} else {
				total += int64(p)
				cur := arr[idx]
				idx++
				if cur.a > 0 {
					heap.Push(pq, node{cost: cur.b, rem: cur.a})
				}
			}
		}
		fmt.Fprintln(out, total)
	}
}
