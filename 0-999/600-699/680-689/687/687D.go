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

type DSU struct {
	p      []int
	rank   []int
	parity []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		p:      make([]int, n+1),
		rank:   make([]int, n+1),
		parity: make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	if d.p[x] != x {
		r, p := d.find(d.p[x])
		d.p[x] = r
		d.parity[x] ^= p
	}
	return d.p[x], d.parity[x]
}

func (d *DSU) union(x, y int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == 1
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.p[ry] = rx
	d.parity[ry] = px ^ py ^ 1
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}

	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		edges[i] = Edge{u, v, w}
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)

		tmp := make([]Edge, r-l+1)
		copy(tmp, edges[l-1:r])
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].w > tmp[j].w })

		dsu := NewDSU(n)
		ans := -1
		for _, e := range tmp {
			if !dsu.union(e.u, e.v) {
				ans = e.w
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
