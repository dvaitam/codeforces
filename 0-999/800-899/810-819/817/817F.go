package main

import (
	"bufio"
	"fmt"
	"os"
)

const LIMIT int64 = 1 << 60

type Node struct {
	left, right *Node
	val         int8 // 0 all zero, 1 all one, -1 mixed
	tag         int8 // 0 none, 1 set0, 2 set1, 3 flip
}

func applySet(n *Node, v int8) {
	n.val = v
	if v == 0 {
		n.tag = 1
	} else {
		n.tag = 2
	}
}

func applyFlip(n *Node) {
	if n.val != -1 {
		n.val ^= 1
	}
	switch n.tag {
	case 0:
		n.tag = 3
	case 1:
		n.tag = 2
	case 2:
		n.tag = 1
	case 3:
		n.tag = 0
	}
}

func push(n *Node) {
	if n.left == nil {
		n.left = &Node{val: n.val}
	}
	if n.right == nil {
		n.right = &Node{val: n.val}
	}
	if n.tag == 0 {
		return
	}
	if n.tag == 1 {
		applySet(n.left, 0)
		applySet(n.right, 0)
	} else if n.tag == 2 {
		applySet(n.left, 1)
		applySet(n.right, 1)
	} else if n.tag == 3 {
		applyFlip(n.left)
		applyFlip(n.right)
	}
	n.tag = 0
}

func update(n *Node, L, R, l, r int64, op int8) {
	if l > R || r < L {
		return
	}
	if l <= L && R <= r {
		if op == 1 {
			applySet(n, 1)
		} else if op == 2 {
			applySet(n, 0)
		} else {
			applyFlip(n)
		}
		return
	}
	push(n)
	mid := (L + R) >> 1
	update(n.left, L, mid, l, r, op)
	update(n.right, mid+1, R, l, r, op)
	if n.left.val == n.right.val {
		n.val = n.left.val
	} else {
		n.val = -1
	}
}

func mex(n *Node, L, R int64) int64 {
	if n == nil {
		return L
	}
	if n.val == 0 {
		return L
	}
	if n.val == 1 {
		return R + 1
	}
	if L == R {
		if n.val == 0 {
			return L
		}
		return L + 1
	}
	push(n)
	mid := (L + R) >> 1
	res := mex(n.left, L, mid)
	if res <= mid {
		return res
	}
	return mex(n.right, mid+1, R)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	root := &Node{val: 0}
	for i := 0; i < n; i++ {
		var t int
		var l, r int64
		fmt.Fscan(in, &t, &l, &r)
		update(root, 1, LIMIT, l, r, int8(t))
		ans := mex(root, 1, LIMIT)
		fmt.Fprintln(out, ans)
	}
}
