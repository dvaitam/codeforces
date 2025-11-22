package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int64
}

type dsu struct {
	parent []int
	sz     []int
	minW   []int64
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	mw := make([]int64, n)
	const inf int64 = 1 << 60
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
		mw[i] = inf
	}
	return &dsu{parent: p, sz: s, minW: mw}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) unite(a, b int, w int64) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		// Even if already connected, we might need to update minW if it is still INF.
		if w < d.minW[ra] {
			d.minW[ra] = w
		}
		return
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.sz[ra] += d.sz[rb]
	// minimum edge weight available inside the merged component
	if d.minW[rb] < d.minW[ra] {
		d.minW[ra] = d.minW[rb]
	}
	if w < d.minW[ra] {
		d.minW[ra] = w
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
			edges[i].u--
			edges[i].v--
		}

		sort.Slice(edges, func(i, j int) bool {
			return edges[i].w < edges[j].w
		})

		d := newDSU(n)
		var ans int64 = -1
		for _, e := range edges {
			d.unite(e.u, e.v, e.w)
			if d.find(0) == d.find(n-1) {
				root := d.find(0)
				ans = e.w + d.minW[root]
				break // further edges only increase max weight
			}
		}

		fmt.Fprintln(out, ans)
	}
}
