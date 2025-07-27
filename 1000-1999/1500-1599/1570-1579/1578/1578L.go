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
}

type DSU struct {
	p  []int
	sz []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n+1), sz: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.p[i] = i
		d.sz[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.sz[ra] += d.sz[rb]
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	// Read candies, but we don't need them for computation
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
	}

	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].w > edges[j].w })

	dsu := NewDSU(n)
	comps := n
	ans := int64(-1)
	for _, e := range edges {
		if dsu.union(e.u, e.v) {
			comps--
			if comps == 1 {
				ans = e.w
				break
			}
		}
	}

	fmt.Println(ans)
}
