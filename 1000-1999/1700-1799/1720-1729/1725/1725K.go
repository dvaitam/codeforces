package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Node struct {
	key    int
	prio   int
	left   *Node
	right  *Node
	idxSet map[int]struct{}
}

type Treap struct{ root *Node }

func search(root *Node, key int) *Node {
	for root != nil {
		if key < root.key {
			root = root.left
		} else if key > root.key {
			root = root.right
		} else {
			return root
		}
	}
	return nil
}

func split(root *Node, key int) (*Node, *Node) {
	if root == nil {
		return nil, nil
	}
	if key <= root.key {
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
	if a.prio > b.prio {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func insert(root, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.prio > root.prio {
		l, r := split(root, node.key)
		node.left = l
		node.right = r
		return node
	}
	if node.key < root.key {
		root.left = insert(root.left, node)
	} else if node.key > root.key {
		root.right = insert(root.right, node)
	}
	return root
}

func erase(root *Node, key int) *Node {
	if root == nil {
		return nil
	}
	if key < root.key {
		root.left = erase(root.left, key)
	} else if key > root.key {
		root.right = erase(root.right, key)
	} else {
		root = merge(root.left, root.right)
	}
	return root
}

func (t *Treap) get(key int) *Node {
	return search(t.root, key)
}

func (t *Treap) ensure(key int) *Node {
	n := search(t.root, key)
	if n != nil {
		return n
	}
	n = &Node{key: key, prio: rand.Int(), idxSet: make(map[int]struct{})}
	t.root = insert(t.root, n)
	return n
}

func (t *Treap) deleteKey(key int) {
	t.root = erase(t.root, key)
}

func collect(n *Node, l, r int, res *[]*Node) {
	if n == nil {
		return
	}
	if l <= n.key {
		collect(n.left, l, r, res)
	}
	if l <= n.key && n.key <= r {
		*res = append(*res, n)
	}
	if n.key <= r {
		collect(n.right, l, r, res)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N int
	if _, err := fmt.Fscan(reader, &N); err != nil {
		return
	}
	heights := make([]int, N+1)
	tree := &Treap{}
	for i := 1; i <= N; i++ {
		fmt.Fscan(reader, &heights[i])
		n := tree.ensure(heights[i])
		n.idxSet[i] = struct{}{}
	}

	var Q int
	fmt.Fscan(reader, &Q)
	for ; Q > 0; Q-- {
		var typ int
		fmt.Fscan(reader, &typ)
		switch typ {
		case 1:
			var k, w int
			fmt.Fscan(reader, &k, &w)
			old := heights[k]
			if old == w {
				continue
			}
			if node := tree.get(old); node != nil {
				delete(node.idxSet, k)
				if len(node.idxSet) == 0 {
					tree.deleteKey(old)
				}
			}
			node := tree.ensure(w)
			node.idxSet[k] = struct{}{}
			heights[k] = w
		case 2:
			var k int
			fmt.Fscan(reader, &k)
			fmt.Fprintln(writer, heights[k])
		case 3:
			var l, r int
			fmt.Fscan(reader, &l, &r)
			nodes := []*Node{}
			collect(tree.root, l, r, &nodes)
			if len(nodes) == 0 {
				continue
			}
			sum := int64(l) + int64(r)
			leftVal := l - 1
			rightVal := r + 1
			for _, nd := range nodes {
				for idx := range nd.idxSet {
					var newVal int
					if int64(nd.key)*2 < sum {
						newVal = leftVal
					} else {
						newVal = rightVal
					}
					heights[idx] = newVal
					node := tree.ensure(newVal)
					node.idxSet[idx] = struct{}{}
				}
				nd.idxSet = nil
				tree.deleteKey(nd.key)
			}
		}
	}
}
