package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var powers []int64

func init() {
	powers = make([]int64, 0)
	p := int64(1)
	for p > 0 && p < 1<<62 {
		powers = append(powers, p)
		if p > (1<<62)/42 {
			break
		}
		p *= 42
	}
}

func nextBadDiff(v int64) int64 {
	// return minimal power of 42 >= v minus v
	for _, p := range powers {
		if p >= v {
			return p - v
		}
	}
	return (1 << 62) - v
}

// interval node for ordered disjoint interval tree (treap)
type Node struct {
	l, r        int
	val         int64
	pri         int
	left, right *Node
}

func newNode(l, r int, val int64) *Node {
	return &Node{l: l, r: r, val: val, pri: rand.Int()}
}

func rotateLeft(x *Node) *Node {
	y := x.right
	x.right = y.left
	y.left = x
	return y
}

func rotateRight(x *Node) *Node {
	y := x.left
	x.left = y.right
	y.right = x
	return y
}

func insert(root *Node, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.l < root.l {
		root.left = insert(root.left, node)
		if root.left.pri < root.pri {
			root = rotateRight(root)
		}
	} else {
		root.right = insert(root.right, node)
		if root.right.pri < root.pri {
			root = rotateLeft(root)
		}
	}
	return root
}

func mergeTrees(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pri < b.pri {
		a.right = mergeTrees(a.right, b)
		return a
	}
	b.left = mergeTrees(a, b.left)
	return b
}

func erase(root *Node, key int) *Node {
	if root == nil {
		return nil
	}
	if key < root.l {
		root.left = erase(root.left, key)
		return root
	} else if key > root.l {
		root.right = erase(root.right, key)
		return root
	}
	// erase this node
	return mergeTrees(root.left, root.right)
}

func lowerBound(root *Node, key int) *Node {
	var res *Node
	for root != nil {
		if root.l >= key {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func predecessor(root *Node, key int) *Node {
	var res *Node
	for root != nil {
		if root.l < key {
			res = root
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func maxNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

func minNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n
}

// split tree into [< key] and [>= key]
func split(root *Node, key int) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	if key <= root.l {
		l, r := split(root.left, key)
		root.left = r
		return l, root
	}
	l, r := split(root.right, key)
	root.right = l
	return root, r
}

// splitByPos ensures there is an interval starting at pos
func splitByPos(root *Node, pos int) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	// ensure interval with start=pos exists
	node := predecessor(root, pos+1)
	if node != nil && node.r >= pos && node.l < pos {
		// split node
		leftPart := newNode(node.l, pos-1, node.val)
		rightPart := newNode(pos, node.r, node.val)
		root = erase(root, node.l)
		if leftPart.l <= leftPart.r {
			root = insert(root, leftPart)
		}
		root = insert(root, rightPart)
	}
	return split(root, pos)
}

// traverse and apply function to each node
func traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	traverse(n.left, f)
	f(n)
	traverse(n.right, f)
}

// add value to all nodes in subtree
func addAll(n *Node, delta int64) {
	traverse(n, func(x *Node) { x.val += delta })
}

func hasBad(n *Node) bool {
	bad := false
	traverse(n, func(x *Node) {
		if !bad {
			d := nextBadDiff(x.val)
			if d == 0 {
				bad = true
			}
		}
	})
	return bad
}

func mergeAdjacent(root *Node, node *Node) *Node {
	// merge with predecessor if same value
	pre := predecessor(root, node.l)
	if pre != nil && pre.r+1 == node.l && pre.val == node.val {
		root = erase(root, pre.l)
		node.l = pre.l
	}
	// merge with successor
	suc := lowerBound(root, node.r+1)
	if suc != nil && node.r+1 == suc.l && suc.val == node.val {
		root = erase(root, suc.l)
		node.r = suc.r
	}
	root = erase(root, node.l)
	root = insert(root, node)
	return root
}

func main() {
	rand.Seed(1)
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &q)
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	// build initial intervals
	var root *Node
	for i := 0; i < n; {
		j := i + 1
		for j < n && arr[j] == arr[i] {
			j++
		}
		root = insert(root, newNode(i+1, j, arr[i]))
		i = j
	}
	out := bufio.NewWriter(os.Stdout)
	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(in, &typ)
		if typ == 1 {
			var idx int
			fmt.Fscan(in, &idx)
			it := predecessor(root, idx+1)
			if it != nil && it.r >= idx {
				fmt.Fprintln(out, it.val)
			}
		} else if typ == 2 {
			var l, r int
			var x int64
			fmt.Fscan(in, &l, &r, &x)
			// split
			left, tmp := splitByPos(root, l)
			_, right := splitByPos(tmp, r+1)
			node := newNode(l, r, x)
			root = mergeTrees(left, node)
			root = mergeTrees(root, right)
			// merge with neighbors
			root = mergeAdjacent(root, node)
		} else if typ == 3 {
			var l, r int
			var x int64
			fmt.Fscan(in, &l, &r, &x)
			left, tmp := splitByPos(root, l)
			mid, right := splitByPos(tmp, r+1)
			if mid != nil {
				for {
					addAll(mid, x)
					if !hasBad(mid) {
						break
					}
				}
			}
			root = mergeTrees(left, mid)
			root = mergeTrees(root, right)
			// merge boundaries if possible
			if mid != nil {
				lt := minNode(mid)
				rt := maxNode(mid)
				if lt != nil {
					root = mergeAdjacent(root, lt)
				}
				if rt != nil && rt != lt {
					root = mergeAdjacent(root, rt)
				}
			}
		}
	}
	out.Flush()
}
