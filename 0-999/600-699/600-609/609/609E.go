package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int64
	idx  int
}

type DSU struct {
	p  []int
	sz []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p: p, sz: sz}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) bool {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return false
	}
	if d.sz[x] < d.sz[y] {
		x, y = y, x
	}
	d.p[y] = x
	d.sz[x] += d.sz[y]
	return true
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

type Node struct {
	to int
	w  int64
}

const LOG = 20

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
		edges[i].idx = i
	}

	sorted := append([]Edge(nil), edges...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].w < sorted[j].w })
	dsu := NewDSU(n)
	mark := make([]bool, m)
	adj := make([][]Node, n+1)
	mstWeight := int64(0)
	for _, e := range sorted {
		if dsu.Union(e.u, e.v) {
			mark[e.idx] = true
			mstWeight += e.w
			adj[e.u] = append(adj[e.u], Node{to: e.v, w: e.w})
			adj[e.v] = append(adj[e.v], Node{to: e.u, w: e.w})
		}
	}

	depth := make([]int, n+1)
	parent := make([][]int, LOG)
	maxUp := make([][]int64, LOG)
	for k := 0; k < LOG; k++ {
		parent[k] = make([]int, n+1)
		maxUp[k] = make([]int64, n+1)
	}

	queue := make([]int, 0, n)
	queue = append(queue, 1)
	parent[0][1] = 0
	depth[1] = 0
	for qi := 0; qi < len(queue); qi++ {
		v := queue[qi]
		for _, nb := range adj[v] {
			if nb.to == parent[0][v] {
				continue
			}
			parent[0][nb.to] = v
			maxUp[0][nb.to] = nb.w
			depth[nb.to] = depth[v] + 1
			queue = append(queue, nb.to)
		}
	}

	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			p := parent[k-1][v]
			parent[k][v] = parent[k-1][p]
			maxUp[k][v] = max64(maxUp[k-1][v], maxUp[k-1][p])
		}
	}

	getMax := func(u, v int) int64 {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		res := int64(0)
		diff := depth[u] - depth[v]
		for k := LOG - 1; k >= 0; k-- {
			if diff>>k&1 == 1 {
				if maxUp[k][u] > res {
					res = maxUp[k][u]
				}
				u = parent[k][u]
			}
		}
		if u == v {
			return res
		}
		for k := LOG - 1; k >= 0; k-- {
			if parent[k][u] != parent[k][v] {
				if maxUp[k][u] > res {
					res = maxUp[k][u]
				}
				if maxUp[k][v] > res {
					res = maxUp[k][v]
				}
				u = parent[k][u]
				v = parent[k][v]
			}
		}
		if maxUp[0][u] > res {
			res = maxUp[0][u]
		}
		if maxUp[0][v] > res {
			res = maxUp[0][v]
		}
		return res
	}

	ans := make([]int64, m)
	for _, e := range edges {
		if mark[e.idx] {
			ans[e.idx] = mstWeight
		} else {
			maxW := getMax(e.u, e.v)
			ans[e.idx] = mstWeight + e.w - maxW
		}
	}

	for i := 0; i < m; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
