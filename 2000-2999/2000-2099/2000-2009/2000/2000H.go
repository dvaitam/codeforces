package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const infVal = 6000005

var rng = rand.New(rand.NewSource(1))

type elemNode struct {
	key         int
	priority    uint32
	left, right *elemNode
}

type gapNode struct {
	start       int
	length      int
	maxLen      int
	priority    uint32
	left, right *gapNode
}

func elemInsert(root *elemNode, key int) *elemNode {
	if root == nil {
		return &elemNode{key: key, priority: rng.Uint32()}
	}
	if key < root.key {
		root.left = elemInsert(root.left, key)
		if root.left.priority < root.priority {
			root = elemRotateRight(root)
		}
	} else if key > root.key {
		root.right = elemInsert(root.right, key)
		if root.right.priority < root.priority {
			root = elemRotateLeft(root)
		}
	}
	return root
}

func elemRotateLeft(root *elemNode) *elemNode {
	r := root.right
	root.right = r.left
	r.left = root
	return r
}

func elemRotateRight(root *elemNode) *elemNode {
	l := root.left
	root.left = l.right
	l.right = root
	return l
}

func elemDelete(root *elemNode, key int) *elemNode {
	if root == nil {
		return nil
	}
	if key < root.key {
		root.left = elemDelete(root.left, key)
	} else if key > root.key {
		root.right = elemDelete(root.right, key)
	} else {
		if root.left == nil {
			return root.right
		}
		if root.right == nil {
			return root.left
		}
		if root.left.priority < root.right.priority {
			root = elemRotateRight(root)
			root.right = elemDelete(root.right, key)
		} else {
			root = elemRotateLeft(root)
			root.left = elemDelete(root.left, key)
		}
	}
	return root
}

func elemPrev(root *elemNode, key int) int {
	res := 0
	for root != nil {
		if key <= root.key {
			root = root.left
		} else {
			res = root.key
			root = root.right
		}
	}
	return res
}

func elemNext(root *elemNode, key int) int {
	res := infVal
	for root != nil {
		if key >= root.key {
			root = root.right
		} else {
			res = root.key
			root = root.left
		}
	}
	return res
}

func gapInsert(root *gapNode, start, length int) *gapNode {
	if length <= 0 {
		return root
	}
	if root == nil {
		return &gapNode{
			start:    start,
			length:   length,
			maxLen:   length,
			priority: rng.Uint32(),
		}
	}
	if start < root.start {
		root.left = gapInsert(root.left, start, length)
		if root.left.priority < root.priority {
			root = gapRotateRight(root)
		}
	} else if start > root.start {
		root.right = gapInsert(root.right, start, length)
		if root.right.priority < root.priority {
			root = gapRotateLeft(root)
		}
	} else {
		root.length = length
	}
	updateGap(root)
	return root
}

func gapDelete(root *gapNode, start int) *gapNode {
	if root == nil {
		return nil
	}
	if start < root.start {
		root.left = gapDelete(root.left, start)
	} else if start > root.start {
		root.right = gapDelete(root.right, start)
	} else {
		if root.left == nil {
			return root.right
		}
		if root.right == nil {
			return root.left
		}
		if root.left.priority < root.right.priority {
			root = gapRotateRight(root)
			root.right = gapDelete(root.right, start)
		} else {
			root = gapRotateLeft(root)
			root.left = gapDelete(root.left, start)
		}
	}
	updateGap(root)
	return root
}

func gapRotateLeft(root *gapNode) *gapNode {
	r := root.right
	root.right = r.left
	r.left = root
	updateGap(root)
	updateGap(r)
	return r
}

func gapRotateRight(root *gapNode) *gapNode {
	l := root.left
	root.left = l.right
	l.right = root
	updateGap(root)
	updateGap(l)
	return l
}

func updateGap(node *gapNode) {
	if node == nil {
		return
	}
	node.maxLen = node.length
	if node.left != nil && node.left.maxLen > node.maxLen {
		node.maxLen = node.left.maxLen
	}
	if node.right != nil && node.right.maxLen > node.maxLen {
		node.maxLen = node.right.maxLen
	}
}

func gapFind(root *gapNode, need int) int {
	cur := root
	for cur != nil {
		if cur.left != nil && cur.left.maxLen >= need {
			cur = cur.left
			continue
		}
		if cur.length >= need {
			return cur.start
		}
		cur = cur.right
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		elemRoot := (*elemNode)(nil)
		gapRoot := (*gapNode)(nil)
		elemRoot = elemInsert(elemRoot, 0)
		elemRoot = elemInsert(elemRoot, infVal)
		gapRoot = gapInsert(gapRoot, 1, infVal-1)

		initial := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &initial[i])
		}
		for _, val := range initial {
			elemRoot, gapRoot = insertValue(elemRoot, gapRoot, val)
		}

		var m int
		fmt.Fscan(in, &m)
		for i := 0; i < m; i++ {
			var op string
			fmt.Fscan(in, &op)
			if op == "+" {
				var x int
				fmt.Fscan(in, &x)
				elemRoot, gapRoot = insertValue(elemRoot, gapRoot, x)
			} else if op == "-" {
				var x int
				fmt.Fscan(in, &x)
				elemRoot, gapRoot = removeValue(elemRoot, gapRoot, x)
			} else {
				var k int
				fmt.Fscan(in, &k)
				ans := gapFind(gapRoot, k)
				if ans == -1 {
					ans = infVal
				}
				fmt.Fprintln(out, ans)
			}
		}
	}
}

func insertValue(elemRoot *elemNode, gapRoot *gapNode, x int) (*elemNode, *gapNode) {
	pre := elemPrev(elemRoot, x)
	suf := elemNext(elemRoot, x)
	gapLen := suf - pre - 1
	if gapLen > 0 {
		gapRoot = gapDelete(gapRoot, pre+1)
	}
	leftLen := x - pre - 1
	if leftLen > 0 {
		gapRoot = gapInsert(gapRoot, pre+1, leftLen)
	}
	rightLen := suf - x - 1
	if rightLen > 0 {
		gapRoot = gapInsert(gapRoot, x+1, rightLen)
	}
	elemRoot = elemInsert(elemRoot, x)
	return elemRoot, gapRoot
}

func removeValue(elemRoot *elemNode, gapRoot *gapNode, x int) (*elemNode, *gapNode) {
	pre := elemPrev(elemRoot, x)
	suf := elemNext(elemRoot, x)
	leftLen := x - pre - 1
	if leftLen > 0 {
		gapRoot = gapDelete(gapRoot, pre+1)
	}
	rightLen := suf - x - 1
	if rightLen > 0 {
		gapRoot = gapDelete(gapRoot, x+1)
	}
	mergedLen := suf - pre - 1
	if mergedLen > 0 {
		gapRoot = gapInsert(gapRoot, pre+1, mergedLen)
	}
	elemRoot = elemDelete(elemRoot, x)
	return elemRoot, gapRoot
}
