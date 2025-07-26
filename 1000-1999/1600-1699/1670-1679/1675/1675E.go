package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{parent: p}
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
	if ra < rb {
		d.parent[rb] = ra
	} else {
		d.parent[ra] = rb
	}
}

func main() {
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		solve()
	}
}

func solve() {
	var n int
	var k int
	fmt.Fscan(reader, &n, &k)
	var s string
	fmt.Fscan(reader, &s)
	d := newDSU(26)
	bytes := []byte(s)
	for i := 0; i < n; i++ {
		idx := int(bytes[i] - 'a')
		root := d.find(idx)
		for root > 0 && k > 0 && d.find(idx) == root {
			d.union(root, root-1)
			k--
			root = d.find(idx)
		}
	}
	for i := 0; i < n; i++ {
		idx := int(bytes[i] - 'a')
		bytes[i] = byte(d.find(idx)) + 'a'
	}
	fmt.Fprintln(writer, string(bytes))
}
