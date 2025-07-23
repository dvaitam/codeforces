package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Query struct {
	t   int
	a   int64
	b   int64
	idx int
	q   int64
}

type SegmentLine struct {
	a, b int64
	l, r int
}

type Line struct {
	m int64
	c int64
}
type LCNode struct {
	ln          Line
	left, right *LCNode
}

func eval(ln Line, x int64) int64 {
	return ln.m*x + ln.c
}

func insert(node *LCNode, l, r int, ln Line) *LCNode {
	if node == nil {
		return &LCNode{ln: ln}
	}
	newNode := &LCNode{ln: node.ln, left: node.left, right: node.right}
	node = newNode
	mid := (l + r) / 2
	midX := xs[mid]
	if eval(ln, midX) > eval(node.ln, midX) {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if eval(ln, xs[l]) > eval(node.ln, xs[l]) {
		node.left = insert(node.left, l, mid, ln)
	} else if eval(ln, xs[r]) > eval(node.ln, xs[r]) {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *LCNode, l, r int, x int64) int64 {
	if node == nil {
		return math.MinInt64
	}
	res := eval(node.ln, x)
	if l == r {
		return res
	}
	mid := (l + r) / 2
	if x <= xs[mid] {
		if v := query(node.left, l, mid, x); v > res {
			res = v
		}
	} else {
		if v := query(node.right, mid+1, r, x); v > res {
			res = v
		}
	}
	return res
}

var (
	seg     [][]int
	lines   []SegmentLine
	queries []Query
	answers []string
	xs      []int64
	m       int
	n       int
)

func addSeg(node, l, r, ql, qr, idx int) {
	if ql <= l && r <= qr {
		seg[node] = append(seg[node], idx)
		return
	}
	mid := (l + r) / 2
	if ql <= mid {
		addSeg(node*2, l, mid, ql, qr, idx)
	}
	if qr > mid {
		addSeg(node*2+1, mid+1, r, ql, qr, idx)
	}
}

func dfs(node, l, r int, root *LCNode) {
	for _, idx := range seg[node] {
		ln := lines[idx]
		root = insert(root, 0, m-1, Line{m: ln.a, c: ln.b})
	}
	if l == r {
		q := queries[l]
		if q.t == 3 {
			if root == nil {
				answers[l] = "EMPTY SET"
			} else {
				val := query(root, 0, m-1, q.q)
				answers[l] = fmt.Sprintf("%d", val)
			}
		}
		return
	}
	mid := (l + r) / 2
	dfs(node*2, l, mid, root)
	dfs(node*2+1, mid+1, r, root)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if n <= 0 {
		return
	}
	queries = make([]Query, n+1)
	lines = make([]SegmentLine, 0)
	addMap := make(map[int]int)
	xsList := make([]int64, 0)

	for i := 1; i <= n; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var a, b int64
			fmt.Fscan(reader, &a, &b)
			queries[i] = Query{t: t}
			lineIdx := len(lines)
			lines = append(lines, SegmentLine{a: a, b: b, l: i, r: n})
			addMap[i] = lineIdx
			queries[i].idx = lineIdx
		} else if t == 2 {
			var idx int
			fmt.Fscan(reader, &idx)
			queries[i] = Query{t: t, idx: idx}
			lineIdx := addMap[idx]
			lines[lineIdx].r = i - 1
		} else {
			var qv int64
			fmt.Fscan(reader, &qv)
			queries[i] = Query{t: t, q: qv}
			xsList = append(xsList, qv)
		}
	}

	// coordinate compression
	if len(xsList) == 0 {
		xsList = append(xsList, 0)
	}
	sort.Slice(xsList, func(i, j int) bool { return xsList[i] < xsList[j] })
	xs = make([]int64, 0, len(xsList))
	for i, v := range xsList {
		if i == 0 || v != xsList[i-1] {
			xs = append(xs, v)
		}
	}
	m = len(xs)

	seg = make([][]int, 4*n+5)
	for idx, ln := range lines {
		addSeg(1, 1, n, ln.l, ln.r, idx)
	}

	answers = make([]string, n+1)
	dfs(1, 1, n, nil)

	for i := 1; i <= n; i++ {
		if queries[i].t == 3 {
			fmt.Fprintln(writer, answers[i])
		}
	}
}
