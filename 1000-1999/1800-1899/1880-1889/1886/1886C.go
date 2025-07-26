package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt.
// Given a string s1, subsequent strings are formed by repeatedly
// removing one character so that the resulting string is
// lexicographically minimal. Concatenating all these strings yields S.
// For a given position pos we output the character S_pos.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		var pos int64
		fmt.Fscan(in, &pos)

		n := len(s)
		prefix := func(k int64) int64 {
			return k * (int64(2*n) - k + 1) / 2
		}
		low, high := 1, n
		for low < high {
			mid := (low + high) / 2
			if prefix(int64(mid)) >= pos {
				high = mid
			} else {
				low = mid + 1
			}
		}
		k := low
		pos -= prefix(int64(k - 1))

		bytes := []byte(s)
		next := make([]int, n)
		prev := make([]int, n)
		for i := 0; i < n; i++ {
			next[i] = i + 1
			prev[i] = i - 1
		}
		next[n-1] = -1

		removed := make([]bool, n)
		h := &intHeap{}
		heap.Init(h)
		for i := 0; i+1 < n; i++ {
			if bytes[i] > bytes[i+1] {
				heap.Push(h, i)
			}
		}
		head, tail := 0, n-1

		for step := 0; step < k-1; step++ {
			j := -1
			for h.Len() > 0 {
				top := heap.Pop(h).(int)
				if removed[top] {
					continue
				}
				nxt := next[top]
				if nxt != -1 && bytes[top] > bytes[nxt] {
					j = top
					break
				}

			}
			if j == -1 {
				j = tail
			}
			removed[j] = true
			p, q := prev[j], next[j]
			if p != -1 {
				next[p] = q
			} else {
				head = q
			}
			if q != -1 {
				prev[q] = p
			} else {
				tail = p
			}
			if p != -1 && q != -1 && bytes[p] > bytes[q] {
				heap.Push(h, p)
			}
		}

		idx := head
		for i := int64(1); i < pos; i++ {
			idx = next[idx]
		}
		fmt.Fprintf(out, "%c", bytes[idx])
	}
}

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
