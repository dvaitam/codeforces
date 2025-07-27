package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	first int
	last  int
	pref  int
	suff  int
	len   int
	ans   int64
}

func makeNode(val int) Node {
	return Node{first: val, last: val, pref: 1, suff: 1, len: 1, ans: 1}
}

func merge(a, b Node) Node {
	if a.len == 0 {
		return b
	}
	if b.len == 0 {
		return a
	}
	res := Node{}
	res.len = a.len + b.len
	res.first = a.first
	res.last = b.last
	res.pref = a.pref
	if a.pref == a.len && a.last <= b.first {
		res.pref = a.len + b.pref
	}
	res.suff = b.suff
	if b.suff == b.len && a.last <= b.first {
		res.suff = b.len + a.suff
	}
	res.ans = a.ans + b.ans
	if a.last <= b.first {
		res.ans += int64(a.suff) * int64(b.pref)
	}
	return res
}

var tree []Node
var arr []int
var n int

func build(idx, l, r int) {
	if l == r {
		tree[idx] = makeNode(arr[l])
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func update(idx, l, r, pos, val int) {
	if l == r {
		arr[pos] = val
		tree[idx] = makeNode(val)
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

func query(idx, l, r, ql, qr int) Node {
	if ql <= l && r <= qr {
		return tree[idx]
	}
	mid := (l + r) / 2
	if qr <= mid {
		return query(idx*2, l, mid, ql, qr)
	} else if ql > mid {
		return query(idx*2+1, mid+1, r, ql, qr)
	}
	left := query(idx*2, l, mid, ql, mid)
	right := query(idx*2+1, mid+1, r, mid+1, qr)
	return merge(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	arr = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	tree = make([]Node, 4*n)
	build(1, 0, n-1)

	for ; q > 0; q-- {
		var t, x, y int
		fmt.Fscan(reader, &t, &x, &y)
		if t == 1 {
			// update: positions are 1-based
			pos := x - 1
			update(1, 0, n-1, pos, y)
		} else {
			l := x - 1
			r := y - 1
			res := query(1, 0, n-1, l, r)
			fmt.Fprintln(writer, res.ans)
		}
	}
}
