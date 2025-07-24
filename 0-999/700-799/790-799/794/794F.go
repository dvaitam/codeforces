package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	cnt  [10]int64
	lazy [10]int
}

var (
	seg []Node
	arr []int
	n   int
)

func (node *Node) init() {
	for i := 0; i < 10; i++ {
		node.lazy[i] = i
	}
}

func build(v, l, r int) {
	seg[v].init()
	if l == r {
		x := arr[l]
		p := 1
		for x > 0 {
			d := x % 10
			seg[v].cnt[d] += int64(p)
			x /= 10
			p *= 10
		}
		return
	}
	m := (l + r) / 2
	build(v*2, l, m)
	build(v*2+1, m+1, r)
	for i := 0; i < 10; i++ {
		seg[v].cnt[i] = seg[v*2].cnt[i] + seg[v*2+1].cnt[i]
	}
}

func apply(v int, mapping [10]int) {
	var newCnt [10]int64
	for d := 0; d < 10; d++ {
		nd := mapping[d]
		newCnt[nd] += seg[v].cnt[d]
	}
	seg[v].cnt = newCnt
	var newLazy [10]int
	for d := 0; d < 10; d++ {
		newLazy[d] = mapping[seg[v].lazy[d]]
	}
	seg[v].lazy = newLazy
}

func isIdentity(lazy [10]int) bool {
	for i := 0; i < 10; i++ {
		if lazy[i] != i {
			return false
		}
	}
	return true
}

func push(v int) {
	if isIdentity(seg[v].lazy) {
		return
	}
	mapping := seg[v].lazy
	apply(v*2, mapping)
	apply(v*2+1, mapping)
	for i := 0; i < 10; i++ {
		seg[v].lazy[i] = i
	}
}

func update(v, l, r, ql, qr, x, y int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		var mapping [10]int
		for i := 0; i < 10; i++ {
			mapping[i] = i
		}
		mapping[x] = y
		apply(v, mapping)
		return
	}
	push(v)
	m := (l + r) / 2
	update(v*2, l, m, ql, qr, x, y)
	update(v*2+1, m+1, r, ql, qr, x, y)
	for i := 0; i < 10; i++ {
		seg[v].cnt[i] = seg[v*2].cnt[i] + seg[v*2+1].cnt[i]
	}
}

func query(v, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		var res int64
		for d := 0; d < 10; d++ {
			res += seg[v].cnt[d] * int64(d)
		}
		return res
	}
	push(v)
	m := (l + r) / 2
	return query(v*2, l, m, ql, qr) + query(v*2+1, m+1, r, ql, qr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &n, &q)
	arr = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	seg = make([]Node, 4*n)
	build(1, 0, n-1)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var l, r, x, y int
			fmt.Fscan(in, &l, &r, &x, &y)
			update(1, 0, n-1, l-1, r-1, x, y)
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			ans := query(1, 0, n-1, l-1, r-1)
			fmt.Fprintln(out, ans)
		}
	}
}
