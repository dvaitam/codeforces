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

func newSegTree(values []int) *segTree {
	n := len(values)
	st := &segTree{
		n:    n,
		min:  make([]int, 4*n),
		lazy: make([]int, 4*n),
	}
	st.build(1, 1, n, values)
	return st
}

func (st *segTree) build(node, l, r int, values []int) {
	if l == r {
		st.min[node] = values[l-1]
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid, values)
	st.build(node<<1|1, mid+1, r, values)
	st.pull(node)
}

func (st *segTree) pull(node int) {
	if st.min[node<<1] < st.min[node<<1|1] {
		st.min[node] = st.min[node<<1]
	} else {
		st.min[node] = st.min[node<<1|1]
	}
}

func (st *segTree) push(node int) {
	if st.lazy[node] != 0 {
		for _, child := range []int{node << 1, node<<1 | 1} {
			st.min[child] += st.lazy[node]
			st.lazy[child] += st.lazy[node]
		}
		st.lazy[node] = 0
	}
}

func (st *segTree) rangeAdd(l, r, val int) {
	st.rangeAddInternal(1, 1, st.n, l, r, val)
}

func (st *segTree) rangeAddInternal(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.min[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	mid := (l + r) >> 1
	st.rangeAddInternal(node<<1, l, mid, ql, qr, val)
	st.rangeAddInternal(node<<1|1, mid+1, r, ql, qr, val)
	st.pull(node)
}

func (st *segTree) rangeMin(l, r int) int {
	return st.rangeMinInternal(1, 1, st.n, l, r)
}

func (st *segTree) rangeMinInternal(node, l, r, ql, qr int) int {
	if ql > r || qr < l {
		return int(1e9)
	}
	if ql <= l && r <= qr {
		return st.min[node]
	}
	st.push(node)
	mid := (l + r) >> 1
	left := st.rangeMinInternal(node<<1, l, mid, ql, qr)
	right := st.rangeMinInternal(node<<1|1, mid+1, r, ql, qr)
	if left < right {
		return left
	}
	return right
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		d := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}

		capacity := make([]int, n)
		for i := 0; i < n; i++ {
			remaining := n - i
			capacity[i] = remaining / (k + 1)
		}

		st := newSegTree(capacity)

		indices := make([]int, n)
		for i := range indices {
			indices[i] = i
		}
		sort.Slice(indices, func(i, j int) bool {
			if d[indices[i]] == d[indices[j]] {
				return indices[i] < indices[j]
			}
			return d[indices[i]] > d[indices[j]]
		})

		var ans int64
		for _, idx := range indices {
			pos := idx + 1
			if st.rangeMin(1, pos) > 0 {
				ans += d[idx]
				st.rangeAdd(1, pos, -1)
			}
		}

		fmt.Fprintln(out, ans)
	}
}
