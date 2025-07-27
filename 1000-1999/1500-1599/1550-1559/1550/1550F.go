package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int
}

type Query struct {
	i, k int
	idx  int
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return
	}
	if d.size[x] < d.size[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.size[x] += d.size[y]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func weight(x, y, d int) int {
	diff := abs(x - y)
	return abs(diff - d)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q, s, d int
	if _, err := fmt.Fscan(reader, &n, &q, &s, &d); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	s--

	edges := make([]Edge, 0, 8*n)
	for i := 0; i < n-1; i++ {
		edges = append(edges, Edge{u: i, v: i + 1, w: weight(a[i], a[i+1], d)})
	}
	for i := 0; i < n; i++ {
		x := a[i]
		j := sort.SearchInts(a, x+d)
		if j < n {
			edges = append(edges, Edge{u: i, v: j, w: weight(x, a[j], d)})
		}
		if j-1 >= 0 {
			edges = append(edges, Edge{u: i, v: j - 1, w: weight(x, a[j-1], d)})
		}
		j = sort.SearchInts(a, x-d)
		if j < n {
			edges = append(edges, Edge{u: i, v: j, w: weight(x, a[j], d)})
		}
		if j-1 >= 0 {
			edges = append(edges, Edge{u: i, v: j - 1, w: weight(x, a[j-1], d)})
		}
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

	queries := make([]Query, q)
	for idx := 0; idx < q; idx++ {
		fmt.Fscan(reader, &queries[idx].i, &queries[idx].k)
		queries[idx].i--
		queries[idx].idx = idx
	}
	sort.Slice(queries, func(i, j int) bool { return queries[i].k < queries[j].k })

	dsu := NewDSU(n)
	ans := make([]string, q)
	eidx := 0
	for _, qu := range queries {
		for eidx < len(edges) && edges[eidx].w <= qu.k {
			dsu.union(edges[eidx].u, edges[eidx].v)
			eidx++
		}
		if dsu.find(s) == dsu.find(qu.i) {
			ans[qu.idx] = "Yes"
		} else {
			ans[qu.idx] = "No"
		}
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
