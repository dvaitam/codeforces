package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	p []int
}

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra != rb {
		d.p[ra] = rb
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, q int
	if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
		return
	}

	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}

	index := make(map[string]int, n)
	for i, w := range words {
		index[w] = i
	}

	dsu := NewDSU(2 * n)

	for i := 0; i < m; i++ {
		var t int
		var a, b string
		fmt.Fscan(reader, &t, &a, &b)
		x := index[a]
		y := index[b]
		if t == 1 {
			if dsu.Find(x) == dsu.Find(y+n) || dsu.Find(x+n) == dsu.Find(y) {
				fmt.Fprintln(writer, "NO")
			} else {
				dsu.Union(x, y)
				dsu.Union(x+n, y+n)
				fmt.Fprintln(writer, "YES")
			}
		} else {
			if dsu.Find(x) == dsu.Find(y) || dsu.Find(x+n) == dsu.Find(y+n) {
				fmt.Fprintln(writer, "NO")
			} else {
				dsu.Union(x, y+n)
				dsu.Union(x+n, y)
				fmt.Fprintln(writer, "YES")
			}
		}
	}

	for i := 0; i < q; i++ {
		var a, b string
		fmt.Fscan(reader, &a, &b)
		x := index[a]
		y := index[b]
		if dsu.Find(x) == dsu.Find(y) {
			fmt.Fprintln(writer, 1)
		} else if dsu.Find(x) == dsu.Find(y+n) {
			fmt.Fprintln(writer, 2)
		} else {
			fmt.Fprintln(writer, 3)
		}
	}
}
