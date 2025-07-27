package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	val       int
	L         int
	R         int
	maxInside int
}

type SegTree struct {
	n    int
	data []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegTree{n: size, data: make([]int, size*2)}
}

func (st *SegTree) Update(pos, val int) {
	idx := pos + st.n - 1
	if val > st.data[idx] {
		st.data[idx] = val
		idx >>= 1
		for idx > 0 {
			left, right := st.data[idx*2], st.data[idx*2+1]
			if left > right {
				st.data[idx] = left
			} else {
				st.data[idx] = right
			}
			idx >>= 1
		}
	}
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.n - 1
	r += st.n - 1
	res := 0
	for l <= r {
		if l&1 == 1 {
			if st.data[l] > res {
				res = st.data[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.data[r] > res {
				res = st.data[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		first := make(map[int]int)
		last := make(map[int]int)
		idxs := make(map[int][]int)
		for i, v := range a {
			if _, ok := first[v]; !ok {
				first[v] = i + 1
			}
			last[v] = i + 1
			idxs[v] = append(idxs[v], i)
		}
		items := make([]Item, 0, len(first))
		for v, L := range first {
			items = append(items, Item{val: v, L: L, R: last[v]})
		}
		sort.Slice(items, func(i, j int) bool {
			if items[i].R == items[j].R {
				return items[i].val > items[j].val
			}
			return items[i].R < items[j].R
		})
		st := NewSegTree(n)
		maxInside := make(map[int]int)
		for i := range items {
			it := &items[i]
			it.maxInside = st.Query(it.L, it.R)
			maxInside[it.val] = it.maxInside
			st.Update(it.L, it.val)
		}
		essential := make(map[int]bool)
		for _, it := range items {
			if maxInside[it.val] <= it.val {
				essential[it.val] = true
			}
		}
		b := make([]int, n)
		for v, positions := range idxs {
			if essential[v] {
				for _, p := range positions {
					b[p] = v
				}
			}
		}
		for v, positions := range idxs {
			if essential[v] {
				continue
			}
			limit := maxInside[v]
			cand := limit - 1
			for cand > 0 && essential[cand] {
				cand--
			}
			if cand <= 0 {
				cand = 1
				for essential[cand] {
					cand++
				}
			}
			for _, p := range positions {
				b[p] = cand
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, b[i])
		}
		fmt.Fprintln(out)
	}
}
