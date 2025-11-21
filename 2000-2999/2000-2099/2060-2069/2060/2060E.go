package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	u, v int
}

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	parent := make([]int, n+1)
	size := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &dsu{parent: parent, size: size}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m1, m2 int
		fmt.Fscan(in, &n, &m1, &m2)

		edges := make([]edge, m1)
		for i := 0; i < m1; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v)
		}

		dsuG := newDSU(n)
		for i := 0; i < m2; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			dsuG.union(u, v)
		}

		dsuF := newDSU(n)
		removals := 0
		for _, e := range edges {
			if dsuG.find(e.u) != dsuG.find(e.v) {
				removals++
			} else {
				dsuF.union(e.u, e.v)
			}
		}

		counts := make([]int, n+1)
		visited := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			fi := dsuF.find(i)
			if !visited[fi] {
				visited[fi] = true
				gi := dsuG.find(i)
				counts[gi]++
			}
		}

		additions := 0
		for i := 1; i <= n; i++ {
			if counts[i] > 0 {
				additions += counts[i] - 1
			}
		}

		fmt.Fprintln(out, removals+additions)
	}
}
