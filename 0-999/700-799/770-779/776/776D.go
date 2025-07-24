package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	rank   []int
	diff   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		diff:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) (int, int) {
	if d.parent[x] == x {
		return x, 0
	}
	root, parity := d.find(d.parent[x])
	parity ^= d.diff[x]
	d.parent[x] = root
	d.diff[x] = parity
	return root, parity
}

func (d *DSU) unite(x, y, val int) bool {
	rx, px := d.find(x)
	ry, py := d.find(y)
	if rx == ry {
		return (px ^ py) == val
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
		px, py = py, px
	}
	d.parent[ry] = rx
	d.diff[ry] = px ^ py ^ val
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	door := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &door[i])
	}

	// each door is controlled by exactly two switches
	switches := make([][2]int, n+1)
	count := make([]int, n+1)
	for i := 1; i <= m; i++ {
		var k int
		fmt.Fscan(reader, &k)
		for j := 0; j < k; j++ {
			var r int
			fmt.Fscan(reader, &r)
			if count[r] < 2 {
				switches[r][count[r]] = i
				count[r]++
			}
		}
	}

	dsu := NewDSU(m + 1) // switches are 1..m
	possible := true
	for i := 1; i <= n && possible; i++ {
		if count[i] != 2 {
			possible = false
			break
		}
		a := switches[i][0]
		b := switches[i][1]
		val := 1 - door[i]
		if !dsu.unite(a, b, val) {
			possible = false
			break
		}
	}

	if possible {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
