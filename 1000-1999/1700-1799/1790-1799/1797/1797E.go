package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXA = 5000000

var phi [MAXA + 1]int
var depth [MAXA + 1]int

func initPhi() {
	for i := 0; i <= MAXA; i++ {
		phi[i] = i
	}
	for i := 2; i <= MAXA; i++ {
		if phi[i] == i {
			for j := i; j <= MAXA; j += i {
				phi[j] -= phi[j] / i
			}
		}
	}
	depth[1] = 0
	for i := 2; i <= MAXA; i++ {
		depth[i] = depth[phi[i]] + 1
	}
}

func lca(x, y int) int {
	for x != y {
		if depth[x] > depth[y] {
			x = phi[x]
		} else if depth[y] > depth[x] {
			y = phi[y]
		} else {
			if x == y {
				break
			}
			x = phi[x]
			y = phi[y]
		}
	}
	return x
}

type Node struct {
	lca int
	sum int
	len int
}

var tree []Node
var a []int
var parent []int
var n int

func merge(a, b Node) Node {
	if a.len == 0 {
		return b
	}
	if b.len == 0 {
		return a
	}
	l := lca(a.lca, b.lca)
	s := a.sum + b.sum + a.len*(depth[a.lca]-depth[l]) + b.len*(depth[b.lca]-depth[l])
	return Node{lca: l, sum: s, len: a.len + b.len}
}

func build(idx, l, r int) {
	if l == r {
		tree[idx] = Node{lca: a[l], sum: 0, len: 1}
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func update(idx, l, r, pos, val int) {
	if l == r {
		tree[idx] = Node{lca: val, sum: 0, len: 1}
		return
	}
	mid := (l + r) / 2
	if pos <= mid {
		update(idx*2, l, mid, pos, val)
	} else {
		update(idx*2+1, mid+1, r, pos, val)
	}
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func query(idx, l, r, L, R int) Node {
	if L <= l && r <= R {
		return tree[idx]
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(idx*2, l, mid, L, R)
	}
	if L > mid {
		return query(idx*2+1, mid+1, r, L, R)
	}
	left := query(idx*2, l, mid, L, R)
	right := query(idx*2+1, mid+1, r, L, R)
	return merge(left, right)
}

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func rangeUpdate(l, r int) {
	for i := find(l); i <= r; i = find(i + 1) {
		a[i] = phi[a[i]]
		update(1, 1, n, i, a[i])
		if a[i] == 1 {
			parent[i] = find(i + 1)
		}
	}
}

func main() {
	initPhi()
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	fmt.Fscan(reader, &n, &m)
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	tree = make([]Node, 4*n+5)
	build(1, 1, n)
	parent = make([]int, n+2)
	for i := 1; i <= n; i++ {
		if a[i] == 1 {
			parent[i] = i + 1
		} else {
			parent[i] = i
		}
	}
	parent[n+1] = n + 1

	for ; m > 0; m-- {
		var t, l, r int
		fmt.Fscan(reader, &t, &l, &r)
		if t == 1 {
			rangeUpdate(l, r)
		} else {
			res := query(1, 1, n, l, r)
			fmt.Fprintln(writer, res.sum)
		}
	}
}
