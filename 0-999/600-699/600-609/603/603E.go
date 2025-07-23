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
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	for x != d.parent[x] {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(a, b int, odd *int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	*odd -= d.size[ra] % 2
	*odd -= d.size[rb] % 2
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	*odd += d.size[ra] % 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
	}

	if n%2 == 1 {
		for i := 0; i < m; i++ {
			fmt.Fprintln(out, -1)
		}
		return
	}

	for i := 1; i <= m; i++ {
		subset := make([]Edge, i)
		copy(subset, edges[:i])
		sort.Slice(subset, func(a, b int) bool { return subset[a].w < subset[b].w })
		d := NewDSU(n)
		odd := n
		ans := int64(-1)
		for _, e := range subset {
			d.Union(e.u, e.v, &odd)
			if odd == 0 {
				ans = e.w
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
