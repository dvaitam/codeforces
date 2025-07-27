package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU structure for disjoint set union operations
type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
		r[i] = 0
	}
	return &DSU{parent: p, rank: r}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.rank[ra] < d.rank[rb] {
		d.parent[ra] = rb
	} else if d.rank[ra] > d.rank[rb] {
		d.parent[rb] = ra
	} else {
		d.parent[rb] = ra
		d.rank[ra]++
	}
}

// Edge represents a built trail
type Edge struct {
	u, v int
	l, r int // sorted endpoints l < r
}

func cross(e1, e2 Edge) bool {
	if e1.u == e2.u || e1.u == e2.v || e1.v == e2.u || e1.v == e2.v {
		return false
	}
	l1, r1 := e1.l, e1.r
	l2, r2 := e2.l, e2.r
	return (l1 < l2 && l2 < r1 && r1 < r2) || (l2 < l1 && l1 < r2 && r2 < r1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	dsu := NewDSU(n)
	edges := make([]Edge, 0)
	var result []byte

	for i := 0; i < m; i++ {
		var e, v, u int
		fmt.Fscan(in, &e, &v, &u)
		if e == 1 {
			if v == u {
				continue
			}
			l, r := v, u
			if l > r {
				l, r = r, l
			}
			newEdge := Edge{u: v, v: u, l: l, r: r}
			dsu.Union(v, u)
			for _, ex := range edges {
				if cross(newEdge, ex) {
					dsu.Union(newEdge.u, ex.u)
				}
			}
			edges = append(edges, newEdge)
		} else if e == 2 {
			if dsu.Find(v) == dsu.Find(u) {
				result = append(result, '1')
			} else {
				result = append(result, '0')
			}
		}
	}

	fmt.Fprintln(out, string(result))
}
