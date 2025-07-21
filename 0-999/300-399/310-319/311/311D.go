package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 95542721
const cycle = 48

type Node struct {
	sums [cycle]int
	lazy int
}

var (
	n    int
	a    []int
	tree []Node
)

func apply(node, d int) {
	d %= cycle
	if d == 0 {
		return
	}
	// rotate sums by d
	var tmp [cycle]int
	for i := 0; i < cycle; i++ {
		tmp[i] = tree[node].sums[(i+d)%cycle]
	}
	tree[node].sums = tmp
	tree[node].lazy = (tree[node].lazy + d) % cycle
}

func push(node int) {
	d := tree[node].lazy
	if d != 0 {
		apply(node*2, d)
		apply(node*2+1, d)
		tree[node].lazy = 0
	}
}

func merge(node int) {
	left, right := node*2, node*2+1
	for i := 0; i < cycle; i++ {
		sum := tree[left].sums[i] + tree[right].sums[i]
		if sum >= mod {
			sum -= mod
		}
		tree[node].sums[i] = sum
	}
}

func build(node, l, r int) {
	if l == r {
		x := a[l] % mod
		for i := 0; i < cycle; i++ {
			tree[node].sums[i] = x
			// x = x^3 mod
			x = int((int64(x) * int64(x) % mod) * int64(x) % mod)
		}
	} else {
		m := (l + r) >> 1
		build(node*2, l, m)
		build(node*2+1, m+1, r)
		merge(node)
	}
}

func update(node, l, r, ql, qr int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		apply(node, 1)
		return
	}
	push(node)
	m := (l + r) >> 1
	update(node*2, l, m, ql, qr)
	update(node*2+1, m+1, r, ql, qr)
	merge(node)
}

func query(node, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return tree[node].sums[0]
	}
	push(node)
	m := (l + r) >> 1
	res := query(node*2, l, m, ql, qr) + query(node*2+1, m+1, r, ql, qr)
	if res >= mod {
		res %= mod
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &n)
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	tree = make([]Node, 4*(n+1))
	build(1, 1, n)
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		t, l, r := 0, 0, 0
		fmt.Fscan(reader, &t, &l, &r)
		if t == 1 {
			ans := query(1, 1, n, l, r)
			fmt.Fprintln(writer, ans)
		} else {
			update(1, 1, n, l, r)
		}
	}
}
