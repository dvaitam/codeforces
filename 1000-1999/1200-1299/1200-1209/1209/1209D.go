package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct{ parent []int }

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

func (d *DSU) Union(x, y int) bool {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return false
	}
	d.parent[rx] = ry
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	dsu := NewDSU(n)
	sad := 0
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if !dsu.Union(x, y) {
			sad++
		}
	}
	fmt.Fprintln(writer, sad)
}
