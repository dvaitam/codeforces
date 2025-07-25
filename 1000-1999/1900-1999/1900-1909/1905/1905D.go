package main

import (
	"bufio"
	"fmt"
	"os"
)

type segtree struct {
	n    int
	sum  []int64
	lazy []int64
	set  []bool
}

func NewSegTree(arr []int64) *segtree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &segtree{
		n:    n,
		sum:  make([]int64, 2*n),
		lazy: make([]int64, 2*n),
		set:  make([]bool, 2*n),
	}
	for i, v := range arr {
		st.sum[n+i] = v
	}
	for i := n - 1; i > 0; i-- {
		st.sum[i] = st.sum[i<<1] + st.sum[i<<1|1]
	}
	return st
}

func (st *segtree) apply(idx int, val int64, length int) {
	st.sum[idx] = val * int64(length)
	st.lazy[idx] = val
	st.set[idx] = true
}

func (st *segtree) push(idx, length int) {
	if st.set[idx] {
		st.apply(idx<<1, st.lazy[idx], length>>1)
		st.apply(idx<<1|1, st.lazy[idx], length>>1)
		st.set[idx] = false
	}
}

func (st *segtree) update(idx, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(idx, val, r-l+1)
		return
	}
	st.push(idx, r-l+1)
	mid := (l + r) >> 1
	st.update(idx<<1, l, mid, ql, qr, val)
	st.update(idx<<1|1, mid+1, r, ql, qr, val)
	st.sum[idx] = st.sum[idx<<1] + st.sum[idx<<1|1]
}

func (st *segtree) Update(l, r int, val int64) {
	if l > r {
		return
	}
	st.update(1, 0, st.n-1, l, r, val)
}

func (st *segtree) QueryAll() int64 {
	return st.sum[1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		pos := make([]int, n)
		for i, v := range p {
			pos[v] = i
		}
		maxPos := make([]int, n+1)
		mx := -1
		for k := 1; k <= n; k++ {
			if pos[k-1] > mx {
				mx = pos[k-1]
			}
			maxPos[k] = mx
		}
		arr := make([]int64, n)
		for k := 1; k <= n; k++ {
			arr[k-1] = int64(maxPos[k] - n)
		}
		st := NewSegTree(arr)
		sumLast := st.QueryAll()
		best := int64(0)
		for s := 0; s < n; s++ {
			cost := int64(s)*int64(n) - sumLast
			if cost > best {
				best = cost
			}
			v := p[s]
			if v+1 <= n {
				st.Update(v, n-1, int64(s))
				sumLast = st.QueryAll()
			}
		}
		fmt.Fprintln(out, best)
	}
}
