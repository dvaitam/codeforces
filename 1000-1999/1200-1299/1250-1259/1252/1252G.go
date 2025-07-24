package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n    int
	tree []int
	lazy []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr)
	st := &SegTree{n: n, tree: make([]int, 4*n), lazy: make([]int, 4*n)}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int) {
	if l == r {
		st.tree[node] = arr[l]
	} else {
		m := (l + r) / 2
		st.build(node*2, l, m, arr)
		st.build(node*2+1, m+1, r, arr)
		if st.tree[node*2] < st.tree[node*2+1] {
			st.tree[node] = st.tree[node*2]
		} else {
			st.tree[node] = st.tree[node*2+1]
		}
	}
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		val := st.lazy[node]
		st.tree[node*2] += val
		st.tree[node*2+1] += val
		st.lazy[node*2] += val
		st.lazy[node*2+1] += val
		st.lazy[node] = 0
	}
}

func (st *SegTree) rangeAdd(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.tree[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	m := (l + r) / 2
	st.rangeAdd(node*2, l, m, ql, qr, val)
	st.rangeAdd(node*2+1, m+1, r, ql, qr, val)
	if st.tree[node*2] < st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *SegTree) RangeAdd(l, r, val int) {
	if l > r {
		return
	}
	st.rangeAdd(1, 0, st.n-1, l, r, val)
}

func (st *SegTree) Min() int {
	return st.tree[1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M, Q int
	if _, err := fmt.Fscan(in, &N, &M, &Q); err != nil {
		return
	}

	A := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &A[i])
	}
	randPerf := A[0]
	lessCount := 0
	for i := 1; i < N; i++ {
		if A[i] < randPerf {
			lessCount++
		}
	}

	R := make([]int, M)
	B := make([][]int, M)
	bCount := make([]int, M)

	for i := 0; i < M; i++ {
		fmt.Fscan(in, &R[i])
		Bi := make([]int, R[i])
		cnt := 0
		for j := 0; j < R[i]; j++ {
			fmt.Fscan(in, &Bi[j])
			if Bi[j] < randPerf {
				cnt++
			}
		}
		B[i] = Bi
		bCount[i] = cnt
	}

	D := make([]int, M)
	for i := 0; i < M; i++ {
		D[i] = bCount[i] - R[i]
	}

	arr := make([]int, M)
	prefix := lessCount
	for i := 0; i < M; i++ {
		arr[i] = prefix - R[i]
		prefix += D[i]
	}

	st := NewSegTree(arr)

	for t := 0; t < Q; t++ {
		var X, Y, Z int
		fmt.Fscan(in, &X, &Y, &Z)
		X--
		Y--
		old := B[X][Y]
		oldLow := 0
		if old < randPerf {
			oldLow = 1
		}
		newLow := 0
		if Z < randPerf {
			newLow = 1
		}
		B[X][Y] = Z
		delta := newLow - oldLow
		if delta != 0 {
			bCount[X] += delta
			D[X] += delta
			if X+1 < M {
				st.RangeAdd(X+1, M-1, delta)
			}
		}
		if st.Min() >= 0 {
			fmt.Fprintln(out, 1)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}
