package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	n    int
	sum  []int
	lazy []int
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1
	st := &SegTree{
		n:    n,
		sum:  make([]int, 4*n+5),
		lazy: make([]int, 4*n+5),
	}
	for i := range st.lazy {
		st.lazy[i] = -1
	}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int) {
	if l == r {
		st.sum[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, arr)
	st.build(node*2+1, mid+1, r, arr)
	st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
}

func (st *SegTree) push(node, l, r int) {
	if st.lazy[node] == -1 || l == r {
		return
	}
	val := st.lazy[node]
	mid := (l + r) / 2
	st.sum[node*2] = (mid - l + 1) * val
	st.sum[node*2+1] = (r - mid) * val
	st.lazy[node*2] = val
	st.lazy[node*2+1] = val
	st.lazy[node] = -1
}

func (st *SegTree) update(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.sum[node] = (r - l + 1) * val
		st.lazy[node] = val
		return
	}
	st.push(node, l, r)
	mid := (l + r) / 2
	st.update(node*2, l, mid, ql, qr, val)
	st.update(node*2+1, mid+1, r, ql, qr, val)
	st.sum[node] = st.sum[node*2] + st.sum[node*2+1]
}

func (st *SegTree) query(node, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return st.sum[node]
	}
	st.push(node, l, r)
	mid := (l + r) / 2
	res := 0
	if ql <= mid {
		res += st.query(node*2, l, mid, ql, qr)
	}
	if qr > mid {
		res += st.query(node*2+1, mid+1, r, ql, qr)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		var sStr, fStr string
		fmt.Fscan(in, &sStr)
		fmt.Fscan(in, &fStr)
		L := make([]int, q)
		R := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &L[i], &R[i])
		}
		arr := make([]int, n+1)
		for i := 0; i < n; i++ {
			if fStr[i] == '1' {
				arr[i+1] = 1
			} else {
				arr[i+1] = 0
			}
		}
		st := NewSegTree(arr)
		ok := true
		for i := q - 1; i >= 0 && ok; i-- {
			l := L[i]
			r := R[i]
			ones := st.query(1, 1, n, l, r)
			length := r - l + 1
			if ones*2 == length {
				ok = false
				break
			}
			if ones*2 > length {
				st.update(1, 1, n, l, r, 1)
			} else {
				st.update(1, 1, n, l, r, 0)
			}
		}
		if ok {
			for i := 0; i < n; i++ {
				val := st.query(1, 1, n, i+1, i+1)
				if val != int(sStr[i]-'0') {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
