package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n   int
	max []int64
	min []int64
}

func NewSegTree(arr []int64) *SegTree {
	n := len(arr)
	st := &SegTree{n: n, max: make([]int64, 4*n), min: make([]int64, 4*n)}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int64) {
	if l == r {
		st.max[node] = arr[l]
		st.min[node] = arr[l]
		return
	}
	m := (l + r) / 2
	st.build(node*2, l, m, arr)
	st.build(node*2+1, m+1, r, arr)
	if st.max[node*2] > st.max[node*2+1] {
		st.max[node] = st.max[node*2]
	} else {
		st.max[node] = st.max[node*2+1]
	}
	if st.min[node*2] < st.min[node*2+1] {
		st.min[node] = st.min[node*2]
	} else {
		st.min[node] = st.min[node*2+1]
	}
}

func (st *SegTree) queryMax(node, l, r, ql, qr int) int64 {
	if ql <= l && r <= qr {
		return st.max[node]
	}
	m := (l + r) / 2
	var res int64 = -1 << 63
	if ql <= m {
		val := st.queryMax(node*2, l, m, ql, qr)
		if val > res {
			res = val
		}
	}
	if qr > m {
		val := st.queryMax(node*2+1, m+1, r, ql, qr)
		if val > res {
			res = val
		}
	}
	return res
}

func (st *SegTree) queryMin(node, l, r, ql, qr int) int64 {
	if ql <= l && r <= qr {
		return st.min[node]
	}
	m := (l + r) / 2
	res := int64(1<<63 - 1)
	if ql <= m {
		val := st.queryMin(node*2, l, m, ql, qr)
		if val < res {
			res = val
		}
	}
	if qr > m {
		val := st.queryMin(node*2+1, m+1, r, ql, qr)
		if val < res {
			res = val
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ps := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			ps[i] = ps[i-1] + a[i]
		}
		seg := NewSegTree(ps)
		left := make([]int, n+1)
		stack := []int{}
		for i := 1; i <= n; i++ {
			for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				left[i] = 0
			} else {
				left[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		right := make([]int, n+1)
		stack = stack[:0]
		for i := n; i >= 1; i-- {
			for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				right[i] = n + 1
			} else {
				right[i] = stack[len(stack)-1]
			}
			stack = append(stack, i)
		}
		ok := true
		for i := 1; i <= n && ok; i++ {
			L := left[i] + 1
			R := right[i] - 1
			if L > R {
				continue
			}
			maxRight := seg.queryMax(1, 0, n, i, R)
			minLeft := seg.queryMin(1, 0, n, L-1, i-1)
			if maxRight-minLeft > a[i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
