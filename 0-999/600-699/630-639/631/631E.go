package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Line represents a linear function y = m*x + b.
type Line struct {
	m, b int64
}

func (ln Line) value(x int64) int64 { return ln.m*x + ln.b }

type Node struct {
	ln          Line
	left, right *Node
}

// insert adds a line into the Li Chao tree on interval [l, r].
func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	leftBetter := ln.value(l) > node.ln.value(l)
	midBetter := ln.value(mid) > node.ln.value(mid)
	if midBetter {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if leftBetter != midBetter {
		node.left = insert(node.left, l, mid, ln)
	} else {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

// query returns maximum value of any line at x on interval [l, r].
func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return math.MinInt64
	}
	res := node.ln.value(x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		val := query(node.left, l, mid, x)
		if val > res {
			res = val
		}
	} else {
		val := query(node.right, mid+1, r, x)
		if val > res {
			res = val
		}
	}
	return res
}

type LiChao struct {
	root        *Node
	left, right int64
}

func NewLiChao(l, r int64) *LiChao {
	return &LiChao{left: l, right: r}
}

func (lc *LiChao) Insert(ln Line)      { lc.root = insert(lc.root, lc.left, lc.right, ln) }
func (lc *LiChao) Query(x int64) int64 { return query(lc.root, lc.left, lc.right, x) }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pre := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pre[i] = pre[i-1] + a[i]
	}
	base := int64(0)
	for i := 1; i <= n; i++ {
		base += a[i] * int64(i)
	}
	ans := base

	tree := NewLiChao(1, int64(n))
	tree.Insert(Line{m: a[1], b: pre[1] - a[1]*1})
	for j := 2; j <= n; j++ {
		val := tree.Query(int64(j)) - pre[j]
		if base+val > ans {
			ans = base + val
		}
		tree.Insert(Line{m: a[j], b: pre[j] - a[j]*int64(j)})
	}

	tree2 := NewLiChao(1, int64(n))
	tree2.Insert(Line{m: a[n], b: pre[n-1] - a[n]*int64(n)})
	for j := n - 1; j >= 1; j-- {
		val := tree2.Query(int64(j)) - pre[j-1]
		if base+val > ans {
			ans = base + val
		}
		tree2.Insert(Line{m: a[j], b: pre[j-1] - a[j]*int64(j)})
	}

	fmt.Println(ans)
}
