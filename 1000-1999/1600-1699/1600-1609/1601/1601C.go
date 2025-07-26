package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segTree struct {
	n    int
	min  []int
	lazy []int
}

func buildSeg(arr []int) *segTree {
	n := len(arr)
	st := &segTree{n: n, min: make([]int, 4*n), lazy: make([]int, 4*n)}
	var build func(int, int, int)
	build = func(idx, l, r int) {
		if l+1 == r {
			st.min[idx] = arr[l]
			return
		}
		m := (l + r) / 2
		build(idx*2, l, m)
		build(idx*2+1, m, r)
		if st.min[idx*2] < st.min[idx*2+1] {
			st.min[idx] = st.min[idx*2]
		} else {
			st.min[idx] = st.min[idx*2+1]
		}
	}
	build(1, 0, n)
	return st
}

func (st *segTree) push(idx int) {
	if st.lazy[idx] != 0 {
		for _, c := range []int{idx * 2, idx*2 + 1} {
			st.min[c] += st.lazy[idx]
			st.lazy[c] += st.lazy[idx]
		}
		st.lazy[idx] = 0
	}
}

func (st *segTree) rangeAdd(idx, l, r, ql, qr, val int) {
	if qr <= l || r <= ql {
		return
	}
	if ql <= l && r <= qr {
		st.min[idx] += val
		st.lazy[idx] += val
		return
	}
	st.push(idx)
	m := (l + r) / 2
	st.rangeAdd(idx*2, l, m, ql, qr, val)
	st.rangeAdd(idx*2+1, m, r, ql, qr, val)
	if st.min[idx*2] < st.min[idx*2+1] {
		st.min[idx] = st.min[idx*2]
	} else {
		st.min[idx] = st.min[idx*2+1]
	}
}

func (st *segTree) addSuffix(pos int, val int) {
	st.rangeAdd(1, 0, st.n, pos, st.n, val)
}

func (st *segTree) queryMin() int {
	return st.min[1]
}

// Fenwick tree for inversion count

type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+2)}
}

func (f *fenwick) add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	s := 0
	for idx > 0 {
		s += f.bit[idx]
		idx -= idx & -idx
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}

		// inversion count of A
		vals := make([]int, n)
		copy(vals, a)
		sort.Ints(vals)
		comp := make(map[int]int, len(vals))
		idx := 1
		for _, v := range vals {
			if _, ok := comp[v]; !ok {
				comp[v] = idx
				idx++
			}
		}
		fw := newFenwick(len(comp) + 2)
		invA := 0
		for i := n - 1; i >= 0; i-- {
			c := comp[a[i]]
			invA += fw.sum(c - 1)
			fw.add(c, 1)
		}

		type pair struct{ val, idx int }
		pairs := make([]pair, n)
		for i := 0; i < n; i++ {
			pairs[i] = pair{a[i], i + 1}
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })

		arr := make([]int, n+1)
		for i := range arr {
			arr[i] = i
		}
		st := buildSeg(arr)

		lessPtr, eqPtr := 0, 0
		sort.Ints(b)
		ans := invA
		for _, x := range b {
			for lessPtr < n && pairs[lessPtr].val < x {
				st.addSuffix(pairs[lessPtr].idx, -1)
				lessPtr++
			}
			for eqPtr < n && pairs[eqPtr].val <= x {
				st.addSuffix(pairs[eqPtr].idx, -1)
				eqPtr++
			}
			ans += st.queryMin() + lessPtr
		}
		fmt.Fprintln(writer, ans)
	}
}
