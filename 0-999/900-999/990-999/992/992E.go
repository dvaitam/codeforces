package main

import (
	"bufio"
	"fmt"
	"os"
)

type segTree struct {
	n    int
	min  []int64
	max  []int64
	lazy []int64
}

func newSegTree(a []int64) *segTree {
	n := len(a)
	size := 1
	for size < n {
		size <<= 1
	}
	st := &segTree{n: n, min: make([]int64, size*2), max: make([]int64, size*2), lazy: make([]int64, size*2)}
	var build func(int, int, int)
	build = func(idx, l, r int) {
		if l == r {
			val := int64(0)
			if l <= n {
				val = a[l-1]
			}
			st.min[idx] = val
			st.max[idx] = val
			return
		}
		m := (l + r) / 2
		build(idx*2, l, m)
		build(idx*2+1, m+1, r)
		st.pull(idx)
	}
	build(1, 1, size)
	return st
}

func (st *segTree) pull(idx int) {
	l := idx * 2
	r := l + 1
	if st.min[l] < st.min[r] {
		st.min[idx] = st.min[l]
	} else {
		st.min[idx] = st.min[r]
	}
	if st.max[l] > st.max[r] {
		st.max[idx] = st.max[l]
	} else {
		st.max[idx] = st.max[r]
	}
}

func (st *segTree) push(idx int) {
	if st.lazy[idx] != 0 {
		val := st.lazy[idx]
		l := idx * 2
		r := l + 1
		st.apply(l, val)
		st.apply(r, val)
		st.lazy[idx] = 0
	}
}

func (st *segTree) apply(idx int, val int64) {
	st.min[idx] += val
	st.max[idx] += val
	st.lazy[idx] += val
}

func (st *segTree) rangeAdd(idx, l, r, ql, qr int, val int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(idx, val)
		return
	}
	st.push(idx)
	m := (l + r) / 2
	st.rangeAdd(idx*2, l, m, ql, qr, val)
	st.rangeAdd(idx*2+1, m+1, r, ql, qr, val)
	st.pull(idx)
}

func (st *segTree) addRange(l, r int, val int64) {
	st.rangeAdd(1, 1, st.size(), l, r, val)
}

func (st *segTree) size() int { return len(st.min) / 2 }

func (st *segTree) findZero(idx, l, r int) int {
	if st.min[idx] > 0 || st.max[idx] < 0 {
		return -1
	}
	if l == r {
		if st.min[idx] == 0 {
			if l <= st.n {
				return l
			}
		}
		return -1
	}
	st.push(idx)
	m := (l + r) / 2
	res := st.findZero(idx*2, l, m)
	if res != -1 {
		return res
	}
	return st.findZero(idx*2+1, m+1, r)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefix := int64(0)
	d := make([]int64, n)
	for i := 0; i < n; i++ {
		d[i] = a[i] - prefix
		prefix += a[i]
	}

	st := newSegTree(d)

	for ; q > 0; q-- {
		var p int
		var x int64
		fmt.Fscan(in, &p, &x)
		old := a[p-1]
		if old != x {
			delta := x - old
			a[p-1] = x
			st.addRange(p, p, delta)
			if p < n {
				st.addRange(p+1, n, -delta)
			}
		}
		ans := st.findZero(1, 1, st.size())
		if ans == -1 || ans > n {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
