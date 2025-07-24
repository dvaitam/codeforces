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
	q   int
	idx int
}

type DSU struct {
	parent []int
	size   []int64
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	size := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) int64 {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return 0
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	pairs := d.size[ra] * d.size[rb]
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return pairs
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	edges := make([]Edge, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
	}

	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &queries[i].q)
		queries[i].idx = i
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	sort.Slice(queries, func(i, j int) bool { return queries[i].q < queries[j].q })

	dsu := NewDSU(n)
	answers := make([]int64, m)
	var pairs int64
	e := 0
	for _, qu := range queries {
		for e < len(edges) && edges[e].w <= qu.q {
			pairs += dsu.Union(edges[e].u, edges[e].v)
			e++
		}
		answers[qu.idx] = pairs
	}

	for i := 0; i < m; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, answers[i])
	}
	fmt.Fprintln(out)
}
