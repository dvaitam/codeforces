package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -(1 << 60)

type segTree struct {
	n    int
	max  []int64
	lazy []int64
}

func newSegTree(arr []int64) *segTree {
	n := len(arr)
	st := &segTree{
		n:    n,
		max:  make([]int64, n*4),
		lazy: make([]int64, n*4),
	}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segTree) build(node, l, r int, arr []int64) {
	if l == r {
		st.max[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, arr)
	st.build(node*2+1, mid+1, r, arr)
	st.max[node] = max64(st.max[node*2], st.max[node*2+1])
}

func (st *segTree) push(node int) {
	if st.lazy[node] == 0 {
		return
	}
	delta := st.lazy[node]
	left := node * 2
	right := left + 1
	st.max[left] += delta
	st.lazy[left] += delta
	st.max[right] += delta
	st.lazy[right] += delta
	st.lazy[node] = 0
}

func (st *segTree) RangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	st.rangeAdd(1, 0, st.n-1, l, r, val)
}

func (st *segTree) rangeAdd(node, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.max[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (l + r) / 2
	st.rangeAdd(node*2, l, mid, ql, qr, val)
	st.rangeAdd(node*2+1, mid+1, r, ql, qr, val)
	st.max[node] = max64(st.max[node*2], st.max[node*2+1])
}

func (st *segTree) RangeMax(l, r int) int64 {
	if l > r {
		return negInf
	}
	return st.rangeMax(1, 0, st.n-1, l, r)
}

func (st *segTree) rangeMax(node, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return negInf
	}
	if ql <= l && r <= qr {
		return st.max[node]
	}
	st.push(node)
	mid := (l + r) / 2
	left := st.rangeMax(node*2, l, mid, ql, qr)
	right := st.rangeMax(node*2+1, mid+1, r, ql, qr)
	if left > right {
		return left
	}
	return right
}

func (st *segTree) GetValue(pos int) int64 {
	return st.getValue(1, 0, st.n-1, pos)
}

func (st *segTree) getValue(node, l, r, pos int) int64 {
	if l == r {
		return st.max[node]
	}
	st.push(node)
	mid := (l + r) / 2
	if pos <= mid {
		return st.getValue(node*2, l, mid, pos)
	}
	return st.getValue(node*2+1, mid+1, r, pos)
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) NextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for err == nil && c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		values := make([]int64, n)
		for i := 0; i < n; i++ {
			values[i] = int64(fs.NextInt())
		}
		alice := make([]int, n)
		for i := 0; i < n; i++ {
			alice[i] = fs.NextInt()
		}
		posB := make([]int, n+1)
		for i := 0; i < n; i++ {
			item := fs.NextInt()
			posB[item] = i + 1
		}

		arr := make([]int64, n+1)
		arr[0] = 0
		for i := 1; i <= n; i++ {
			arr[i] = negInf
		}
		st := newSegTree(arr)

		for _, item := range alice {
			r := posB[item]
			val := values[item-1]
			best := st.RangeMax(0, r)
			curr := st.GetValue(r)
			st.RangeAdd(r, r, best-curr)
			st.RangeAdd(0, r-1, val)
		}
		ans := st.RangeMax(0, n)
		fmt.Fprintln(out, ans)
	}
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
