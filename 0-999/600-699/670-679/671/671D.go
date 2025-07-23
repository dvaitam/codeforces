package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	val         int64
	to          int
	add         int64
	left, right *Node
}

func apply(n *Node, d int64) {
	if n != nil {
		n.val += d
		n.add += d
	}
}

func push(n *Node) {
	if n != nil && n.add != 0 {
		apply(n.left, n.add)
		apply(n.right, n.add)
		n.add = 0
	}
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.val > b.val {
		a, b = b, a
	}
	push(a)
	a.right = merge(a.right, b)
	a.left, a.right = a.right, a.left
	return a
}

var (
	g       [][]int
	roots   []*Node
	visited []bool
	top     []int
	weight  []int64
	ans     int64
)

func dfs(u, p int) {
	for _, v := range g[u] {
		if v != p {
			dfs(v, u)
			roots[u] = merge(roots[u], roots[v])
		}
	}
	visited[u] = true
	if u == 1 {
		return
	}
	for roots[u] != nil {
		push(roots[u])
		if visited[roots[u].to] {
			roots[u] = merge(roots[u].left, roots[u].right)
		} else {
			break
		}
	}
	if roots[u] == nil {
		fmt.Println(-1)
		os.Exit(0)
	}
	ans += roots[u].val
	apply(roots[u], -roots[u].val)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	roots = make([]*Node, n+1)
	top = make([]int, m+1)
	weight = make([]int64, m+1)
	for i := 1; i <= m; i++ {
		var u, v int
		var c int64
		fmt.Fscan(reader, &u, &v, &c)
		top[i] = v
		weight[i] = c
		node := &Node{val: c, to: v}
		roots[u] = merge(roots[u], node)
	}
	visited = make([]bool, n+1)
	dfs(1, 0)
	fmt.Println(ans)
}
