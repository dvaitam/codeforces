package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type segTree struct {
	n    int
	tree []int
	lazy []int
}

func newSegTree(n int) *segTree {
	return &segTree{n: n, tree: make([]int, 4*n), lazy: make([]int, 4*n)}
}

func (st *segTree) apply(node, val int) {
	st.tree[node] += val
	st.lazy[node] += val
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		v := st.lazy[node]
		st.apply(node*2, v)
		st.apply(node*2+1, v)
		st.lazy[node] = 0
	}
}

func (st *segTree) rangeAdd(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(node, val)
		return
	}
	st.push(node)
	m := (l + r) / 2
	st.rangeAdd(node*2, l, m, ql, qr, val)
	st.rangeAdd(node*2+1, m+1, r, ql, qr, val)
	if st.tree[node*2] > st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *segTree) Add(l, r, val int) {
	if l > r {
		return
	}
	st.rangeAdd(1, 0, st.n-1, l, r, val)
}

func (st *segTree) Max() int {
	return st.tree[1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		segs := make([]struct{ l, r, li, ri int }, n)
		coords := []int{1, m + 1}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
			coords = append(coords, segs[i].l, segs[i].r+1)
		}
		sort.Ints(coords)
		coords = unique(coords)
		idx := func(x int) int {
			return sort.SearchInts(coords, x)
		}
		k := len(coords) - 1
		for i := 0; i < n; i++ {
			segs[i].li = idx(segs[i].l)
			segs[i].ri = idx(segs[i].r+1) - 1
		}

		segsByR := make([]struct{ l, r, li, ri int }, n)
		copy(segsByR, segs)
		sort.Slice(segsByR, func(i, j int) bool { return segsByR[i].r < segsByR[j].r })

		pref := make([]int, len(coords))
		st := newSegTree(k)
		p := 0
		for i := 0; i < len(coords); i++ {
			y := coords[i]
			for p < n && segsByR[p].r < y {
				st.Add(segsByR[p].li, segsByR[p].ri, 1)
				p++
			}
			if k > 0 {
				pref[i] = st.Max()
			}
		}

		segsByL := make([]struct{ l, r, li, ri int }, n)
		copy(segsByL, segs)
		sort.Slice(segsByL, func(i, j int) bool { return segsByL[i].l > segsByL[j].l })
		suf := make([]int, len(coords))
		st2 := newSegTree(k)
		p = 0
		for i := len(coords) - 1; i >= 0; i-- {
			y := coords[i]
			for p < n && segsByL[p].l > y {
				st2.Add(segsByL[p].li, segsByL[p].ri, 1)
				p++
			}
			if k > 0 {
				suf[i] = st2.Max()
			}
		}

		ans := 0
		for i := 0; i < len(coords) && coords[i] <= m; i++ {
			if pref[i] > ans {
				ans = pref[i]
			}
			if suf[i] > ans {
				ans = suf[i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func unique(a []int) []int {
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}
