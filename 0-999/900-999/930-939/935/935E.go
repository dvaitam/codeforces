package main

import (
	"bufio"
	"fmt"
	"os"
)

// Node represents a parsed expression tree node
type Node struct {
	val    int
	left   *Node
	right  *Node
	ops    int
	maxVal []int
	minVal []int
}

// parse parses expression s starting at index *idx and returns the Node
func parse(s string, idx *int) *Node {
	if s[*idx] >= '0' && s[*idx] <= '9' {
		val := int(s[*idx] - '0')
		*idx++
		return &Node{val: val, ops: 0}
	}
	// assume '(' at s[idx]
	*idx++
	left := parse(s, idx)
	*idx++ // skip '?'
	right := parse(s, idx)
	*idx++ // skip ')'
	return &Node{left: left, right: right, ops: left.ops + right.ops + 1}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// solve computes dp arrays for node
func solve(n *Node) {
	if n.ops == 0 {
		n.maxVal = []int{n.val}
		n.minVal = []int{n.val}
		return
	}
	solve(n.left)
	solve(n.right)

	l, r := n.left, n.right
	n.maxVal = make([]int, n.ops+1)
	n.minVal = make([]int, n.ops+1)
	for i := range n.maxVal {
		n.maxVal[i] = -1 << 60
		n.minVal[i] = 1 << 60
	}

	for i := 0; i <= l.ops; i++ {
		for j := 0; j <= r.ops; j++ {
			// plus operator
			if i+j+1 <= n.ops {
				vMax := l.maxVal[i] + r.maxVal[j]
				vMin := l.minVal[i] + r.minVal[j]
				idx := i + j + 1
				if vMax > n.maxVal[idx] {
					n.maxVal[idx] = vMax
				}
				if vMin < n.minVal[idx] {
					n.minVal[idx] = vMin
				}
			}
			// minus operator
			if i+j <= n.ops {
				vMax := l.maxVal[i] - r.minVal[j]
				vMin := l.minVal[i] - r.maxVal[j]
				idx := i + j
				if vMax > n.maxVal[idx] {
					n.maxVal[idx] = vMax
				}
				if vMin < n.minVal[idx] {
					n.minVal[idx] = vMin
				}
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var expr string
	fmt.Fscan(in, &expr)
	var P, M int
	fmt.Fscan(in, &P, &M)
	idx := 0
	root := parse(expr, &idx)
	solve(root)
	fmt.Println(root.maxVal[P])
}
