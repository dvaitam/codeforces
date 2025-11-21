package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type scientist struct {
	total     int
	processed int
	x, y, m   int64
}

type treapNode struct {
	key      int64
	priority uint32
	left     *treapNode
	right    *treapNode
	ids      []int
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func newNode(key int64, id int) *treapNode {
	return &treapNode{
		key:      key,
		priority: rng.Uint32(),
		ids:      []int{id},
	}
}

func rotateRight(root *treapNode) *treapNode {
	left := root.left
	root.left = left.right
	left.right = root
	return left
}

func rotateLeft(root *treapNode) *treapNode {
	right := root.right
	root.right = right.left
	right.left = root
	return right
}

func insert(root *treapNode, key int64, id int) *treapNode {
	if root == nil {
		return newNode(key, id)
	}
	if key == root.key {
		root.ids = append(root.ids, id)
		return root
	}
	if key < root.key {
		root.left = insert(root.left, key, id)
		if root.left.priority < root.priority {
			root = rotateRight(root)
		}
	} else {
		root.right = insert(root.right, key, id)
		if root.right.priority < root.priority {
			root = rotateLeft(root)
		}
	}
	return root
}

func deleteKey(root *treapNode, key int64) *treapNode {
	if root == nil {
		return nil
	}
	if key < root.key {
		root.left = deleteKey(root.left, key)
	} else if key > root.key {
		root.right = deleteKey(root.right, key)
	} else {
		if root.left == nil {
			return root.right
		}
		if root.right == nil {
			return root.left
		}
		if root.left.priority < root.right.priority {
			root = rotateRight(root)
			root.right = deleteKey(root.right, key)
		} else {
			root = rotateLeft(root)
			root.left = deleteKey(root.left, key)
		}
	}
	return root
}

func findMin(root *treapNode) *treapNode {
	if root == nil {
		return nil
	}
	for root.left != nil {
		root = root.left
	}
	return root
}

func lowerBound(root *treapNode, key int64) *treapNode {
	var res *treapNode
	cur := root
	for cur != nil {
		if cur.key >= key {
			res = cur
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	sci := make([]scientist, n)
	totalProblems := 0
	var root *treapNode

	for i := 0; i < n; i++ {
		var k int
		var a1, x, y, m int64
		fmt.Fscan(in, &k, &a1, &x, &y, &m)
		sci[i] = scientist{
			total: k,
			x:     x,
			y:     y,
			m:     m,
		}
		totalProblems += k
		if k > 0 {
			root = insert(root, a1, i)
		}
	}

	needPrint := totalProblems <= 200000
	var outVals []int64
	var outIDs []int
	if needPrint {
		outVals = make([]int64, 0, totalProblems)
		outIDs = make([]int, 0, totalProblems)
	}

	var badPairs int64
	var last int64
	hasLast := false
	processed := 0

	for processed < totalProblems {
		var node *treapNode
		if !hasLast {
			node = findMin(root)
		} else {
			node = lowerBound(root, last)
			if node == nil {
				node = findMin(root)
			}
		}
		if node == nil {
			break
		}

		idx := len(node.ids) - 1
		id := node.ids[idx]
		node.ids = node.ids[:idx]
		value := node.key
		if len(node.ids) == 0 {
			root = deleteKey(root, node.key)
		}

		if hasLast && last > value {
			badPairs++
		}
		last = value
		hasLast = true
		processed++

		if needPrint {
			outVals = append(outVals, value)
			outIDs = append(outIDs, id+1)
		}

		s := &sci[id]
		s.processed++
		if s.processed < s.total {
			nextVal := (value*s.x + s.y) % s.m
			root = insert(root, nextVal, id)
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, badPairs)
	if needPrint {
		for i := range outVals {
			fmt.Fprintln(writer, outVals[i], outIDs[i])
		}
	}
}
