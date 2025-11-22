package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const inf = int(1e9 + 7)

type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type node struct {
	left, right *node
	sz          int
	rev         bool
	h           minHeap
	pri         int
}

func newNode() *node {
	return &node{
		sz:  1,
		pri: rand.Int(),
	}
}

func size(nd *node) int {
	if nd == nil {
		return 0
	}
	return nd.sz
}

func topValid(nd *node, deleted []bool) int {
	if nd == nil {
		return inf
	}
	for nd.h.Len() > 0 {
		v := nd.h[0]
		if deleted[v] {
			heap.Pop(&nd.h)
		} else {
			return v
		}
	}
	return inf
}

func pull(nd *node) {
	nd.sz = 1 + size(nd.left) + size(nd.right)
}

func push(nd *node) {
	if nd != nil && nd.rev {
		nd.left, nd.right = nd.right, nd.left
		if nd.left != nil {
			nd.left.rev = !nd.left.rev
		}
		if nd.right != nil {
			nd.right.rev = !nd.right.rev
		}
		nd.rev = false
	}
}

// split root into [0..k-1], [k..]
func split(nd *node, k int) (*node, *node) {
	if nd == nil {
		return nil, nil
	}
	push(nd)
	if size(nd.left) >= k {
		l, r := split(nd.left, k)
		nd.left = r
		pull(nd)
		return l, nd
	}
	l, r := split(nd.right, k-size(nd.left)-1)
	nd.right = l
	pull(nd)
	return nd, r
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pri < b.pri {
		push(a)
		a.right = merge(a.right, b)
		pull(a)
		return a
	}
	push(b)
	b.left = merge(a, b.left)
	pull(b)
	return b
}

func addValue(nd *node, val int) {
	heap.Push(&nd.h, val)
}

func queryKth(nd *node, k int, deleted []bool) int {
	best := inf
	for nd != nil {
		push(nd)
		if v := topValid(nd, deleted); v < best {
			best = v
		}
		ls := size(nd.left)
		if k <= ls {
			nd = nd.left
		} else if k == ls+1 {
			break
		} else {
			k -= ls + 1
			nd = nd.right
		}
	}
	if best == inf {
		return 0
	}
	return best
}

func build(n int) *node {
	var root *node
	for i := 0; i < n; i++ {
		root = merge(root, newNode())
	}
	return root
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	rand.Seed(time.Now().UnixNano())

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	root := build(n)
	deleted := make([]bool, q+2)

	ansPrev := 0
	for i := 1; i <= q; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		if a == 1 {
			r := (b+ansPrev-1)%n + 1
			left, right := split(root, r)
			addValue(left, i)
			root = merge(left, right)
		} else if a == 2 {
			r := (b+ansPrev-1)%n + 1
			left, right := split(root, r)
			if left != nil {
				left.rev = !left.rev
			}
			root = merge(left, right)
		} else { // delete
			x := (b+ansPrev-1)%q + 1
			deleted[x] = true
		}

		p := (c+ansPrev-1)%n + 1
		ansPrev = queryKth(root, p, deleted)
		fmt.Fprintln(out, ansPrev)
	}
}
