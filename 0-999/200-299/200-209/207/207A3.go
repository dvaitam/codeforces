package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type item struct {
	val int64
	id  int
}

type node struct {
	key         item
	priority    uint32
	left, right *node
}

func less(a, b item) bool {
	if a.val != b.val {
		return a.val < b.val
	}
	return a.id < b.id
}

func rotateLeft(root *node) *node {
	r := root.right
	root.right = r.left
	r.left = root
	return r
}

func rotateRight(root *node) *node {
	l := root.left
	root.left = l.right
	l.right = root
	return l
}

var rng = rand.New(rand.NewSource(1))

func insert(root *node, it item) *node {
	if root == nil {
		return &node{key: it, priority: rng.Uint32()}
	}
	if less(it, root.key) {
		root.left = insert(root.left, it)
		if root.left.priority < root.priority {
			root = rotateRight(root)
		}
	} else {
		root.right = insert(root.right, it)
		if root.right.priority < root.priority {
			root = rotateLeft(root)
		}
	}
	return root
}

func merge(left, right *node) *node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority < right.priority {
		left.right = merge(left.right, right)
		return left
	}
	right.left = merge(left, right.left)
	return right
}

func erase(root *node, it item) *node {
	if root == nil {
		return nil
	}
	if root.key == it {
		return merge(root.left, root.right)
	}
	if less(it, root.key) {
		root.left = erase(root.left, it)
	} else {
		root.right = erase(root.right, it)
	}
	return root
}

func findMin(root *node) *node {
	if root == nil {
		return nil
	}
	for root.left != nil {
		root = root.left
	}
	return root
}

func lowerBound(root *node, target item) *node {
	var res *node
	for root != nil {
		if !less(root.key, target) {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func extractMin(root *node) (*node, item, bool) {
	if root == nil {
		return nil, item{}, false
	}
	if root.left == nil {
		return root.right, root.key, true
	}
	var extracted item
	root.left, extracted, _ = extractMin(root.left)
	return root, extracted, true
}

func extractLowerBound(root *node, target item) (*node, item, bool) {
	n := lowerBound(root, target)
	if n == nil {
		return root, item{}, false
	}
	root = erase(root, n.key)
	return root, n.key, true
}

type scientist struct {
	k         int
	x, y, m   int64
	lastValue int64
	produced  int
}

type outputEntry struct {
	val int64
	id  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	scis := make([]scientist, n)
	total := 0
	var root *node

	for i := 0; i < n; i++ {
		var (
			k  int
			a1 int64
			x  int64
			y  int64
			m  int64
		)
		fmt.Fscan(in, &k, &a1, &x, &y, &m)
		scis[i] = scientist{
			k:         k,
			x:         x,
			y:         y,
			m:         m,
			lastValue: a1,
			produced:  0,
		}
		total += k
		if k > 0 {
			scis[i].produced = 1
			root = insert(root, item{val: a1, id: i + 1})
		}
	}

	const limit = 200000
	var order []outputEntry
	if total <= limit {
		order = make([]outputEntry, 0, total)
	}

	var lastVal int64
	hasLast := false
	var bad int64
	processed := 0

	for processed < total {
		var current item
		var ok bool
		if !hasLast {
			root, current, ok = extractMin(root)
			if !ok {
				break
			}
			hasLast = true
		} else {
			root, current, ok = extractLowerBound(root, item{val: lastVal, id: 0})
			if !ok {
				root, current, ok = extractMin(root)
				if !ok {
					break
				}
				bad++
			}
		}

		lastVal = current.val
		processed++

		if total <= limit {
			order = append(order, outputEntry{val: current.val, id: current.id})
		}

		s := &scis[current.id-1]
		if s.produced < s.k {
			s.lastValue = (s.lastValue*s.x + s.y) % s.m
			s.produced++
			root = insert(root, item{val: s.lastValue, id: current.id})
		}
	}

	fmt.Fprintln(out, bad)
	if total <= limit {
		for _, e := range order {
			fmt.Fprintf(out, "%d %d\n", e.val, e.id)
		}
	}
}
