package main

import (
	"bufio"
	"fmt"
	"os"
)

// Segment tree with lazy propagation supporting range add and finding the
// earliest index with non-zero value. It also supports an O(1) clear using
// a global timestamp.
type SegTree struct {
	n    int
	min  []int64
	max  []int64
	add  []int64
	time []int
	cur  int
}

func NewSegTree(n int) *SegTree {
	size := 4 * (n + 2)
	return &SegTree{
		n:    n,
		min:  make([]int64, size),
		max:  make([]int64, size),
		add:  make([]int64, size),
		time: make([]int, size),
		cur:  1,
	}
}

func (st *SegTree) touch(i int) {
	if st.time[i] != st.cur {
		st.time[i] = st.cur
		st.min[i] = 0
		st.max[i] = 0
		st.add[i] = 0
	}
}

func (st *SegTree) push(i int) {
	st.touch(i)
	if st.add[i] != 0 {
		val := st.add[i]
		left, right := i<<1, i<<1|1
		st.touch(left)
		st.touch(right)
		st.add[left] += val
		st.min[left] += val
		st.max[left] += val
		st.add[right] += val
		st.min[right] += val
		st.max[right] += val
		st.add[i] = 0
	}
}

func (st *SegTree) pull(i int) {
	left, right := i<<1, i<<1|1
	if st.min[left] < st.min[right] {
		st.min[i] = st.min[left]
	} else {
		st.min[i] = st.min[right]
	}
	if st.max[left] > st.max[right] {
		st.max[i] = st.max[left]
	} else {
		st.max[i] = st.max[right]
	}
}

func (st *SegTree) rangeAdd(i, l, r, ql, qr int, val int64) {
	st.touch(i)
	if ql <= l && r <= qr {
		st.add[i] += val
		st.min[i] += val
		st.max[i] += val
		return
	}
	st.push(i)
	mid := (l + r) >> 1
	if ql <= mid {
		st.rangeAdd(i<<1, l, mid, ql, qr, val)
	}
	if qr > mid {
		st.rangeAdd(i<<1|1, mid+1, r, ql, qr, val)
	}
	st.pull(i)
}

func (st *SegTree) firstNonZero(i, l, r int) (int, int64, bool) {
	st.touch(i)
	if st.min[i] == 0 && st.max[i] == 0 {
		return 0, 0, false
	}
	if l == r {
		return l, st.min[i], true
	}
	st.push(i)
	mid := (l + r) >> 1
	if idx, val, ok := st.firstNonZero(i<<1, l, mid); ok {
		return idx, val, true
	}
	return st.firstNonZero(i<<1|1, mid+1, r)
}

func (st *SegTree) Clear() {
	st.cur++
	st.touch(1)
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
		var q int
		fmt.Fscan(in, &q)
		ops := make([]struct {
			l, r int
			x    int64
		}, q)
		seg := NewSegTree(n)
		best := 0
		for j := 1; j <= q; j++ {
			var l, r int
			var x int64
			fmt.Fscan(in, &l, &r, &x)
			ops[j-1] = struct {
				l, r int
				x    int64
			}{l, r, x}
			seg.rangeAdd(1, 1, n, l, r, x)
			if idx, val, ok := seg.firstNonZero(1, 1, n); ok {
				if val < 0 {
					best = j
					seg.Clear()
				} else if val > 0 {
					// best array remains
				} else {
					_ = idx // arrays equal
				}
			}
		}
		// compute array after best operations
		diff := make([]int64, n+2)
		for i := 0; i < best; i++ {
			op := ops[i]
			diff[op.l] += op.x
			diff[op.r+1] -= op.x
		}
		cur := int64(0)
		for i := 1; i <= n; i++ {
			cur += diff[i]
			a[i] += cur
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, a[i])
		}
		fmt.Fprintln(out)
	}
}
