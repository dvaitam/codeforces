package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Rope struct {
	l int
	r int
}

type Query struct {
	x   int
	y   int
	idx int
}

type SegTree struct {
	n    int
	tree []int
}

func NewSegTree(n int) *SegTree {
	size := 1
	for size < n+2 {
		size <<= 1
	}
	return &SegTree{n: size, tree: make([]int, 2*size)}
}

func (st *SegTree) Update(pos, val int) {
	pos += st.n
	if st.tree[pos] >= val {
		return
	}
	st.tree[pos] = val
	for pos > 1 {
		pos >>= 1
		if st.tree[pos] >= val {
			break
		}
		st.tree[pos] = val
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.n
	r += st.n
	res := 0
	for l <= r {
		if l&1 == 1 {
			res = max(res, st.tree[l])
			l++
		}
		if r&1 == 0 {
			res = max(res, st.tree[r])
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

	var n, m int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &m)

	ropes := make([]Rope, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &ropes[i].l, &ropes[i].r)
	}

	sort.Slice(ropes, func(i, j int) bool { return ropes[i].r < ropes[j].r })

	var q int
	fmt.Fscan(in, &q)
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i].x, &queries[i].y)
		queries[i].idx = i
	}

	sort.Slice(queries, func(i, j int) bool { return queries[i].y < queries[j].y })

	st := NewSegTree(n)
	ans := make([]int, q)
	j := 0
	for _, qu := range queries {
		for j < m && ropes[j].r <= qu.y {
			if ropes[j].l <= ropes[j].r {
				st.Update(ropes[j].l, ropes[j].r)
			}
			j++
		}
		h := qu.x
		if h > qu.y {
			ans[qu.idx] = qu.x
			continue
		}
		for {
			mx := st.Query(qu.x, h)
			if mx > h {
				h = mx
				if h >= qu.y {
					if h > qu.y {
						h = qu.y
					}
					break
				}
			} else {
				break
			}
		}
		if h > qu.y {
			h = qu.y
		}
		ans[qu.idx] = h
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
