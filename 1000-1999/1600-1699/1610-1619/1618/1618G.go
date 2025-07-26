package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct {
	parent []int
	right  []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), right: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.right[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	d.parent[a] = b
	if d.right[a] > d.right[b] {
		d.right[b] = d.right[a]
	}
}

type item struct {
	val   int
	fromA bool
}

type edge struct {
	diff int
	idx  int
}

type query struct {
	k   int
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}
	items := make([]item, 0, n+m)
	maxPosA := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		items = append(items, item{val: x, fromA: true})
	}
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(reader, &x)
		items = append(items, item{val: x, fromA: false})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].val < items[j].val
	})
	for i := range items {
		if items[i].fromA {
			maxPosA = i
		}
	}
	N := len(items)
	pref := make([]int64, N+1)
	for i := 0; i < N; i++ {
		pref[i+1] = pref[i] + int64(items[i].val)
	}
	edges := make([]edge, 0, N-1)
	for i := 0; i < N-1; i++ {
		diff := items[i+1].val - items[i].val
		edges = append(edges, edge{diff: diff, idx: i})
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].diff < edges[j].diff
	})
	queries := make([]query, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &queries[i].k)
		queries[i].idx = i
	}
	sort.Slice(queries, func(i, j int) bool {
		return queries[i].k < queries[j].k
	})
	d := newDSU(N)
	e := 0
	ans := make([]int64, q)
	for _, qu := range queries {
		for e < len(edges) && edges[e].diff <= qu.k {
			d.union(edges[e].idx, edges[e].idx+1)
			e++
		}
		root := d.find(maxPosA)
		r := d.right[root]
		sum := pref[r+1] - pref[r+1-n]
		ans[qu.idx] = sum
	}
	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, ans[i])
	}
}
