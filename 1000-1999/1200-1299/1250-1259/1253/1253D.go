package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	right  []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		r[i] = i
	}
	return &DSU{parent: p, right: r}
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
	if d.right[ra] < d.right[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.right[rb] > d.right[ra] {
		d.right[ra] = d.right[rb]
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	dsu := NewDSU(n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		dsu.Union(u, v)
	}

	ans := 0
	for i := 1; i <= n; {
		r := dsu.right[dsu.Find(i)]
		j := i + 1
		for j <= r {
			if dsu.Find(j) != dsu.Find(i) {
				dsu.Union(i, j)
				ans++
			}
			r = max(r, dsu.right[dsu.Find(i)])
			j++
		}
		i = r + 1
	}

	fmt.Fprintln(writer, ans)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
