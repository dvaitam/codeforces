package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) bool {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return false
	}
	d.parent[a] = b
	return true
}

const maxN = 10

type Node struct {
	cnt   int
	left  [maxN]int
	right [maxN]int
}

var (
	n, m, q int
	grid    [][]int
	tree    []Node
)

func buildLeaf(col int) Node {
	dsu := NewDSU(n)
	for i := 1; i < n; i++ {
		if grid[i][col] == grid[i-1][col] {
			dsu.Union(i, i-1)
		}
	}
	mp := make(map[int]int)
	nextID := 0
	comp := make([]int, n)
	for i := 0; i < n; i++ {
		root := dsu.Find(i)
		id, ok := mp[root]
		if !ok {
			id = nextID
			mp[root] = id
			nextID++
		}
		comp[i] = id
	}
	var node Node
	node.cnt = nextID
	for i := 0; i < n; i++ {
		node.left[i] = comp[i]
		node.right[i] = comp[i]
	}
	return node
}

func merge(a Node, b Node, mid int) Node {
	dsu := NewDSU(a.cnt + b.cnt)
	merged := 0
	for i := 0; i < n; i++ {
		if grid[i][mid] == grid[i][mid+1] {
			if dsu.Union(a.right[i], b.left[i]+a.cnt) {
				merged++
			}
		}
	}
	res := Node{}
	res.cnt = a.cnt + b.cnt - merged
	mp := make(map[int]int)
	nextID := 0
	for i := 0; i < n; i++ {
		root := dsu.Find(a.left[i])
		id, ok := mp[root]
		if !ok {
			id = nextID
			mp[root] = id
			nextID++
		}
		res.left[i] = id
	}
	for i := 0; i < n; i++ {
		root := dsu.Find(b.right[i] + a.cnt)
		id, ok := mp[root]
		if !ok {
			id = nextID
			mp[root] = id
			nextID++
		}
		res.right[i] = id
	}
	return res
}

func build(idx, l, r int) {
	if l == r {
		tree[idx] = buildLeaf(l)
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1], mid)
}

func query(idx, l, r, ql, qr int) Node {
	if ql == l && qr == r {
		return tree[idx]
	}
	mid := (l + r) / 2
	if qr <= mid {
		return query(idx*2, l, mid, ql, qr)
	}
	if ql > mid {
		return query(idx*2+1, mid+1, r, ql, qr)
	}
	leftNode := query(idx*2, l, mid, ql, mid)
	rightNode := query(idx*2+1, mid+1, r, mid+1, qr)
	return merge(leftNode, rightNode, mid)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &q)
	grid = make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}
	tree = make([]Node, 4*m)
	build(1, 0, m-1)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		r--
		ans := query(1, 0, m-1, l, r)
		fmt.Fprintln(out, ans.cnt)
	}
}
