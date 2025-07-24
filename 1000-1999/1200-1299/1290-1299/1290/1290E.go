package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	key   int
	left  *Node
	right *Node
	size  int
	sum   int64
}

func update(t *Node) {
	if t == nil {
		return
	}
	ls, lsum := 0, int64(0)
	if t.left != nil {
		ls = t.left.size
		lsum = t.left.sum
	}
	rs, rsum := 0, int64(0)
	if t.right != nil {
		rs = t.right.size
		rsum = t.right.sum
	}
	t.size = ls + rs + 1
	t.sum = int64(t.size) + lsum + rsum
}

func split(t *Node, key int) (*Node, *Node) {
	if t == nil {
		return nil, nil
	}
	if t.key <= key {
		var r *Node
		t.right, r = split(t.right, key)
		update(t)
		return t, r
	}
	var l *Node
	l, t.left = split(t.left, key)
	update(t)
	return l, t
}

func newNode(key int) *Node {
	n := &Node{key: key, size: 1, sum: 1}
	return n
}

func nextInt(scanner *bufio.Reader) int {
	sign, val := 1, 0
	c, _ := scanner.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = scanner.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = scanner.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = scanner.ReadByte()
	}
	return sign * val
}

func main() {
	scanner := bufio.NewReader(os.Stdin)
	n := nextInt(scanner)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt(scanner)
	}
	pos := make([]int, n+1)
	for i, v := range a {
		pos[v] = i
	}
	var root *Node
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		p := pos[i]
		left, right := split(root, p-1)
		node := newNode(p)
		node.left = left
		node.right = right
		update(node)
		root = node
		ans[i-1] = root.sum
	}
	writer := bufio.NewWriter(os.Stdout)
	for _, v := range ans {
		fmt.Fprintln(writer, v)
	}
	writer.Flush()
}
