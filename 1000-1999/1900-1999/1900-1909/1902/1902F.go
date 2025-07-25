package main

import (
	"bufio"
	"fmt"
	"os"
)

// Basis implements xor linear basis for numbers up to 20 bits
// (values up to 2^20-1).
type Basis struct {
	b [20]int
}

// Add inserts x into the basis.
func (bs *Basis) Add(x int) {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] != 0 {
			x ^= bs.b[i]
		} else {
			bs.b[i] = x
			return
		}
	}
}

// Merge merges another basis into this one.
func (bs *Basis) Merge(o *Basis) {
	for i := 19; i >= 0; i-- {
		if o.b[i] != 0 {
			bs.Add(o.b[i])
		}
	}
}

// Contains checks whether x can be represented as xor of basis elements.
func (bs *Basis) Contains(x int) bool {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			return false
		}
		x ^= bs.b[i]
	}
	return true
}

type SegTree struct {
	n    int
	tree []Basis
}

func NewSegTree(arr []int) *SegTree {
	n := len(arr) - 1
	st := &SegTree{n: n, tree: make([]Basis, 4*(n+2))}
	st.build(1, 1, n, arr)
	return st
}

func (st *SegTree) build(p, l, r int, arr []int) {
	if l == r {
		if arr[l] != 0 {
			st.tree[p].Add(arr[l])
		}
		return
	}
	mid := (l + r) >> 1
	st.build(p<<1, l, mid, arr)
	st.build(p<<1|1, mid+1, r, arr)
	st.tree[p] = Basis{}
	st.tree[p].Merge(&st.tree[p<<1])
	st.tree[p].Merge(&st.tree[p<<1|1])
}

func (st *SegTree) query(p, l, r, ql, qr int, res *Basis) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		res.Merge(&st.tree[p])
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		st.query(p<<1, l, mid, ql, qr, res)
	}
	if qr > mid {
		st.query(p<<1|1, mid+1, r, ql, qr, res)
	}
}

type HLD struct {
	n      int
	adj    [][]int
	parent []int
	depth  []int
	heavy  []int
	head   []int
	pos    []int
	cur    int
}

func NewHLD(n int, adj [][]int) *HLD {
	h := &HLD{n: n, adj: adj}
	h.parent = make([]int, n+1)
	h.depth = make([]int, n+1)
	h.heavy = make([]int, n+1)
	h.head = make([]int, n+1)
	h.pos = make([]int, n+1)
	h.dfs()
	h.decompose()
	return h
}

func (h *HLD) dfs() {
	order := make([]int, 0, h.n)
	stack := []int{1}
	h.parent[1] = 0
	h.depth[1] = 0
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, u)
		for _, v := range h.adj[u] {
			if v != h.parent[u] {
				h.parent[v] = u
				h.depth[v] = h.depth[u] + 1
				stack = append(stack, v)
			}
		}
	}
	size := make([]int, h.n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		size[u] = 1
		maxSize := 0
		heavyChild := 0
		for _, v := range h.adj[u] {
			if v != h.parent[u] {
				size[u] += size[v]
				if size[v] > maxSize {
					maxSize = size[v]
					heavyChild = v
				}
			}
		}
		h.heavy[u] = heavyChild
	}
}

func (h *HLD) decompose() {
	h.cur = 1
	type pair struct{ u, head int }
	stack := []pair{{1, 1}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		u, hd := p.u, p.head
		for {
			h.head[u] = hd
			h.pos[u] = h.cur
			h.cur++
			for i := len(h.adj[u]) - 1; i >= 0; i-- {
				v := h.adj[u][i]
				if v != h.parent[u] && v != h.heavy[u] {
					stack = append(stack, pair{v, v})
				}
			}
			if h.heavy[u] == 0 {
				break
			}
			u = h.heavy[u]
		}
	}
}

func (h *HLD) queryPath(u, v int, st *SegTree) Basis {
	var res Basis
	for h.head[u] != h.head[v] {
		if h.depth[h.head[u]] < h.depth[h.head[v]] {
			u, v = v, u
		}
		st.query(1, 1, st.n, h.pos[h.head[u]], h.pos[u], &res)
		u = h.parent[h.head[u]]
	}
	if h.depth[u] > h.depth[v] {
		u, v = v, u
	}
	st.query(1, 1, st.n, h.pos[u], h.pos[v], &res)
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	values := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &values[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	h := NewHLD(n, adj)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[h.pos[i]] = values[i]
	}
	st := NewSegTree(arr)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var x, y, k int
		fmt.Fscan(reader, &x, &y, &k)
		bs := h.queryPath(x, y, st)
		if bs.Contains(k) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
