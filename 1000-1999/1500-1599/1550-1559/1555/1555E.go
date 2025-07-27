package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type SegTree struct {
	n    int
	min  []int
	lazy []int
}

func NewSegTree(n int) *SegTree {
	size := 4 * n
	st := &SegTree{n: n, min: make([]int, size), lazy: make([]int, size)}
	return st
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		val := st.lazy[node]
		left := node << 1
		right := left | 1
		st.min[left] += val
		st.min[right] += val
		st.lazy[left] += val
		st.lazy[right] += val
		st.lazy[node] = 0
	}
}

func (st *SegTree) add(node, l, r, ql, qr, val int) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.min[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	m := (l + r) >> 1
	st.add(node<<1, l, m, ql, qr, val)
	st.add(node<<1|1, m+1, r, ql, qr, val)
	if st.min[node<<1] < st.min[node<<1|1] {
		st.min[node] = st.min[node<<1]
	} else {
		st.min[node] = st.min[node<<1|1]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	segs := make([][3]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &segs[i][0], &segs[i][1], &segs[i][2])
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i][2] < segs[j][2] })

	if m == 1 {
		fmt.Fprintln(out, 0)
		return
	}

	st := NewSegTree(m - 1)
	ans := int(1e9 + 7)
	l := 0
	for r := 0; r < n; r++ {
		st.add(1, 1, m-1, segs[r][0], segs[r][1]-1, 1)
		for l <= r && st.min[1] > 0 {
			diff := segs[r][2] - segs[l][2]
			if diff < ans {
				ans = diff
			}
			st.add(1, 1, m-1, segs[l][0], segs[l][1]-1, -1)
			l++
		}
	}
	fmt.Fprintln(out, ans)
}
