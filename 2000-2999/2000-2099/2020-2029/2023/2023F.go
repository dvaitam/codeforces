package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	size int
	min  []int64
	max  []int64
}

func newSegTree(arr []int64) *SegTree {
	n := len(arr)
	size := 1
	for size < n {
		size <<= 1
	}
	inf := int64(1 << 62)
	min := make([]int64, 2*size)
	max := make([]int64, 2*size)
	for i := 0; i < 2*size; i++ {
		min[i] = inf
	}
	for i := 0; i < n; i++ {
		min[size+i] = arr[i]
		max[size+i] = arr[i]
	}
	for i := size - 1; i > 0; i-- {
		left := i << 1
		right := left | 1
		if min[left] < min[right] {
			min[i] = min[left]
		} else {
			min[i] = min[right]
		}
		if max[left] > max[right] {
			max[i] = max[left]
		} else {
			max[i] = max[right]
		}
	}
	return &SegTree{size: size, min: min, max: max}
}

func (st *SegTree) queryMin(l, r int) int64 {
	l += st.size
	r += st.size + 1
	res := int64(1 << 62)
	for l < r {
		if l&1 == 1 {
			if st.min[l] < res {
				res = st.min[l]
			}
			l++
		}
		if r&1 == 1 {
			r--
			if st.min[r] < res {
				res = st.min[r]
			}
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func (st *SegTree) queryMax(l, r int) int64 {
	l += st.size
	r += st.size + 1
	res := int64(-1 << 62)
	for l < r {
		if l&1 == 1 {
			if st.max[l] > res {
				res = st.max[l]
			}
			l++
		}
		if r&1 == 1 {
			r--
			if st.max[r] > res {
				res = st.max[r]
			}
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func (st *SegTree) findFirstLE(l, r int, val int64) int {
	if st.queryMin(l, r) > val {
		return -1
	}
	l += st.size
	r += st.size + 1
	for l < r {
		if l&1 == 1 {
			if st.min[l] <= val {
				return st.findFirstLEInNode(l, val)
			}
			l++
		}
		if r&1 == 1 {
			r--
			if st.min[r] <= val {
				return st.findFirstLEInNode(r, val)
			}
		}
		l >>= 1
		r >>= 1
	}
	return -1
}

func (st *SegTree) findFirstLEInNode(node int, val int64) int {
	for node < st.size {
		left := node << 1
		right := left | 1
		if st.min[left] <= val {
			node = left
		} else {
			node = right
		}
	}
	return node - st.size
}

func (st *SegTree) findLastGE(l, r int, val int64) int {
	if st.queryMax(l, r) < val {
		return -1
	}
	l += st.size
	r += st.size + 1
	for l < r {
		if r&1 == 1 {
			r--
			if st.max[r] >= val {
				return st.findLastGEInNode(r, val)
			}
		}
		if l&1 == 1 {
			l++
		}
		l >>= 1
		r >>= 1
	}
	return -1
}

func (st *SegTree) findLastGEInNode(node int, val int64) int {
	for node < st.size {
		left := node << 1
		right := left | 1
		if st.max[right] >= val {
			node = right
		} else {
			node = left
		}
	}
	return node - st.size
}

func solveFromLeft(l, r int, pref []int64, st *SegTree) int64 {
	// Move from l to r. If prefix never dips below base, a single sweep works.
	// Otherwise, we must go from r back to the earliest position where the minimum prefix appears.
	minPref := st.queryMin(l, r)
	base := pref[l-1]
	if minPref >= base {
		return int64(r - l)
	}
	posMin := st.findFirstLE(l, r, minPref)
	if posMin == -1 {
		return int64(1 << 60)
	}
	return int64(r-l) + int64(r-posMin)
}

func solveFromRight(l, r int, pref []int64, st *SegTree) int64 {
	// Symmetric to solveFromLeft: traverse from r to l, then return to the first point
	// where the minimum suffix (i.e. maximum prefix) is attained.
	maxPref := st.queryMax(l-1, r-1)
	minSuffix := pref[r] - maxPref
	if minSuffix >= 0 {
		return int64(r - l)
	}
	posMax := st.findLastGE(l-1, r-1, maxPref)
	if posMax == -1 {
		return int64(1 << 60)
	}
	pos := posMax + 1 // convert pref index back to position
	return int64(r-l) + int64(pos-l)
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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		pref := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] + a[i-1]
		}

		tree := newSegTree(pref)

		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			total := pref[r] - pref[l-1]
			if total < 0 {
				fmt.Fprintln(out, -1)
				continue
			}

			best := solveFromLeft(l, r, pref, tree)
			fromRight := solveFromRight(l, r, pref, tree)
			if fromRight < best {
				best = fromRight
			}
			if best >= int64(1<<60) {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, best)
			}
		}
	}
}
