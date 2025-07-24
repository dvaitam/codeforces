package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

type matrix struct {
	a, b, c, d int64
}

type node struct {
	mat  matrix
	lazy bool
}

var (
	n   int
	q   int
	s   string
	seg []node
)

func mul(x, y matrix) matrix {
	return matrix{
		a: (x.a*y.a + x.b*y.c) % mod,
		b: (x.a*y.b + x.b*y.d) % mod,
		c: (x.c*y.a + x.d*y.c) % mod,
		d: (x.c*y.b + x.d*y.d) % mod,
	}
}

func toggle(m matrix) matrix {
	return matrix{a: m.d, b: m.c, c: m.b, d: m.a}
}

func apply(idx int) {
	seg[idx].mat = toggle(seg[idx].mat)
	seg[idx].lazy = !seg[idx].lazy
}

func build(idx, l, r int) {
	seg[idx].lazy = false
	if l == r {
		if s[l] == 'A' {
			seg[idx].mat = matrix{1, 1, 0, 1}
		} else {
			seg[idx].mat = matrix{1, 0, 1, 1}
		}
		return
	}
	m := (l + r) >> 1
	build(idx<<1, l, m)
	build(idx<<1|1, m+1, r)
	seg[idx].mat = mul(seg[idx<<1|1].mat, seg[idx<<1].mat)
}

func push(idx int) {
	if seg[idx].lazy {
		apply(idx << 1)
		apply(idx<<1 | 1)
		seg[idx].lazy = false
	}
}

func update(idx, l, r, ql, qr int) {
	if ql <= l && r <= qr {
		apply(idx)
		return
	}
	push(idx)
	m := (l + r) >> 1
	if ql <= m {
		update(idx<<1, l, m, ql, qr)
	}
	if qr > m {
		update(idx<<1|1, m+1, r, ql, qr)
	}
	seg[idx].mat = mul(seg[idx<<1|1].mat, seg[idx<<1].mat)
}

func query(idx, l, r, ql, qr int) matrix {
	if ql <= l && r <= qr {
		return seg[idx].mat
	}
	push(idx)
	m := (l + r) >> 1
	if qr <= m {
		return query(idx<<1, l, m, ql, qr)
	}
	if ql > m {
		return query(idx<<1|1, m+1, r, ql, qr)
	}
	left := query(idx<<1, l, m, ql, m)
	right := query(idx<<1|1, m+1, r, m+1, qr)
	return mul(right, left)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	fmt.Fscan(reader, &s)

	seg = make([]node, 4*n)
	build(1, 0, n-1)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			update(1, 0, n-1, l, r)
		} else {
			var l, r int
			var A, B int64
			fmt.Fscan(reader, &l, &r, &A, &B)
			l--
			r--
			mat := query(1, 0, n-1, l, r)
			resA := (mat.a*A + mat.b*B) % mod
			resB := (mat.c*A + mat.d*B) % mod
			fmt.Fprintf(writer, "%d %d\n", resA, resB)
		}
	}
}
