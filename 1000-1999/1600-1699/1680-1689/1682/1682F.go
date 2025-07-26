package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 1000000007

type Node struct {
	left, right int
	w           int64
	bw          int64
}

var tree []Node

func newNode() int {
	tree = append(tree, Node{})
	return len(tree) - 1
}

func update(prev, l, r, idx int, addW, addBW int64) int {
	cur := newNode()
	tree[cur] = tree[prev]
	tree[cur].w += addW
	tree[cur].bw = (tree[cur].bw + addBW) % MOD
	if l != r {
		mid := (l + r) >> 1
		if idx <= mid {
			tree[cur].left = update(tree[prev].left, l, mid, idx, addW, addBW)
		} else {
			tree[cur].right = update(tree[prev].right, mid+1, r, idx, addW, addBW)
		}
	}
	return cur
}

func query(node, l, r, ql, qr int) (int64, int64) {
	if node == 0 || ql > r || qr < l {
		return 0, 0
	}
	if ql <= l && r <= qr {
		return tree[node].w, tree[node].bw
	}
	mid := (l + r) >> 1
	w1, bw1 := query(tree[node].left, l, mid, ql, qr)
	w2, bw2 := query(tree[node].right, mid+1, r, ql, qr)
	return w1 + w2, (bw1 + bw2) % MOD
}

func modMul(a, b int64) int64 {
	a %= MOD
	if a < 0 {
		a += MOD
	}
	b %= MOD
	if b < 0 {
		b += MOD
	}
	return (a * b) % MOD
}

func calc(t int, x int64, roots []int, prefW, prefBW []int64, coord map[int64]int, vals []int64) int64 {
	if t <= 0 {
		return 0
	}
	pos := coord[x]
	wge, bwge := query(roots[t], 1, len(vals), pos, len(vals))
	wtot := prefW[t]
	bwTot := prefBW[t]
	temp1 := (2*bwge%MOD - bwTot) % MOD
	if temp1 < 0 {
		temp1 += MOD
	}
	temp2 := (wtot - 2*wge) % MOD
	if temp2 < 0 {
		temp2 += MOD
	}
	xmod := x % MOD
	if xmod < 0 {
		xmod += MOD
	}
	return (temp1 + xmod*temp2%MOD) % MOD
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	prefB := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefB[i] = prefB[i-1] + b[i]
	}

	vals := append([]int64(nil), prefB...)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	uniq := vals[:1]
	for i := 1; i < len(vals); i++ {
		if vals[i] != vals[i-1] {
			uniq = append(uniq, vals[i])
		}
	}
	vals = uniq
	coord := make(map[int64]int, len(vals))
	for i, v := range vals {
		coord[v] = i + 1
	}

	tree = append(tree, Node{})
	roots := make([]int, n)
	prefW := make([]int64, n)
	prefBW := make([]int64, n)
	for i := 1; i <= n-1; i++ {
		pos := coord[prefB[i]]
		w := a[i+1] - a[i]
		bw := modMul(w, prefB[i])
		roots[i] = update(roots[i-1], 1, len(vals), pos, w, bw)
		prefW[i] = prefW[i-1] + w
		prefBW[i] = (prefBW[i-1] + bw) % MOD
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		x := prefB[r]
		ans := calc(r-1, x, roots, prefW, prefBW, coord, vals)
		ans -= calc(l-1, x, roots, prefW, prefBW, coord, vals)
		ans %= MOD
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(writer, ans)
	}
}
