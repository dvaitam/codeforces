package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MinHeap struct {
	data []int
	s    []int64
}

func (h MinHeap) Len() int { return len(h.data) }
func (h MinHeap) Less(i, j int) bool {
	return h.s[h.data[i]] < h.s[h.data[j]]
}
func (h MinHeap) Swap(i, j int) { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *MinHeap) Push(x interface{}) {
	h.data = append(h.data, x.(int))
}
func (h *MinHeap) Pop() interface{} {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

var (
	n, m, p, q int
	parent     []int
	ssum       []int64
)

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func main() {
	rdr := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	// read n, m, p, q_input
	var qInput int
	if _, err := fmt.Fscan(rdr, &n, &m, &p, &qInput); err != nil {
		return
	}
	q = n - qInput
	parent = make([]int, n+1)
	ssum = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		ssum[i] = 0
	}
	// initial unions
	for i := 0; i < m; i++ {
		var a, b int
		var d int64
		fmt.Fscan(rdr, &a, &b, &d)
		x := find(a)
		y := find(b)
		if x != y {
			q--
			ssum[x] += ssum[y] + d
			parent[y] = x
		} else {
			ssum[x] += d
		}
	}
	// check feasibility
	if q < 0 || q > p || (q == 0 && m == 0 && p > 0) {
		fmt.Fprintln(w, "NO")
		return
	}
	fmt.Fprintln(w, "YES")
	// build heap of roots
	h := &MinHeap{data: make([]int, 0, n), s: ssum}
	for i := 1; i <= n; i++ {
		if find(i) == i {
			h.data = append(h.data, i)
		}
	}
	heap.Init(h)
	// primary merges
	for q > 0 {
		if h.Len() < 2 {
			break
		}
		x := heap.Pop(h).(int)
		y := heap.Pop(h).(int)
		fmt.Fprintf(w, "%d %d\n", x, y)
		// update ssum and union
		sxy := ssum[x] + ssum[y]
		inc := sxy + 1
		if inc > 1e9 {
			inc = 1e9
		}
		ssum[x] = sxy + inc
		parent[y] = x
		heap.Push(h, x)
		q--
		p--
	}
	// extra merges if p remains
	if p > 0 {
		// find any non-root
		for i := 1; i <= n; i++ {
			if find(i) != i {
				r := find(i)
				for j := 0; j < p; j++ {
					fmt.Fprintf(w, "%d %d\n", r, i)
				}
				break
			}
		}
	}
}
