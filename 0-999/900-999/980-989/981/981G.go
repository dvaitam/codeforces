package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

type node struct {
	sum int64
	mul int64
	add int64
}

type segtree struct {
	n    int
	tree []node
}

func newSegTree(n int) *segtree {
	t := &segtree{n: n, tree: make([]node, 4*n)}
	for i := range t.tree {
		t.tree[i].mul = 1
	}
	return t
}

func (t *segtree) apply(idx int, mul, add int64, length int) {
	nd := &t.tree[idx]
	nd.sum = (nd.sum*mul + add*int64(length)) % mod
	nd.mul = nd.mul * mul % mod
	nd.add = (nd.add*mul + add) % mod
}

func (t *segtree) push(idx, l, r int) {
	nd := &t.tree[idx]
	if nd.mul == 1 && nd.add == 0 {
		return
	}
	mid := (l + r) / 2
	t.apply(idx*2, nd.mul, nd.add, mid-l+1)
	t.apply(idx*2+1, nd.mul, nd.add, r-mid)
	nd.mul = 1
	nd.add = 0
}

func (t *segtree) update(idx, l, r, ql, qr int, mul, add int64) {
	if ql <= l && r <= qr {
		t.apply(idx, mul, add, r-l+1)
		return
	}
	t.push(idx, l, r)
	mid := (l + r) / 2
	if ql <= mid {
		t.update(idx*2, l, mid, ql, qr, mul, add)
	}
	if qr > mid {
		t.update(idx*2+1, mid+1, r, ql, qr, mul, add)
	}
	t.tree[idx].sum = (t.tree[idx*2].sum + t.tree[idx*2+1].sum) % mod
}

func (t *segtree) UpdateMul(l, r int) {
	if l > r {
		return
	}
	t.update(1, 1, t.n, l, r, 2, 0)
}

func (t *segtree) UpdateAdd(l, r int) {
	if l > r {
		return
	}
	t.update(1, 1, t.n, l, r, 1, 1)
}

func (t *segtree) query(idx, l, r, ql, qr int) int64 {
	if ql <= l && r <= qr {
		return t.tree[idx].sum
	}
	t.push(idx, l, r)
	mid := (l + r) / 2
	var res int64
	if ql <= mid {
		res += t.query(idx*2, l, mid, ql, qr)
	}
	if qr > mid {
		res += t.query(idx*2+1, mid+1, r, ql, qr)
	}
	return res % mod
}

func (t *segtree) Query(l, r int) int64 {
	if l > r {
		return 0
	}
	return t.query(1, 1, t.n, l, r)
}

type interval struct{ l, r int }

func mergeIntervals(arr []interval) []interval {
	if len(arr) == 0 {
		return arr
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].l < arr[j].l })
	res := []interval{arr[0]}
	for _, iv := range arr[1:] {
		last := &res[len(res)-1]
		if iv.l <= last.r+1 {
			if iv.r > last.r {
				last.r = iv.r
			}
		} else {
			res = append(res, iv)
		}
	}
	return res
}

func addNumber(x, l, r int, st *segtree, sets map[int][]interval) {
	arr := sets[x]
	res := make([]interval, 0, len(arr)+1)
	i := 0
	for i < len(arr) && arr[i].r < l {
		res = append(res, arr[i])
		i++
	}
	start, end := l, r
	cur := l
	for i < len(arr) && arr[i].l <= r {
		iv := arr[i]
		if cur < iv.l {
			st.UpdateAdd(cur, min(iv.l-1, r))
		}
		interL := max(iv.l, l)
		interR := min(iv.r, r)
		if interL <= interR {
			st.UpdateMul(interL, interR)
		}
		if iv.l < start {
			start = iv.l
		}
		if iv.r > end {
			end = iv.r
		}
		cur = iv.r + 1
		i++
	}
	if cur <= r {
		st.UpdateAdd(cur, r)
	}
	res = append(res, interval{start, end})
	res = append(res, arr[i:]...)
	sets[x] = mergeIntervals(res)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	st := newSegTree(n)
	sets := make(map[int][]interval)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			addNumber(x, l, r, st, sets)
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			ans := st.Query(l, r)
			fmt.Fprintln(writer, ans)
		}
	}
}
