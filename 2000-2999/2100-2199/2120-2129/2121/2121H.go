package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segTree struct {
	n    int
	seg  []int
	lazy []int
}

func newSegTree(n int) *segTree {
	return &segTree{
		n:    n,
		seg:  make([]int, 4*n),
		lazy: make([]int, 4*n),
	}
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		delta := st.lazy[node]
		child := node * 2
		st.seg[child] += delta
		st.lazy[child] += delta
		child++
		st.seg[child] += delta
		st.lazy[child] += delta
		st.lazy[node] = 0
	}
}

func (st *segTree) rangeAdd(l, r, val int) {
	if l > r {
		return
	}
	st.rangeAddRec(1, 0, st.n-1, l, r, val)
}

func (st *segTree) rangeAddRec(node, nl, nr, l, r, val int) {
	if r < nl || nr < l {
		return
	}
	if l <= nl && nr <= r {
		st.seg[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (nl + nr) / 2
	st.rangeAddRec(node*2, nl, mid, l, r, val)
	st.rangeAddRec(node*2+1, mid+1, nr, l, r, val)
	if st.seg[node*2] > st.seg[node*2+1] {
		st.seg[node] = st.seg[node*2]
	} else {
		st.seg[node] = st.seg[node*2+1]
	}
}

func (st *segTree) rangeMax(l, r int) int {
	if l > r {
		return 0
	}
	return st.rangeMaxRec(1, 0, st.n-1, l, r)
}

func (st *segTree) rangeMaxRec(node, nl, nr, l, r int) int {
	if r < nl || nr < l {
		return 0
	}
	if l <= nl && nr <= r {
		return st.seg[node]
	}
	st.push(node)
	mid := (nl + nr) / 2
	left := st.rangeMaxRec(node*2, nl, mid, l, r)
	right := st.rangeMaxRec(node*2+1, mid+1, nr, l, r)
	if left > right {
		return left
	}
	return right
}

func uniqueSorted(a []int) []int {
	if len(a) == 0 {
		return a
	}
	w := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[w-1] {
			a[w] = a[i]
			w++
		}
	}
	return a[:w]
}

func upperBound(a []int, target int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] > target })
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		l := make([]int, n)
		r := make([]int, n)
		vals := make([]int, 1, n+1)
		vals[0] = 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
			vals = append(vals, l[i])
		}
		sort.Ints(vals)
		vals = uniqueSorted(vals)
		m := len(vals)
		index := make(map[int]int, m)
		for i, v := range vals {
			index[v] = i
		}
		st := newSegTree(m)
		ans := make([]int, n)
		for i := 0; i < n; i++ {
			li := l[i]
			ri := r[i]
			idxL := index[li]
			idxR := upperBound(vals, ri) - 1
			if idxR >= m {
				idxR = m - 1
			}
			idxLimit := idxL - 1
			if idxLimit > idxR {
				idxLimit = idxR
			}
			mx := 0
			if idxLimit >= 0 {
				mx = st.rangeMax(0, idxLimit)
			}
			if idxR >= idxL {
				st.rangeAdd(idxL, idxR, 1)
			}
			cand := mx + 1
			cur := st.rangeMax(idxL, idxL)
			if cand > cur {
				st.rangeAdd(idxL, idxL, cand-cur)
			}
			ans[i] = st.seg[1]
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[i])
		}
		fmt.Fprintln(writer)
	}
}
