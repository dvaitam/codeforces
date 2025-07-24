package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

// Treap node storing point (y,x)
type node struct {
	y, x        int64
	pr          int
	left, right *node
}

func split(root *node, y int64) (l, r *node) {
	if root == nil {
		return nil, nil
	}
	if y < root.y {
		l, root.left = split(root.left, y)
		return l, root
	}
	root.right, r = split(root.right, y)
	return root, r
}

func merge(a, b *node) *node {
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

func insert(root, n *node) *node {
	if root == nil {
		return n
	}
	if n.pr < root.pr {
		l, r := split(root, n.y)
		n.left = l
		n.right = r
		return n
	}
	if n.y < root.y {
		root.left = insert(root.left, n)
	} else {
		root.right = insert(root.right, n)
	}
	return root
}

func erase(root *node, y int64) *node {
	if root == nil {
		return nil
	}
	if y < root.y {
		root.left = erase(root.left, y)
		return root
	}
	if y > root.y {
		root.right = erase(root.right, y)
		return root
	}
	return merge(root.left, root.right)
}

func find(root *node, y int64) *node {
	for root != nil {
		if y < root.y {
			root = root.left
		} else if y > root.y {
			root = root.right
		} else {
			return root
		}
	}
	return nil
}

func predecessor(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y <= root.y {
			root = root.left
		} else {
			res = root
			root = root.right
		}
	}
	return res
}

func successor(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y >= root.y {
			root = root.right
		} else {
			res = root
			root = root.left
		}
	}
	return res
}

func lowerBound(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y <= root.y {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func minNode(root *node) *node {
	if root == nil {
		return nil
	}
	for root.left != nil {
		root = root.left
	}
	return root
}

func maxNode(root *node) *node {
	if root == nil {
		return nil
	}
	for root.right != nil {
		root = root.right
	}
	return root
}

func cross(a, b, c *node) int64 {
	return (b.y-a.y)*(c.x-b.x) - (b.x-a.x)*(c.y-b.y)
}

var root *node

func isBad(n *node) bool {
	if n == nil || n.y == 0 {
		return false
	}
	p := predecessor(root, n.y)
	r := successor(root, n.y)
	if p == nil || r == nil {
		return false
	}
	if p.y == 0 && r == nil {
		return false
	}
	return cross(p, n, r) >= 0
}

func addSpell(x, y int64) {
	if existing := find(root, y); existing != nil {
		if existing.x >= x {
			return
		}
		root = erase(root, y)
	}
	n := &node{y: y, x: x, pr: rand.Int()}
	root = insert(root, n)
	if isBad(n) {
		root = erase(root, n.y)
		return
	}
	for {
		p := predecessor(root, n.y)
		if p == nil || p.y == 0 {
			break
		}
		if isBad(p) {
			root = erase(root, p.y)
		} else {
			break
		}
	}
	for {
		r := successor(root, n.y)
		if r == nil {
			break
		}
		if isBad(r) {
			root = erase(root, r.y)
		} else {
			break
		}
	}
}

func hullValue(z float64) float64 {
	if root == nil {
		return 0
	}
	if z <= 0 {
		return 0
	}
	mx := maxNode(root)
	if z >= float64(mx.y) {
		return float64(mx.x)
	}
	r := lowerBound(root, int64(math.Ceil(z)))
	if r == nil {
		r = mx
	}
	if float64(r.y) == z {
		return float64(r.x)
	}
	l := predecessor(root, r.y)
	if l == nil {
		return (float64(r.x) / float64(r.y)) * z
	}
	return float64(l.x) + (float64(r.x-l.x))*(z-float64(l.y))/float64(r.y-l.y)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	var m int64
	if _, err := fmt.Fscan(reader, &q, &m); err != nil {
		return
	}
	// initialize hull with origin
	root = &node{y: 0, x: 0, pr: rand.Int()}

	lastOK := 0
	for i := 1; i <= q; i++ {
		var k int
		var a, b int64
		fmt.Fscan(reader, &k, &a, &b)
		a = (a+int64(lastOK))%1000000 + 1
		b = (b+int64(lastOK))%1000000 + 1
		if k == 1 {
			addSpell(a, b) // x=a, y=b
		} else {
			t := a
			h := b
			rate := hullValue(float64(m) / float64(t))
			if rate*float64(t) >= float64(h) {
				lastOK = i
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
