package main

import (
	"bufio"
	"fmt"
	"os"
)

type dsu struct {
	p  []int
	sz []int
}

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n), sz: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.sz[a] < d.sz[b] {
		a, b = b, a
	}
	d.p[b] = a
	d.sz[a] += d.sz[b]
}

type Edge struct {
	u, v int
	w    int
}

func can(edges []Edge, n int, x int) bool {
	d := newDSU(n)
	for _, e := range edges {
		if e.w <= x {
			d.union(e.u, e.v)
		}
	}
	compID := make(map[int]int)
	id := 0
	for i := 0; i < n; i++ {
		r := d.find(i)
		if _, ok := compID[r]; !ok {
			compID[r] = id
			id++
		}
	}
	indeg := make([]int, id)
	for _, e := range edges {
		if e.w > x {
			a := compID[d.find(e.u)]
			b := compID[d.find(e.v)]
			if a != b {
				indeg[b]++
			}
		}
	}
	cnt := 0
	for i := 0; i < id; i++ {
		if indeg[i] == 0 {
			cnt++
			if cnt > 1 {
				return false
			}
		}
	}
	return cnt == 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
			edges[i].u--
			edges[i].v--
		}
		// check connectivity ignoring directions
		d := newDSU(n)
		for _, e := range edges {
			d.union(e.u, e.v)
		}
		root := d.find(0)
		connected := true
		for i := 1; i < n; i++ {
			if d.find(i) != root {
				connected = false
				break
			}
		}
		if !connected {
			fmt.Fprintln(out, -1)
			continue
		}
		lo, hi := 0, int(1e9)
		for lo < hi {
			mid := (lo + hi) / 2
			if can(edges, n, mid) {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(out, lo)
	}
}
