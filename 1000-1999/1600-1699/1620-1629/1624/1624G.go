package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v int
	w    int
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
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
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func (d *dsu) connected(n int) bool {
	root := d.find(1)
	for i := 2; i <= n; i++ {
		if d.find(i) != root {
			return false
		}
	}
	return true
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
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
		}

		cur := edges
		ans := 0
		for bit := 30; bit >= 0; bit-- {
			d := newDSU(n)
			nextEdges := make([]edge, 0, len(cur))
			mask := 1 << uint(bit)
			for _, e := range cur {
				if e.w&mask == 0 {
					d.union(e.u, e.v)
					nextEdges = append(nextEdges, e)
				}
			}
			if d.connected(n) {
				cur = nextEdges
			} else {
				ans |= mask
			}
		}
		fmt.Fprintln(out, ans)
	}
}
