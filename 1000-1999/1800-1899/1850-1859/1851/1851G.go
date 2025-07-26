package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func (d *DSU) same(a, b int) bool {
	return d.find(a) == d.find(b)
}

type Edge struct {
	u, v int
	h    int
}

type Query struct {
	a, b int
	t    int
	idx  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		h := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &h[i])
		}
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			mh := h[u]
			if h[v] > mh {
				mh = h[v]
			}
			edges[i] = Edge{u, v, mh}
		}
		sort.Slice(edges, func(i, j int) bool { return edges[i].h < edges[j].h })

		var q int
		fmt.Fscan(reader, &q)
		queries := make([]Query, q)
		for i := 0; i < q; i++ {
			var a, b, e int
			fmt.Fscan(reader, &a, &b, &e)
			t := h[a] + e
			queries[i] = Query{a: a, b: b, t: t, idx: i}
		}
		sort.Slice(queries, func(i, j int) bool { return queries[i].t < queries[j].t })

		ans := make([]bool, q)
		dsu := NewDSU(n)
		ei := 0
		for _, qu := range queries {
			for ei < m && edges[ei].h <= qu.t {
				dsu.union(edges[ei].u, edges[ei].v)
				ei++
			}
			if h[qu.b] <= qu.t && dsu.same(qu.a, qu.b) {
				ans[qu.idx] = true
			} else {
				ans[qu.idx] = false
			}
		}
		for i := 0; i < q; i++ {
			if ans[i] {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
