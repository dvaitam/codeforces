package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segTree struct {
	n   int
	tr  [][]int
	arr []int
}

func newSegTree(a []int) *segTree {
	n := len(a)
	size := 1
	for size < n {
		size <<= 1
	}
	tr := make([][]int, size<<1)
	st := &segTree{n: n, tr: tr, arr: a}
	st.build(1, 0, n-1)
	return st
}

func (st *segTree) build(v, l, r int) {
	if l == r {
		st.tr[v] = []int{st.arr[l]}
		return
	}
	mid := (l + r) >> 1
	st.build(v<<1, l, mid)
	st.build(v<<1|1, mid+1, r)
	left := st.tr[v<<1]
	right := st.tr[v<<1|1]
	merged := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			merged = append(merged, left[i])
			i++
		} else {
			merged = append(merged, right[j])
			j++
		}
	}
	merged = append(merged, left[i:]...)
	merged = append(merged, right[j:]...)
	st.tr[v] = merged
}

// count of elements < x in [ql, qr]
func (st *segTree) countLess(ql, qr, x int) int {
	return st.countLessRec(1, 0, st.n-1, ql, qr, x)
}

func (st *segTree) countLessRec(v, l, r, ql, qr, x int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return sort.SearchInts(st.tr[v], x)
	}
	mid := (l + r) >> 1
	return st.countLessRec(v<<1, l, mid, ql, qr, x) + st.countLessRec(v<<1|1, mid+1, r, ql, qr, x)
}

func processQuery(l, r, k int, pos []int, p []int, st *segTree) int {
	kPos := pos[k]
	if kPos < l || kPos > r {
		return -1
	}

	cntLess := st.countLess(l, r, k)
	lenSeg := r - l + 1
	cntGreater := lenSeg - cntLess - 1 // one element equals k

	var rg, rl int         // required > , required <
	var cg, cl, bg, bl int // correct and bad counts
	L, R := l, r
	for L <= R {
		mid := (L + R) >> 1
		if mid == kPos {
			break
		}
		if kPos < mid {
			rg++
			if p[mid] > k {
				cg++
			} else {
				bg++
			}
			R = mid - 1
		} else {
			rl++
			if p[mid] < k {
				cl++
			} else {
				bl++
			}
			L = mid + 1
		}
	}

	if cntGreater < rg || cntLess < rl {
		return -1
	}

	greaterOthers := cntGreater - cg - bl
	lessOthers := cntLess - cl - bg

	if bg > bl {
		if greaterOthers < bg-bl {
			return -1
		}
		return 2 * bg
	}
	if bl > bg {
		if lessOthers < bl-bg {
			return -1
		}
		return 2 * bl
	}
	return 2 * bg
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		p := make([]int, n+1) // 1-indexed for positions, but segTree uses 0-index
		pos := make([]int, n+1)
		arr := make([]int, n)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i
			arr[i-1] = p[i]
		}
		st := newSegTree(arr)

		for ; q > 0; q-- {
			var l, r, k int
			fmt.Fscan(in, &l, &r, &k)
			// convert to 0-based indices for segment tree queries
			ans := processQuery(l-1, r-1, k, pos, p, st)
			fmt.Fprintln(out, ans)
		}
	}
}
