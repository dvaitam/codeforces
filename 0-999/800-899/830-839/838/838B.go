package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	from int
	to   int
	w    int64
}

// Fenwick tree supporting range add and point query
type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) add(idx int, val int64) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *Fenwick) rangeAdd(l, r int, val int64) {
	if l > r {
		return
	}
	f.add(l, val)
	if r+1 <= f.n {
		f.add(r+1, -val)
	}
}

func (f *Fenwick) prefixSum(idx int) int64 {
	res := int64(0)
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

// Segment tree for range add and range min query
type SegTree struct {
	n    int
	tree []int64
	lazy []int64
}

func NewSegTree(arr []int64) *SegTree {
	n := len(arr) - 1
	st := &SegTree{n: n, tree: make([]int64, 4*(n+2)), lazy: make([]int64, 4*(n+2))}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(node, l, r int, arr []int64) {
	if l == r {
		st.tree[node] = arr[l]
		return
	}
	m := (l + r) / 2
	st.build(node*2, l, m, arr)
	st.build(node*2+1, m+1, r, arr)
	if st.tree[node*2] < st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *SegTree) push(node int) {
	if st.lazy[node] != 0 {
		val := st.lazy[node]
		st.tree[node*2] += val
		st.lazy[node*2] += val
		st.tree[node*2+1] += val
		st.lazy[node*2+1] += val
		st.lazy[node] = 0
	}
}

func (st *SegTree) rangeAdd(node, l, r, L, R int, val int64) {
	if L > r || R < l {
		return
	}
	if L <= l && r <= R {
		st.tree[node] += val
		st.lazy[node] += val
		return
	}
	st.push(node)
	m := (l + r) / 2
	st.rangeAdd(node*2, l, m, L, R, val)
	st.rangeAdd(node*2+1, m+1, r, L, R, val)
	if st.tree[node*2] < st.tree[node*2+1] {
		st.tree[node] = st.tree[node*2]
	} else {
		st.tree[node] = st.tree[node*2+1]
	}
}

func (st *SegTree) queryMin(node, l, r, L, R int) int64 {
	if L > r || R < l {
		return 1<<63 - 1
	}
	if L <= l && r <= R {
		return st.tree[node]
	}
	st.push(node)
	m := (l + r) / 2
	left := st.queryMin(node*2, l, m, L, R)
	right := st.queryMin(node*2+1, m+1, r, L, R)
	if left < right {
		return left
	}
	return right
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

var (
	n, q    int
	edges   []Edge
	child   [][]int
	starIdx []int
	wStar   []int64
	in      []int
	out     []int
	base    []int64
	timer   int
	fw      *Fenwick
	st      *SegTree
)

func dfs(v int) {
	timer++
	in[v] = timer
	for _, eIdx := range child[v] {
		e := edges[eIdx]
		base[e.to] = base[v] + e.w
		dfs(e.to)
	}
	out[v] = timer
}

func main() {
	inReader := bufio.NewReader(os.Stdin)
	outWriter := bufio.NewWriter(os.Stdout)
	defer outWriter.Flush()

	if _, err := fmt.Fscan(inReader, &n, &q); err != nil {
		return
	}
	m := 2*n - 2
	edges = make([]Edge, m+1)
	child = make([][]int, n+1)
	starIdx = make([]int, n+1)
	wStar = make([]int64, n+1)

	for i := 1; i <= m; i++ {
		var a, b int
		var c int64
		fmt.Fscan(inReader, &a, &b, &c)
		edges[i] = Edge{from: a, to: b, w: c}
		if i < n {
			child[a] = append(child[a], i)
		} else {
			starIdx[a] = i
			wStar[a] = c
		}
	}

	in = make([]int, n+1)
	out = make([]int, n+1)
	base = make([]int64, n+1)
	timer = 0
	dfs(1)

	arr := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		arr[in[i]] = base[i] + wStar[i]
	}

	fw = NewFenwick(n + 2)
	st = NewSegTree(arr)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(inReader, &t)
		if t == 1 {
			var idx int
			var w int64
			fmt.Fscan(inReader, &idx, &w)
			diff := w - edges[idx].w
			edges[idx].w = w
			if idx < n {
				c := edges[idx].to
				fw.rangeAdd(in[c], out[c], diff)
				st.rangeAdd(1, 1, n, in[c], out[c], diff)
			} else {
				x := edges[idx].from
				wStar[x] = w
				st.rangeAdd(1, 1, n, in[x], in[x], diff)
			}
		} else {
			var u, v int
			fmt.Fscan(inReader, &u, &v)
			du := base[u] + fw.prefixSum(in[u])
			dv := base[v] + fw.prefixSum(in[v])
			ans := int64(1<<63 - 1)
			if in[u] <= in[v] && out[v] <= out[u] {
				ans = dv - du
			}
			minVal := st.queryMin(1, 1, n, in[u], out[u])
			alt := dv + minVal - du
			if alt < ans {
				ans = alt
			}
			fmt.Fprintln(outWriter, ans)
		}
	}
}
