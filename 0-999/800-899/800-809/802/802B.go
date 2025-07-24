package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Item struct {
	next int
	book int
}

type MaxHeap []Item

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].next > h[j].next }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	next := make([]int, n)
	last := make(map[int]int)
	for i := 0; i < n; i++ {
		last[a[i]] = n
	}
	for i := n - 1; i >= 0; i-- {
		next[i] = last[a[i]]
		last[a[i]] = i
	}

	library := make(map[int]int)
	pq := &MaxHeap{}
	heap.Init(pq)
	cost := 0

	for i := 0; i < n; i++ {
		b := a[i]
		nxt := next[i]
		if _, ok := library[b]; ok {
			library[b] = nxt
			heap.Push(pq, Item{nxt, b})
			continue
		}
		cost++
		if len(library) >= k {
			for len(*pq) > 0 {
				item := heap.Pop(pq).(Item)
				if cur, ok := library[item.book]; ok && cur == item.next {
					delete(library, item.book)
					break
				}
			}
		}
		if k > 0 {
			library[b] = nxt
			heap.Push(pq, Item{nxt, b})
		}
	}

	fmt.Fprintln(out, cost)
}
