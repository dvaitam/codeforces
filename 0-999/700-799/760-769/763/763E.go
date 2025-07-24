package main

import (
	"bufio"
	"fmt"
	"os"
)

const kMax = 5

var (
	n, K  int
	edges [][]bool
	tree  []Node
)

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n)}
	for i := 0; i < n; i++ {
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

func hasEdge(u, v int) bool {
	if u > v {
		u, v = v, u
	}
	diff := v - u
	if diff > K || diff <= 0 {
		return false
	}
	return edges[u][diff]
}

type Node struct {
	length int
	cnt    int
	lidx   [kMax]int
	lcomp  [kMax]int
	lsize  int
	ridx   [kMax]int
	rcomp  [kMax]int
	rsize  int
}

func makeLeaf(idx int) Node {
	var node Node
	node.length = 1
	node.cnt = 1
	node.lidx[0] = idx
	node.lcomp[0] = 0
	node.lsize = 1
	node.ridx[0] = idx
	node.rcomp[0] = 0
	node.rsize = 1
	return node
}

func merge(a, b Node) Node {
	if a.length == 0 {
		return b
	}
	if b.length == 0 {
		return a
	}
	var res Node
	res.length = a.length + b.length

	dsu := NewDSU(a.cnt + b.cnt)
	merges := 0
	for i := 0; i < a.rsize; i++ {
		for j := 0; j < b.lsize; j++ {
			if hasEdge(a.ridx[i], b.lidx[j]) {
				if dsu.Union(a.rcomp[i], b.lcomp[j]+a.cnt) {
					merges++
				}
			}
		}
	}
	res.cnt = a.cnt + b.cnt - merges

	leftNeed := res.length
	if leftNeed > K {
		leftNeed = K
	}
	pos := 0
	if a.length >= leftNeed {
		for i := 0; i < leftNeed; i++ {
			res.lidx[pos] = a.lidx[i]
			res.lcomp[pos] = a.lcomp[i]
			pos++
		}
	} else {
		for i := 0; i < a.lsize; i++ {
			res.lidx[pos] = a.lidx[i]
			res.lcomp[pos] = a.lcomp[i]
			pos++
		}
		need := leftNeed - a.length
		for j := 0; j < need; j++ {
			res.lidx[pos] = b.lidx[j]
			res.lcomp[pos] = b.lcomp[j] + a.cnt
			pos++
		}
	}
	res.lsize = leftNeed

	rightNeed := res.length
	if rightNeed > K {
		rightNeed = K
	}
	pos = 0
	if b.length >= rightNeed {
		start := b.rsize - rightNeed
		for i := start; i < b.rsize; i++ {
			res.ridx[pos] = b.ridx[i]
			res.rcomp[pos] = b.rcomp[i] + a.cnt
			pos++
		}
	} else {
		need := rightNeed - b.length
		startA := a.rsize - need
		for i := startA; i < a.rsize; i++ {
			res.ridx[pos] = a.ridx[i]
			res.rcomp[pos] = a.rcomp[i]
			pos++
		}
		for i := 0; i < b.rsize; i++ {
			res.ridx[pos] = b.ridx[i]
			res.rcomp[pos] = b.rcomp[i] + a.cnt
			pos++
		}
	}
	res.rsize = rightNeed

	rootToID := make(map[int]int)
	nextID := 0
	for i := 0; i < res.lsize; i++ {
		root := dsu.Find(res.lcomp[i])
		id, ok := rootToID[root]
		if !ok {
			id = nextID
			rootToID[root] = id
			nextID++
		}
		res.lcomp[i] = id
	}
	for i := 0; i < res.rsize; i++ {
		root := dsu.Find(res.rcomp[i])
		id, ok := rootToID[root]
		if !ok {
			id = nextID
			rootToID[root] = id
			nextID++
		}
		res.rcomp[i] = id
	}

	return res
}

func build(idx, l, r int) {
	if l == r {
		tree[idx] = makeLeaf(l)
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
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
	left := query(idx*2, l, mid, ql, mid)
	right := query(idx*2+1, mid+1, r, mid+1, qr)
	return merge(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &K)
	edges = make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		edges[i] = make([]bool, K+1)
	}
	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		if u > v {
			u, v = v, u
		}
		diff := v - u
		if diff <= K {
			edges[u][diff] = true
		}
	}
	tree = make([]Node, 4*n+5)
	build(1, 1, n)

	var q int
	fmt.Fscan(reader, &q)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		res := query(1, 1, n, l, r)
		fmt.Fprintln(writer, res.cnt)
	}
}
