package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

// IntHeap is a max-heap of ints.
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

// Node represents a trie node.
type Node struct {
	ch     [2]*Node
	isWord bool
}

var (
	reader *bufio.Reader
	writer *bufio.Writer
	N      int
	root   *Node
)

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &N)
	root = &Node{}
	for i := 0; i < N; i++ {
		var x int
		fmt.Fscan(reader, &x)
		bits := make([]int, 0)
		for x > 0 {
			bits = append(bits, x&1)
			x >>= 1
		}
		u := root
		for j := len(bits) - 1; j >= 0; j-- {
			b := bits[j]
			if u.ch[b] == nil {
				u.ch[b] = &Node{}
			}
			u = u.ch[b]
		}
		u.isWord = true
	}
	h := dfs(root, 0, true)
	// Extract N values
	ans := make([]int, N)
	for i := 0; i < N; i++ {
		ans[i] = heap.Pop(h).(int)
	}
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		writer.WriteString(strconv.Itoa(v))
	}
	writer.WriteByte('\n')
}

// dfs processes the trie and returns a heap of values.
// isRoot indicates if node is the root.
func dfs(u *Node, val int, isRoot bool) *IntHeap {
	if u == nil {
		return nil
	}
	left := dfs(u.ch[0], val<<1, false)
	right := dfs(u.ch[1], (val<<1)|1, false)
	var cur *IntHeap
	switch {
	case left == nil && right == nil:
		cur = &IntHeap{}
		heap.Init(cur)
	case left != nil && right == nil:
		cur = left
	case left == nil && right != nil:
		cur = right
	default:
		// both non-nil, merge smaller into larger
		if left.Len() > right.Len() {
			left, right = right, left
		}
		for left.Len() > 0 {
			x := heap.Pop(left).(int)
			heap.Push(right, x)
		}
		cur = right
	}
	if u.isWord {
		heap.Push(cur, val)
	} else if !isRoot {
		if cur.Len() > 0 {
			heap.Pop(cur)
		}
		heap.Push(cur, val)
	}
	return cur
}
