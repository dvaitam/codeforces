package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct{ u, v int }

type DSU struct {
	parent []int
	parity []int
	size   []int
	good   bool
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), parity: make([]int, n+1), size: make([]int, n+1), good: true}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	p := 0
	for x != d.parent[x] {
		p ^= d.parity[x]
		x = d.parent[x]
	}
	return x, p
}

func (d *DSU) Union(a, b int) {
	ra, pa := d.find(a)
	rb, pb := d.find(b)
	if ra == rb {
		if pa == pb {
			d.good = false
		}
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
		pa, pb = pb, pa
	}
	d.parent[rb] = ra
	d.parity[rb] = pa ^ pb ^ 1
	d.size[ra] += d.size[rb]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(in, &n, &m, &q); err != nil {
		return
	}
	edges := make([]Edge, m+1)
	for i := 1; i <= m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v)
	}
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		d := NewDSU(n)
		for i := 1; i < l; i++ {
			d.Union(edges[i].u, edges[i].v)
			if !d.good {
				break
			}
		}
		if d.good {
			for i := r + 1; i <= m; i++ {
				d.Union(edges[i].u, edges[i].v)
				if !d.good {
					break
				}
			}
		}
		if d.good {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
