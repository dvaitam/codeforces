package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	left, right int
	sum         int
}

var (
	seg    []Node
	values []int
	g      [][]int
	a      []int
	root   []int
	parent [][]int
	depth  []int
	n      int
	LOG    int
)

func newNode() int {
	seg = append(seg, Node{})
	return len(seg) - 1
}

func update(prev, l, r, pos int) int {
	idx := newNode()
	seg[idx] = seg[prev]
	seg[idx].sum++
	if l != r {
		mid := (l + r) >> 1
		if pos <= mid {
			seg[idx].left = update(seg[prev].left, l, mid, pos)
		} else {
			seg[idx].right = update(seg[prev].right, mid+1, r, pos)
		}
	}
	return idx
}

func build(u, p int) {
	parent[0][u] = p
	depth[u] = depth[p] + 1
	root[u] = update(root[p], 1, 100000, a[u])
	for _, v := range g[u] {
		if v == p {
			continue
		}
		build(v, u)
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for k := LOG - 1; k >= 0; k-- {
		if diff>>uint(k)&1 == 1 {
			u = parent[k][u]
		}
	}
	if u == v {
		return u
	}
	for k := LOG - 1; k >= 0; k-- {
		if parent[k][u] != parent[k][v] {
			u = parent[k][u]
			v = parent[k][v]
		}
	}
	return parent[0][u]
}

func diff(ru1, rv1, rl1, ru2, rv2, rl2 int, valL1, valL2, l, r int) int {
	c1 := seg[ru1].sum + seg[rv1].sum - 2*seg[rl1].sum
	c2 := seg[ru2].sum + seg[rv2].sum - 2*seg[rl2].sum
	if valL1 >= l && valL1 <= r {
		c1++
	}
	if valL2 >= l && valL2 <= r {
		c2++
	}
	if c1 == c2 {
		return -1
	}
	if l == r {
		return l
	}
	mid := (l + r) >> 1
	res := diff(seg[ru1].left, seg[rv1].left, seg[rl1].left,
		seg[ru2].left, seg[rv2].left, seg[rl2].left,
		valL1, valL2, l, mid)
	if res != -1 {
		return res
	}
	return diff(seg[ru1].right, seg[rv1].right, seg[rl1].right,
		seg[ru2].right, seg[rv2].right, seg[rl2].right,
		valL1, valL2, mid+1, r)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	LOG = 0
	for (1 << LOG) <= n {
		LOG++
	}
	parent = make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		parent[i] = make([]int, n+1)
	}
	depth = make([]int, n+1)
	root = make([]int, n+1)
	seg = make([]Node, 1)
	build(1, 0)
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			parent[k][v] = parent[k-1][parent[k-1][v]]
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var u1, v1, u2, v2, k int
		fmt.Fscan(reader, &u1, &v1, &u2, &v2, &k)
		l1 := lca(u1, v1)
		l2 := lca(u2, v2)
		ans := diff(root[u1], root[v1], root[l1], root[u2], root[v2], root[l2], a[l1], a[l2], 1, 100000)
		if ans == -1 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, 1, ans)
		}
	}
}
