package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Treap node representing an interval [l,r]
type Node struct {
	l, r  int
	color int
	val   int64
	pri   int
	left  *Node
	right *Node
}

// merge two treaps
func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pri < b.pri {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

// split treap by key (interval start)
func split(t *Node, key int) (l, r *Node) {
	if t == nil {
		return nil, nil
	}
	if key <= t.l {
		lsub, rsub := split(t.left, key)
		t.left = rsub
		return lsub, t
	}
	lsub, rsub := split(t.right, key)
	t.right = lsub
	return t, rsub
}

func insert(t *Node, x *Node) *Node {
	if t == nil {
		return x
	}
	if x.pri < t.pri {
		l, r := split(t, x.l)
		x.left, x.right = l, r
		return x
	}
	if x.l < t.l {
		t.left = insert(t.left, x)
	} else {
		t.right = insert(t.right, x)
	}
	return t
}

func erase(t *Node, key int) *Node {
	if t == nil {
		return nil
	}
	if key == t.l {
		return merge(t.left, t.right)
	}
	if key < t.l {
		t.left = erase(t.left, key)
	} else {
		t.right = erase(t.right, key)
	}
	return t
}

func predecessor(t *Node, key int) *Node {
	var res *Node
	for t != nil {
		if t.l < key {
			res = t
			t = t.right
		} else {
			t = t.left
		}
	}
	return res
}

func successor(t *Node, key int) *Node {
	var res *Node
	for t != nil {
		if t.l > key {
			res = t
			t = t.left
		} else {
			t = t.right
		}
	}
	return res
}

func findSegment(t *Node, pos int) *Node {
	n := predecessor(t, pos+1)
	if n != nil && n.r >= pos {
		return n
	}
	return nil
}

// ensure there is a segment starting at pos
func splitPos(root **Node, pos int) *Node {
	if *root == nil {
		return nil
	}
	n := findSegment(*root, pos)
	if n == nil || n.l == pos {
		return n
	}
	// remove current interval and split
	*root = erase(*root, n.l)
	left := &Node{l: n.l, r: pos - 1, color: n.color, val: n.val, pri: rand.Int()}
	right := &Node{l: pos, r: n.r, color: n.color, val: n.val, pri: rand.Int()}
	*root = insert(*root, left)
	*root = insert(*root, right)
	return right
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

	add := make([]int64, n+2)
	root := &Node{l: 1, r: n, color: 1, val: 0, pri: rand.Int()}

	for ; q > 0; q-- {
		var op string
		if _, err := fmt.Fscan(in, &op); err != nil {
			return
		}
		switch op[0] {
		case 'C':
			var l, r, c int
			fmt.Fscan(in, &l, &r, &c)
			if l > r {
				l, r = r, l
			}
			splitPos(&root, r+1)
			splitPos(&root, l)
			for node := findSegment(root, l); node != nil && node.l <= r; node = successor(root, node.l) {
				actual := node.val + add[node.color]
				node.val = actual - add[c]
				node.color = c
			}
		case 'A':
			var c int
			var x int64
			fmt.Fscan(in, &c, &x)
			add[c] += x
		case 'Q':
			var idx int
			fmt.Fscan(in, &idx)
			node := findSegment(root, idx)
			if node != nil {
				fmt.Fprintln(out, node.val+add[node.color])
			} else {
				fmt.Fprintln(out, 0)
			}
		}
	}
}
