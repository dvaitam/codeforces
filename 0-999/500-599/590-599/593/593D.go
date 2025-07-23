package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1000000000000000001

type Edge struct {
	to  int
	idx int
}

type SegTree struct {
	n   int
	val []int64
}

func mul(a, b int64) int64 {
	if a >= INF || b >= INF {
		return INF
	}
	if a > INF/b {
		return INF
	}
	return a * b
}

func NewSegTree(arr []int64) *SegTree {
	st := &SegTree{n: len(arr), val: make([]int64, 4*len(arr))}
	var build func(p, l, r int)
	build = func(p, l, r int) {
		if l == r {
			st.val[p] = arr[l]
			return
		}
		m := (l + r) >> 1
		build(p<<1, l, m)
		build(p<<1|1, m+1, r)
		st.val[p] = mul(st.val[p<<1], st.val[p<<1|1])
	}
	build(1, 0, st.n-1)
	return st
}

func (st *SegTree) Update(idx int, v int64) {
	var rec func(p, l, r int)
	rec = func(p, l, r int) {
		if l == r {
			st.val[p] = v
			return
		}
		m := (l + r) >> 1
		if idx <= m {
			rec(p<<1, l, m)
		} else {
			rec(p<<1|1, m+1, r)
		}
		st.val[p] = mul(st.val[p<<1], st.val[p<<1|1])
	}
	rec(1, 0, st.n-1)
}

func (st *SegTree) Query(l, r int) int64 {
	if l > r {
		return 1
	}
	var rec func(p, L, R int) int64
	rec = func(p, L, R int) int64 {
		if l <= L && R <= r {
			return st.val[p]
		}
		if R < l || L > r {
			return 1
		}
		m := (L + R) >> 1
		left := rec(p<<1, L, m)
		right := rec(p<<1|1, m+1, R)
		return mul(left, right)
	}
	return rec(1, 0, st.n-1)
}

var (
	n, m     int
	edgesU   []int
	edgesV   []int
	edgesW   []int64
	g        [][]Edge
	parent   []int
	depth    []int
	heavy    []int
	sizeArr  []int
	head     []int
	pos      []int
	value    []int64
	cur      int
	seg      *SegTree
	edgeNode []int
)

func dfs(u, p int) {
	sizeArr[u] = 1
	heavy[u] = 0
	for _, e := range g[u] {
		v := e.to
		if v == p {
			continue
		}
		parent[v] = u
		depth[v] = depth[u] + 1
		value[v] = edgesW[e.idx]
		edgeNode[e.idx] = v
		dfs(v, u)
		sizeArr[u] += sizeArr[v]
		if heavy[u] == 0 || sizeArr[v] > sizeArr[heavy[u]] {
			heavy[u] = v
		}
	}
}

func decompose(u, h int) {
	head[u] = h
	pos[u] = cur
	cur++
	if heavy[u] != 0 {
		decompose(heavy[u], h)
	}
	for _, e := range g[u] {
		v := e.to
		if v == parent[u] || v == heavy[u] {
			continue
		}
		decompose(v, v)
	}
}

func queryPath(u, v int, limit int64) int64 {
	res := int64(1)
	for head[u] != head[v] {
		if depth[head[u]] > depth[head[v]] {
			res = mul(res, seg.Query(pos[head[u]], pos[u]))
			if res > limit {
				return res
			}
			u = parent[head[u]]
		} else {
			res = mul(res, seg.Query(pos[head[v]], pos[v]))
			if res > limit {
				return res
			}
			v = parent[head[v]]
		}
	}
	if depth[u] > depth[v] {
		u, v = v, u
	}
	if u != v {
		res = mul(res, seg.Query(pos[u]+1, pos[v]))
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edgesU = make([]int, n)
	edgesV = make([]int, n)
	edgesW = make([]int64, n)
	g = make([][]Edge, n+1)
	for i := 1; i < n; i++ {
		var u, v int
		var w int64
		fmt.Fscan(reader, &u, &v, &w)
		edgesU[i] = u
		edgesV[i] = v
		edgesW[i] = w
		g[u] = append(g[u], Edge{to: v, idx: i})
		g[v] = append(g[v], Edge{to: u, idx: i})
	}
	parent = make([]int, n+1)
	depth = make([]int, n+1)
	heavy = make([]int, n+1)
	sizeArr = make([]int, n+1)
	head = make([]int, n+1)
	pos = make([]int, n+1)
	value = make([]int64, n+1)
	edgeNode = make([]int, n)

	value[1] = 1
	dfs(1, 0)
	cur = 0
	decompose(1, 1)

	arr := make([]int64, n)
	for i := 1; i <= n; i++ {
		arr[pos[i]] = value[i]
	}
	seg = NewSegTree(arr)

	for i := 0; i < m; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var a, b int
			var y int64
			fmt.Fscan(reader, &a, &b, &y)
			prod := queryPath(a, b, y)
			if prod > y {
				fmt.Fprintln(writer, 0)
			} else {
				fmt.Fprintln(writer, y/prod)
			}
		} else {
			var p int
			var c int64
			fmt.Fscan(reader, &p, &c)
			edgesW[p] = c
			node := edgeNode[p]
			seg.Update(pos[node], c)
		}
	}
}
