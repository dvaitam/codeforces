package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

type PQItem struct {
	wDiff int64
	id    int
}

// Priority queue for min weight
type MinHeap []PQItem

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].wDiff < h[j].wDiff }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(PQItem)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// Fenwick tree supporting range add and point query
type Fenwick struct {
	n   int
	bit []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) add(i int, v int64) {
	for i <= f.n {
		f.bit[i] += v
		i += i & -i
	}
}

func (f *Fenwick) RangeAdd(l, r int, v int64) {
	if l > r {
		return
	}
	f.add(l, v)
	f.add(r+1, -v)
}

func (f *Fenwick) Query(i int) int64 {
	res := int64(0)
	for i > 0 {
		res += f.bit[i]
		i -= i & -i
	}
	return res
}

// Segment tree for range add and range min with index
type SegNode struct {
	val  int64
	idx  int
	lazy int64
}

type SegTree struct {
	n    int
	tree []SegNode
}

func NewSegTree(arr []int64, id []int) *SegTree {
	n := len(arr) - 1
	st := &SegTree{n: n, tree: make([]SegNode, 4*(n+2))}
	var build func(p, l, r int)
	build = func(p, l, r int) {
		if l == r {
			st.tree[p] = SegNode{val: arr[l], idx: id[l]}
			return
		}
		mid := (l + r) >> 1
		build(p<<1, l, mid)
		build(p<<1|1, mid+1, r)
		st.tree[p] = st.combine(st.tree[p<<1], st.tree[p<<1|1])
	}
	build(1, 1, n)
	return st
}

func (st *SegTree) combine(a, b SegNode) SegNode {
	if a.val < b.val {
		return SegNode{val: a.val, idx: a.idx}
	} else if b.val < a.val {
		return SegNode{val: b.val, idx: b.idx}
	}
	if a.idx < b.idx {
		return SegNode{val: a.val, idx: a.idx}
	}
	return SegNode{val: b.val, idx: b.idx}
}

func (st *SegTree) apply(p int, v int64) {
	st.tree[p].val += v
	st.tree[p].lazy += v
}

func (st *SegTree) push(p int) {
	if st.tree[p].lazy != 0 {
		v := st.tree[p].lazy
		st.apply(p<<1, v)
		st.apply(p<<1|1, v)
		st.tree[p].lazy = 0
	}
}

func (st *SegTree) UpdateRange(p, l, r, ql, qr int, v int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		st.apply(p, v)
		return
	}
	st.push(p)
	mid := (l + r) >> 1
	st.UpdateRange(p<<1, l, mid, ql, qr, v)
	st.UpdateRange(p<<1|1, mid+1, r, ql, qr, v)
	st.tree[p] = st.combine(st.tree[p<<1], st.tree[p<<1|1])
}

func (st *SegTree) UpdatePoint(p, l, r, idx int, val int64) {
	if l == r {
		st.tree[p].val = val
		return
	}
	st.push(p)
	mid := (l + r) >> 1
	if idx <= mid {
		st.UpdatePoint(p<<1, l, mid, idx, val)
	} else {
		st.UpdatePoint(p<<1|1, mid+1, r, idx, val)
	}
	st.tree[p] = st.combine(st.tree[p<<1], st.tree[p<<1|1])
}

