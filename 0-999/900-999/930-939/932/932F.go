package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Line struct {
	m int64
	c int64
}

type Node struct {
	ln    Line
	left  *Node
	right *Node
}

func eval(ln Line, x int64) int64 {
	return ln.m*x + ln.c
}

func insert(node *Node, l, r int64, ln Line) *Node {
	if node == nil {
		return &Node{ln: ln}
	}
	mid := (l + r) >> 1
	leftBetter := eval(ln, mid) < eval(node.ln, mid)
	if leftBetter {
		node.ln, ln = ln, node.ln
	}
	if l == r {
		return node
	}
	if eval(ln, l) < eval(node.ln, l) {
		node.left = insert(node.left, l, mid, ln)
	} else if eval(ln, r) < eval(node.ln, r) {
		node.right = insert(node.right, mid+1, r, ln)
	}
	return node
}

func query(node *Node, l, r, x int64) int64 {
	if node == nil {
		return math.MaxInt64
	}
	res := eval(node.ln, x)
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	if x <= mid {
		v := query(node.left, l, mid, x)
		if v < res {
			res = v
		}
	} else {
		v := query(node.right, mid+1, r, x)
		if v < res {
			res = v
		}
	}
	return res
}

type LiChao struct {
	root *Node
	size int
	l    int64
	r    int64
}

func NewLiChao(l, r int64) *LiChao {
	return &LiChao{l: l, r: r}
}

func (lc *LiChao) Insert(ln Line) {
	lc.root = insert(lc.root, lc.l, lc.r, ln)
	lc.size++
}

func (lc *LiChao) Query(x int64) int64 {
	return query(lc.root, lc.l, lc.r, x)
}

func (lc *LiChao) gather(n *Node) {
	if n == nil {
		return
	}
	lc.Insert(n.ln)
	lc.gather(n.left)
	lc.gather(n.right)
}

func (lc *LiChao) Merge(other *LiChao) {
	if other == nil || other.size == 0 {
		return
	}
	if lc.size < other.size {
		lc.root, other.root = other.root, lc.root
		lc.size, other.size = other.size, lc.size
	}
	lc.gather(other.root)
	lc.size += other.size
}

var (
	n    int
	a    []int64
	b    []int64
	adj  [][]int
	dp   []int64
	minA int64
	maxA int64
)

func dfs(u, p int) *LiChao {
	tree := NewLiChao(minA, maxA)
	isLeaf := true
	for _, v := range adj[u] {
		if v == p {
			continue
		}
		isLeaf = false
		child := dfs(v, u)
		tree.Merge(child)
	}
	if isLeaf && u != 1 {
		dp[u] = 0
	} else {
		dp[u] = tree.Query(a[u])
	}
	tree.Insert(Line{m: b[u], c: dp[u]})
	return tree
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a = make([]int64, n+1)
	b = make([]int64, n+1)
	adj = make([][]int, n+1)
	dp = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		if i == 1 {
			minA = a[i]
			maxA = a[i]
		} else {
			if a[i] < minA {
				minA = a[i]
			}
			if a[i] > maxA {
				maxA = a[i]
			}
		}
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dfs(1, 0)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, dp[i])
	}
	fmt.Fprint(writer, "\n")
}
