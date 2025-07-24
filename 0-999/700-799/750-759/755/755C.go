package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	fx := d.Find(x)
	fy := d.Find(y)
	if fx != fy {
		d.parent[fx] = fy
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	dsu := NewDSU(n)
	for i := 1; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		dsu.Union(i, p)
	}

	seen := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		seen[dsu.Find(i)] = struct{}{}
	}
	fmt.Fprintln(writer, len(seen))
}