func (st *SegTree) Query(p, l, r, ql, qr int) SegNode {
	if ql > r || qr < l {
		return SegNode{val: INF, idx: 0}
	}
	if ql <= l && r <= qr {
		return st.tree[p]
	}
	st.push(p)
	mid := (l + r) >> 1
	left := st.Query(p<<1, l, mid, ql, qr)
	right := st.Query(p<<1|1, mid+1, r, ql, qr)
	return st.combine(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	// girls living junctions
	girlsAt := make([]int, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(reader, &girlsAt[i])
	}

	// prepare HLD
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	size := make([]int, n+1)
	heavy := make([]int, n+1)
	// iterative dfs1
	type stackEntry struct{ u, idx int }
	stack := []stackEntry{{1, 0}}
	parent[1] = 0
	depth[1] = 0
	order := make([]int, 0, n)
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		u := top.u
		if top.idx == 0 {
			order = append(order, u)
		}
		if top.idx < len(adj[u]) {
			v := adj[u][top.idx]
			top.idx++
			if v != parent[u] {
				parent[v] = u
				depth[v] = depth[u] + 1
				stack = append(stack, stackEntry{v, 0})
			}
		} else {
			size[u] = 1
			maxSize := 0
			heavy[u] = 0
			for _, v := range adj[u] {
				if v != parent[u] {
					size[u] += size[v]
					if size[v] > maxSize {
						maxSize = size[v]
						heavy[u] = v
					}
				}
			}
			stack = stack[:len(stack)-1]
		}
	}

	head := make([]int, n+1)
	pos := make([]int, n+1)
	rev := make([]int, n+1)
	curPos := 1
	type hldEntry struct{ u, h int }
	hstack := []hldEntry{{1, 1}}
	for len(hstack) > 0 {
		ent := hstack[len(hstack)-1]
		hstack = hstack[:len(hstack)-1]
		u, h := ent.u, ent.h
		head[u] = h
		pos[u] = curPos
		rev[curPos] = u
		curPos++
		// push light children first
		for i := len(adj[u]) - 1; i >= 0; i-- {
			v := adj[u][i]
			if v != parent[u] && v != heavy[u] {
				hstack = append(hstack, hldEntry{v, v})
			}
		}
		if heavy[u] != 0 {
			hstack = append(hstack, hldEntry{heavy[u], h})
		}
	}

	// fenwick for node additions
	fw := NewFenwick(n + 2)
	// heaps per node
	heaps := make([]MinHeap, n+1)
	for i := 1; i <= n; i++ {
		heaps[i] = MinHeap{}
	}
	for id := 1; id <= m; id++ {
		u := girlsAt[id]
		heap.Push(&heaps[u], PQItem{wDiff: int64(id), id: id})
	}
	// initial array for segment tree
	arr := make([]int64, n+1)
	ids := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if len(heaps[i]) > 0 {
			arr[pos[i]] = heaps[i][0].wDiff
		} else {
			arr[pos[i]] = INF
		}
		ids[pos[i]] = i
	}
	st := NewSegTree(arr, ids)

	// helper to get subtree range
	tout := make([]int, n+1)
	for i := 1; i <= n; i++ {
		tout[i] = pos[i] + size[i] - 1
	}

	better := func(a, b SegNode) SegNode {
		if a.val < b.val {
			return a
		} else if b.val < a.val {
			return b
		}
		if a.idx < b.idx {
			return a
		}
		return b
	}

	queryPath := func(u, v int) SegNode {
		res := SegNode{val: INF, idx: 0}
		for head[u] != head[v] {
			if depth[head[u]] < depth[head[v]] {
				u, v = v, u
			}
			cur := st.Query(1, 1, n, pos[head[u]], pos[u])
			res = better(res, cur)
			u = parent[head[u]]
		}
		if depth[u] > depth[v] {
			u, v = v, u
		}
		cur := st.Query(1, 1, n, pos[u], pos[v])
		res = better(res, cur)
		return res
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var v, u, k int
			fmt.Fscan(reader, &v, &u, &k)
			ans := make([]int, 0)
			for k > 0 {
				best := queryPath(v, u)
				if best.val >= INF {
					break
				}
				node := best.idx
				nodeAdd := fw.Query(pos[node])
				item := heap.Pop(&heaps[node]).(PQItem)
				ans = append(ans, item.id)
				if len(heaps[node]) > 0 {
					newVal := heaps[node][0].wDiff + nodeAdd
					st.UpdatePoint(1, 1, n, pos[node], newVal)
				} else {
					st.UpdatePoint(1, 1, n, pos[node], INF)
				}
				k--
			}
			fmt.Fprint(writer, len(ans))
			for _, id := range ans {
				fmt.Fprint(writer, " ", id)
			}
			fmt.Fprintln(writer)
		} else {
			var v int
			var add int64
			fmt.Fscan(reader, &v, &add)
			fw.RangeAdd(pos[v], tout[v], add)
			st.UpdateRange(1, 1, n, pos[v], tout[v], add)
		}
	}
}
