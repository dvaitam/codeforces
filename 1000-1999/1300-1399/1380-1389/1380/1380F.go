package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

type Matrix struct {
	a11, a12, a21, a22 int64
}

func mul(a, b Matrix) Matrix {
	return Matrix{
		a11: (a.a11*b.a11 + a.a12*b.a21) % MOD,
		a12: (a.a11*b.a12 + a.a12*b.a22) % MOD,
		a21: (a.a21*b.a11 + a.a22*b.a21) % MOD,
		a22: (a.a21*b.a12 + a.a22*b.a22) % MOD,
	}
}

func countSum(s int) int64 {
	if s <= 9 {
		return int64(s + 1)
	}
	return int64(19 - s)
}

type SegTree struct {
	n    int
	tree []Matrix
	arr  []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1 // arr is 1-indexed
	st := &SegTree{n: n, tree: make([]Matrix, 4*n+4), arr: arr}
	st.build(1, 1, n)
	return st
}

func (st *SegTree) build(node, l, r int) {
	if l == r {
		st.tree[node] = st.makeMatrix(l)
		return
	}
	m := (l + r) >> 1
	st.build(node<<1, l, m)
	st.build(node<<1|1, m+1, r)
	st.tree[node] = mul(st.tree[node<<1|1], st.tree[node<<1])
}

func (st *SegTree) makeMatrix(i int) Matrix {
	cnt1 := countSum(st.arr[i])
	cnt2 := int64(0)
	if i > 1 {
		val := st.arr[i-1]*10 + st.arr[i]
		if val >= 10 && val <= 18 {
			cnt2 = countSum(val)
		}
	}
	return Matrix{a11: cnt1 % MOD, a12: cnt2 % MOD, a21: 1, a22: 0}
}

func (st *SegTree) update(pos int) {
	st.updateRec(1, 1, st.n, pos)
}

func (st *SegTree) updateRec(node, l, r, pos int) {
	if l == r {
		st.tree[node] = st.makeMatrix(l)
		return
	}
	m := (l + r) >> 1
	if pos <= m {
		st.updateRec(node<<1, l, m, pos)
	} else {
		st.updateRec(node<<1|1, m+1, r, pos)
	}
	st.tree[node] = mul(st.tree[node<<1|1], st.tree[node<<1])
}

func (st *SegTree) query() Matrix {
	return st.tree[1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = int(s[i-1] - '0')
	}
	st := NewSegTree(arr)
	for ; m > 0; m-- {
		var x, d int
		fmt.Fscan(reader, &x, &d)
		arr[x] = d
		st.update(x)
		if x < n {
			st.update(x + 1)
		}
		res := st.query()
		fmt.Fprintln(writer, res.a11%MOD)
	}
}
