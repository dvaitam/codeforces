package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Treap node for ordered set of ints
type Node struct {
	key         int
	prio        int
	left, right *Node
}

func split(root *Node, key int) (l, r *Node) {
	if root == nil {
		return nil, nil
	}
	if root.key < key {
		var sr *Node
		root.right, sr = split(root.right, key)
		return root, sr
	}
	var sl *Node
	sl, root.left = split(root.left, key)
	return sl, root
}

func merge(l, r *Node) *Node {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		return l
	}
	r.left = merge(l, r.left)
	return r
}

func insert(root *Node, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.prio > root.prio {
		l, r := split(root, node.key)
		node.left, node.right = l, r
		return node
	}
	if node.key < root.key {
		root.left = insert(root.left, node)
	} else {
		root.right = insert(root.right, node)
	}
	return root
}

func predecessor(root *Node, key int) int {
	res := -1 << 60
	for root != nil {
		if root.key < key {
			if root.key > res {
				res = root.key
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *Node, key int) int {
	res := 1<<60 - 1
	for root != nil {
		if root.key > key {
			if root.key < res {
				res = root.key
			}
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

type Set struct{ root *Node }

func (s *Set) Insert(key int) {
	s.root = insert(s.root, &Node{key: key, prio: rand.Int()})
}

func (s *Set) Prev(key int) int { return predecessor(s.root, key) }
func (s *Set) Next(key int) int { return successor(s.root, key) }

func calc(len, a int) int {
	if len < 0 {
		return 0
	}
	return (len + 1) / (a + 1)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, a int
	if _, err := fmt.Fscan(reader, &n, &k, &a); err != nil {
		return
	}
	var m int
	fmt.Fscan(reader, &m)
	shots := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &shots[i])
	}

	s := &Set{}
	s.Insert(0)
	s.Insert(n + 1)
	total := calc(n, a)

	for i, x := range shots {
		prev := s.Prev(x)
		next := s.Next(x)
		total -= calc(next-prev-1, a)
		total += calc(x-prev-1, a)
		total += calc(next-x-1, a)
		s.Insert(x)
		if total < k {
			fmt.Fprintln(writer, i+1)
			return
		}
	}

	fmt.Fprintln(writer, -1)
}
