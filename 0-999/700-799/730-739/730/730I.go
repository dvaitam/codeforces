package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// entry represents a heap element with weight and position
type entry struct {
	w, pos int
}

// maxHeap implements a max-heap of entry
type maxHeap []entry

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].w > h[j].w }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
func (h maxHeap) Peek() entry { return h[0] }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, p, s int
	if _, err := fmt.Fscan(reader, &n, &p, &s); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	sum := make([]int, n+1)
	type node struct{ a, b, pos int }
	x := make([]node, n)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
		sum[i] = b[i] - a[i]
		x[i-1] = node{a: a[i], b: b[i], pos: i}
	}
	// sort by a descending
	sort.Slice(x, func(i, j int) bool {
		return x[i].a > x[j].a
	})
	// initialize heaps
	q := &maxHeap{}
	h1 := &maxHeap{}
	h2 := &maxHeap{}
	heap.Init(q)
	heap.Init(h1)
	heap.Init(h2)
	t1 := make([]bool, n+1)
	t2 := make([]bool, n+1)
	used := make([]bool, n+1)
	ans := 0
	// select first p by highest a
	for i := 0; i < p; i++ {
		node := x[i]
		heap.Push(q, entry{w: node.b - node.a, pos: node.pos})
		ans += node.a
		t1[node.pos] = true
	}
	// remaining into h1 and h2
	for i := p; i < n; i++ {
		node := x[i]
		heap.Push(h1, entry{w: node.a, pos: node.pos})
		heap.Push(h2, entry{w: node.b, pos: node.pos})
	}
	// process s selections
	for i := 0; i < s; i++ {
		// get top of h1 not used
		var e1 entry
		for {
			e1 = h1.Peek()
			if used[e1.pos] {
				heap.Pop(h1)
				continue
			}
			break
		}
		// get top of h2 not used
		var e2 entry
		for {
			e2 = h2.Peek()
			if used[e2.pos] {
				heap.Pop(h2)
				continue
			}
			break
		}
		e3 := q.Peek()
		w1, pos1 := e1.w, e1.pos
		w2, pos2 := e2.w, e2.pos
		w3, pos3 := e3.w, e3.pos
		if w1+w3 > w2 {
			heap.Pop(q)
			t1[pos3] = false
			t2[pos3] = true
			heap.Push(q, entry{w: sum[pos1], pos: pos1})
			used[pos1] = true
			t1[pos1] = true
			heap.Pop(h1)
			ans += w1 + w3
		} else {
			heap.Pop(h2)
			t2[pos2] = true
			used[pos2] = true
			ans += w2
		}
	}
	// output
	fmt.Fprintln(writer, ans)
	for i := 1; i <= n; i++ {
		if t1[i] {
			fmt.Fprint(writer, i, " ")
		}
	}
	fmt.Fprintln(writer)
	for i := 1; i <= n; i++ {
		if t2[i] {
			fmt.Fprint(writer, i, " ")
		}
	}
	fmt.Fprintln(writer)
}
