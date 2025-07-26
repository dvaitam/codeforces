package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type Pair struct {
	idx int
	typ int
}

type Node struct {
	key   Pair
	pr    int
	left  *Node
	right *Node
}

func less(a, b Pair) bool {
	if a.idx != b.idx {
		return a.idx < b.idx
	}
	return a.typ < b.typ
}

func split(root *Node, key Pair) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	if less(key, root.key) {
		l, r := split(root.left, key)
		root.left = r
		return l, root
	}
	l, r := split(root.right, key)
	root.right = l
	return root, r
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr < b.pr {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func insert(root *Node, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.pr < root.pr {
		l, r := split(root, node.key)
		node.left = l
		node.right = r
		return node
	}
	if less(node.key, root.key) {
		root.left = insert(root.left, node)
	} else {
		root.right = insert(root.right, node)
	}
	return root
}

func erase(root *Node, key Pair) *Node {
	if root == nil {
		return nil
	}
	if root.key == key {
		return merge(root.left, root.right)
	}
	if less(key, root.key) {
		root.left = erase(root.left, key)
	} else {
		root.right = erase(root.right, key)
	}
	return root
}

func lowerBound(root *Node, key Pair) *Node {
	var res *Node
	for root != nil {
		if less(root.key, key) {
			root = root.right
		} else {
			res = root
			root = root.left
		}
	}
	return res
}

func predecessor(root *Node, key Pair) *Node {
	var res *Node
	for root != nil {
		if less(root.key, key) {
			res = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *Node, key Pair) *Node {
	var res *Node
	for root != nil {
		if less(key, root.key) {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func maximum(root *Node) *Node {
	if root == nil {
		return nil
	}
	for root.right != nil {
		root = root.right
	}
	return root
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		a[i]--
	}

	stmn := []struct{ val, idx int }{{-1, -1}}
	stmx := []struct{ val, idx int }{{n, -1}}
	var root *Node
	root = insert(root, &Node{key: Pair{-1, 0}, pr: rand.Int()})
	root = insert(root, &Node{key: Pair{-1, 1}, pr: rand.Int()})
	length := 0
	var ans int64

	for i, x := range a {
		for stmn[len(stmn)-1].val > x {
			idx := stmn[len(stmn)-1].idx
			stmn = stmn[:len(stmn)-1]
			it := lowerBound(root, Pair{idx, 0})
			prv := predecessor(root, Pair{idx, 0})
			nxt := successor(root, Pair{idx, 0})
			if it != nil && prv != nil {
				length -= it.key.idx - prv.key.idx
			}
			if nxt != nil && prv != nil && nxt.key.typ == 0 {
				length += nxt.key.idx - prv.key.idx
			}
			root = erase(root, Pair{idx, 0})
		}
		maxNode := maximum(root)
		length += i - maxNode.key.idx
		root = insert(root, &Node{key: Pair{i, 0}, pr: rand.Int()})
		stmn = append(stmn, struct{ val, idx int }{x, i})

		for stmx[len(stmx)-1].val < x {
			idx := stmx[len(stmx)-1].idx
			stmx = stmx[:len(stmx)-1]
			it := lowerBound(root, Pair{idx, 1})
			prv := predecessor(root, Pair{idx, 1})
			nxt := successor(root, Pair{idx, 1})
			if nxt != nil && nxt.key.typ == 0 && prv != nil {
				length += it.key.idx - prv.key.idx
			}
			root = erase(root, Pair{idx, 1})
		}
		root = insert(root, &Node{key: Pair{i, 1}, pr: rand.Int()})
		stmx = append(stmx, struct{ val, idx int }{x, i})

		ans += int64(length)
	}
	fmt.Fprintln(out, ans-int64(n))
}
