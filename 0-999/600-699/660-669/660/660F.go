package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Line struct {
	m int64
	b int64
}

type Node struct {
	ln    Line
	left  *Node
	right *Node
}

func eval(ln Line, x int64) int64 {
	return ln.m*x + ln.b
}

func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	if eval(ln, mid) > eval(node.ln, mid) {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if eval(ln, l) > eval(node.ln, l) {
		node.left = insert(node.left, l, mid, ln)
	} else if eval(ln, r) > eval(node.ln, r) {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return math.MinInt64
	}
	res := eval(node.ln, x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		v := query(node.left, l, mid, x)
		if v > res {
			res = v
		}
	} else {
		v := query(node.right, mid+1, r, x)
		if v > res {
			res = v
		}
	}
	return res
}

type LiChao struct {
	root *Node
	l    int64
	r    int64
}

func NewLiChao(l, r int64) *LiChao {
	return &LiChao{l: l, r: r}
}

func (lc *LiChao) Insert(ln Line) {
	lc.root = insert(lc.root, lc.l, lc.r, ln)
}

func (lc *LiChao) Query(x int64) int64 {
	return query(lc.root, lc.l, lc.r, x)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	S := make([]int64, n+1)
	T := make([]int64, n+1)
	minS, maxS := int64(0), int64(0)
	for i := 1; i <= n; i++ {
		S[i] = S[i-1] + a[i]
		if S[i] < minS {
			minS = S[i]
		}
		if S[i] > maxS {
			maxS = S[i]
		}
		T[i] = T[i-1] + int64(i)*a[i]
	}

	lc := NewLiChao(minS, maxS)
	lc.Insert(Line{m: 0, b: 0})
	ans := int64(0)
	for r := 1; r <= n; r++ {
		val := lc.Query(S[r])
		cand := T[r] + val
		if cand > ans {
			ans = cand
		}
		line := Line{m: -int64(r), b: int64(r)*S[r] - T[r]}
		lc.Insert(line)
	}
	fmt.Fprintln(writer, ans)
}
