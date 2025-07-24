package main

import (
	"bufio"
	"fmt"
	"os"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	res := 0
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

const negInf = -1 << 60

type SegTree struct {
	size int
	val  []int
	add  []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	val := make([]int, 2*size)
	add := make([]int, 2*size)
	for i := range val {
		val[i] = negInf
	}
	return &SegTree{size: size, val: val, add: add}
}

func (st *SegTree) apply(idx, v int) {
	st.val[idx] += v
	st.add[idx] += v
}

func (st *SegTree) push(idx int) {
	if st.add[idx] != 0 {
		st.apply(idx*2, st.add[idx])
		st.apply(idx*2+1, st.add[idx])
		st.add[idx] = 0
	}
}

func (st *SegTree) rangeAdd(l, r, v, idx, lo, hi int) {
	if l > hi || r < lo {
		return
	}
	if l <= lo && hi <= r {
		st.apply(idx, v)
		return
	}
	st.push(idx)
	mid := (lo + hi) / 2
	st.rangeAdd(l, r, v, idx*2, lo, mid)
	st.rangeAdd(l, r, v, idx*2+1, mid+1, hi)
	if st.val[idx*2] > st.val[idx*2+1] {
		st.val[idx] = st.val[idx*2]
	} else {
		st.val[idx] = st.val[idx*2+1]
	}
}

func (st *SegTree) Add(l, r, v int) {
	if l <= r {
		st.rangeAdd(l, r, v, 1, 1, st.size)
	}
}

func (st *SegTree) pointSet(i, v, idx, lo, hi int) {
	if lo == hi {
		st.val[idx] = v
		st.add[idx] = 0
		return
	}
	st.push(idx)
	mid := (lo + hi) / 2
	if i <= mid {
		st.pointSet(i, v, idx*2, lo, mid)
	} else {
		st.pointSet(i, v, idx*2+1, mid+1, hi)
	}
	if st.val[idx*2] > st.val[idx*2+1] {
		st.val[idx] = st.val[idx*2]
	} else {
		st.val[idx] = st.val[idx*2+1]
	}
}

func (st *SegTree) Set(i, v int) { st.pointSet(i, v, 1, 1, st.size) }

func (st *SegTree) rangeMax(l, r, idx, lo, hi int) int {
	if l > hi || r < lo {
		return negInf
	}
	if l <= lo && hi <= r {
		return st.val[idx]
	}
	st.push(idx)
	mid := (lo + hi) / 2
	left := st.rangeMax(l, r, idx*2, lo, mid)
	right := st.rangeMax(l, r, idx*2+1, mid+1, hi)
	if left > right {
		return left
	}
	return right
}

func (st *SegTree) Max(l, r int) int { return st.rangeMax(l, r, 1, 1, st.size) }

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	order := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &order[i])
	}
	bit := NewFenwick(n)
	seg := NewSegTree(n)
	res := make([]int, n+1)
	inserted := 0
	for step := 0; step <= n; step++ {
		k := inserted
		bound := n - k
		lo, hi := 1, n
		t := 0
		for lo <= hi {
			mid := (lo + hi) / 2
			cond := mid - bit.Sum(mid)
			if cond < bound {
				t = mid
				lo = mid + 1
			} else {
				hi = mid - 1
			}
		}
		if t > 0 && bit.Sum(t) > 0 {
			val := seg.Max(1, t)
			res[step] = bound + val
		} else {
			res[step] = 1
		}
		if step == n {
			break
		}
		pos := order[step]
		bit.Add(pos, 1)
		seg.Add(pos, n, 2)
		val := 2*bit.Sum(pos) - pos
		seg.Set(pos, val)
		inserted++
	}
	out := bufio.NewWriter(os.Stdout)
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
	out.Flush()
}
