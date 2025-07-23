package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// IntHeap implements a max-heap for ints.
type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func merge(a, b *IntHeap) *IntHeap {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.Len() < b.Len() {
		a, b = b, a
	}
	for b.Len() > 0 {
		heap.Push(a, heap.Pop(b))
	}
	return a
}

type Node struct {
	child [26]int
	end   bool
	depth int
}

func solve(words []string) int {
	nodes := []Node{{}}
	// build trie
	for _, s := range words {
		idx := 0
		for _, ch := range s {
			c := int(ch - 'a')
			if nodes[idx].child[c] == 0 {
				nodes[idx].child[c] = len(nodes)
				nodes = append(nodes, Node{depth: nodes[idx].depth + 1})
			}
			idx = nodes[idx].child[c]
		}
		nodes[idx].end = true
	}
	// build processing order (preorder)
	order := []int{0}
	for i := 0; i < len(order); i++ {
		u := order[i]
		for _, v := range nodes[u].child {
			if v != 0 {
				order = append(order, v)
			}
		}
	}
	heaps := make([]*IntHeap, len(nodes))
	// process in reverse order
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		var h *IntHeap
		for _, v := range nodes[u].child {
			if v != 0 {
				h = merge(h, heaps[v])
			}
		}
		if h == nil {
			h = &IntHeap{}
		}
		if nodes[u].end {
			heap.Push(h, nodes[u].depth)
		} else if u != 0 && h.Len() > 0 && (*h)[0] > nodes[u].depth {
			heap.Pop(h)
			heap.Push(h, nodes[u].depth)
		}
		heaps[u] = h
	}
	ans := 0
	for _, v := range *heaps[0] {
		ans += v
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}
	fmt.Fprintln(writer, solve(words))
}
