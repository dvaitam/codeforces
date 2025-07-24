package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
	mx     []int
	need   []bool
	hasOne []bool
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		size:   make([]int, n),
		mx:     make([]int, n),
		need:   make([]bool, n),
		hasOne: make([]bool, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		d.mx[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) int {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return ra
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	if d.mx[rb] > d.mx[ra] {
		d.mx[ra] = d.mx[rb]
	}
	d.need[ra] = d.need[ra] || d.need[rb]
	d.hasOne[ra] = d.hasOne[ra] || d.hasOne[rb]
	return ra
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	type Edge struct{ u, v, x int }
	edges := make([]Edge, q)
	for i := 0; i < q; i++ {
		var u, v, x int
		fmt.Fscan(reader, &u, &v, &x)
		u--
		v--
		edges[i] = Edge{u, v, x}
	}

	ans := make([]int, n)
	for bit := 0; bit < 30; bit++ {
		forced0 := make([]bool, n)
		for _, e := range edges {
			if ((e.x >> bit) & 1) == 0 {
				forced0[e.u] = true
				forced0[e.v] = true
			}
		}
		dsu := NewDSU(n)
		forced1 := make([]bool, n)
		for _, e := range edges {
			if ((e.x >> bit) & 1) == 1 {
				fu := forced0[e.u]
				fv := forced0[e.v]
				if fu && fv {
					// contradiction impossible
				} else if fu && !fv {
					forced1[e.v] = true
				} else if fv && !fu {
					forced1[e.u] = true
				} else {
					root := dsu.union(e.u, e.v)
					dsu.need[root] = true
				}
			}
		}
		for i := 0; i < n; i++ {
			if forced0[i] {
				continue
			}
			r := dsu.find(i)
			if dsu.mx[r] < i {
				dsu.mx[r] = i
			}
		}
		for i := 0; i < n; i++ {
			if forced1[i] {
				r := dsu.find(i)
				dsu.hasOne[r] = true
			}
		}
		seen := make(map[int]bool)
		for i := 0; i < n; i++ {
			if forced0[i] {
				continue
			}
			r := dsu.find(i)
			if seen[r] {
				continue
			}
			seen[r] = true
			if dsu.need[r] && !dsu.hasOne[r] {
				idx := dsu.mx[r]
				forced1[idx] = true
				dsu.hasOne[r] = true
			}
		}
		for i := 0; i < n; i++ {
			if forced1[i] {
				ans[i] |= (1 << bit)
			}
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, ans[i])
	}
	writer.WriteByte('\n')
}
