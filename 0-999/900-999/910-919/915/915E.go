package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	left, right *Node
	sum         int64
	lazy        int8 // -1 none, 0 or 1 set
}

func newNode() *Node {
	return &Node{lazy: -1}
}

func push(node *Node, l, r int64) {
	if node.lazy == -1 || l == r {
		return
	}
	mid := (l + r) >> 1
	if node.left == nil {
		node.left = newNode()
	}
	if node.right == nil {
		node.right = newNode()
	}
	val := int64(node.lazy)
	node.left.sum = (mid - l + 1) * val
	node.left.lazy = node.lazy
	node.right.sum = (r - mid) * val
	node.right.lazy = node.lazy
	node.lazy = -1
}

func update(node *Node, l, r, ql, qr int64, val int8) {
	if ql > r || qr < l || node == nil {
		return
	}
	if ql <= l && r <= qr {
		node.sum = (r - l + 1) * int64(val)
		node.lazy = val
		node.left, node.right = nil, nil
		return
	}
	push(node, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		if node.left == nil {
			node.left = newNode()
		}
		update(node.left, l, mid, ql, qr, val)
	}
	if qr > mid {
		if node.right == nil {
			node.right = newNode()
		}
		update(node.right, mid+1, r, ql, qr, val)
	}
	node.sum = 0
	if node.left != nil {
		node.sum += node.left.sum
	}
	if node.right != nil {
		node.sum += node.right.sum
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int64
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	root := &Node{sum: n, lazy: 1}

	for i := int64(0); i < q; i++ {
		var l, r, k int64
		fmt.Fscan(reader, &l, &r, &k)
		val := int8(0)
		if k == 2 {
			val = 1
		}
		update(root, 1, n, l, r, val)
		fmt.Fprintln(writer, root.sum)
	}
}
