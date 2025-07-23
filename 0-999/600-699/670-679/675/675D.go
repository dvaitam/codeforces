package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Node represents a treap node
type Node struct {
	key   int
	prio  int32
	left  *Node
	right *Node
}

func rotateRight(t *Node) *Node {
	l := t.left
	t.left = l.right
	l.right = t
	return l
}

func rotateLeft(t *Node) *Node {
	r := t.right
	t.right = r.left
	r.left = t
	return r
}

func insert(t *Node, key int) *Node {
	if t == nil {
		return &Node{key: key, prio: rand.Int31()}
	}
	if key < t.key {
		t.left = insert(t.left, key)
		if t.left.prio < t.prio {
			t = rotateRight(t)
		}
	} else {
		t.right = insert(t.right, key)
		if t.right.prio < t.prio {
			t = rotateLeft(t)
		}
	}
	return t
}

func predecessor(t *Node, key int) *Node {
	var res *Node
	for t != nil {
		if key > t.key {
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
		if key < t.key {
			res = t
			t = t.left
		} else {
			t = t.right
		}
	}
	return res
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	idx := make(map[int]int, n)
	root := &Node{key: a[0], prio: rand.Int31()}
	idx[a[0]] = 0
	res := make([]int, n-1)

	for i := 1; i < n; i++ {
		x := a[i]
		p := predecessor(root, x)
		s := successor(root, x)
		var parent int
		if p == nil {
			parent = s.key
		} else if s == nil {
			parent = p.key
		} else if idx[p.key] > idx[s.key] {
			parent = p.key
		} else {
			parent = s.key
		}
		res[i-1] = parent
		idx[x] = i
		root = insert(root, x)
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
