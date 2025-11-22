package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	states = 5
	negInf = int64(-4e18)
)

type matrix [states][states]int64

func newMatrix() matrix {
	var m matrix
	for i := 0; i < states; i++ {
		for j := 0; j < states; j++ {
			m[i][j] = negInf
		}
	}
	return m
}

func leafMatrix(a, b int64) matrix {
	m := newMatrix()
	// stay in neutral states without touching this position
	m[0][0] = 0
	m[2][2] = 0
	m[4][4] = 0

	// start/continue/finish first segment
	m[0][1] = a + b   // start first and keep it open
	m[0][2] = a + 2*b // start and end first here
	m[1][1] = a       // continue first
	m[1][2] = a + b   // end first here

	// start/continue/finish second segment
	m[2][3] = a + b   // start second and keep it open
	m[2][4] = a + 2*b // start and end second here
	m[3][3] = a       // continue second
	m[3][4] = a + b   // end second here

	return m
}

func multiply(A, B matrix) matrix {
	res := newMatrix()
	for i := 0; i < states; i++ {
		for k := 0; k < states; k++ {
			if A[i][k] == negInf {
				continue
			}
			ak := A[i][k]
			for j := 0; j < states; j++ {
				v := ak + B[k][j]
				if v > res[i][j] {
					res[i][j] = v
				}
			}
		}
	}
	return res
}

type segTree struct {
	n    int
	tr   []matrix
	arrA []int64
	arrB []int64
}

func newSegTree(a, b []int64) *segTree {
	n := len(a)
	tr := make([]matrix, 4*n)
	st := &segTree{n: n, tr: tr, arrA: a, arrB: b}
	st.build(1, 0, n-1)
	return st
}

func (st *segTree) build(v, l, r int) {
	if l == r {
		st.tr[v] = leafMatrix(st.arrA[l], st.arrB[l])
		return
	}
	m := (l + r) >> 1
	st.build(v<<1, l, m)
	st.build(v<<1|1, m+1, r)
	st.tr[v] = multiply(st.tr[v<<1], st.tr[v<<1|1])
}

func (st *segTree) update(v, l, r, pos int, isA bool, val int64) {
	if l == r {
		if isA {
			st.arrA[pos] = val
		} else {
			st.arrB[pos] = val
		}
		st.tr[v] = leafMatrix(st.arrA[pos], st.arrB[pos])
		return
	}
	m := (l + r) >> 1
	if pos <= m {
		st.update(v<<1, l, m, pos, isA, val)
	} else {
		st.update(v<<1|1, m+1, r, pos, isA, val)
	}
	st.tr[v] = multiply(st.tr[v<<1], st.tr[v<<1|1])
}

func (st *segTree) query(v, l, r, ql, qr int) matrix {
	if ql == l && qr == r {
		return st.tr[v]
	}
	m := (l + r) >> 1
	if qr <= m {
		return st.query(v<<1, l, m, ql, qr)
	}
	if ql > m {
		return st.query(v<<1|1, m+1, r, ql, qr)
	}
	left := st.query(v<<1, l, m, ql, m)
	right := st.query(v<<1|1, m+1, r, m+1, qr)
	return multiply(left, right)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	seg := newSegTree(a, b)

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 || t == 2 {
			var p int
			var x int64
			fmt.Fscan(in, &p, &x)
			p--
			seg.update(1, 0, n-1, p, t == 1, x)
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			r--
			mat := seg.query(1, 0, n-1, l, r)
			fmt.Fprintln(out, mat[0][4])
		}
	}
}
