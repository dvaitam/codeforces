package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct{ parent map[int]int }

func NewDSU() *DSU { return &DSU{parent: make(map[int]int)} }
func (d *DSU) find(x int) int {
	if p, ok := d.parent[x]; ok {
		if p != x {
			d.parent[x] = d.find(p)
		}
		return d.parent[x]
	}
	return x
}
func (d *DSU) union(x, y int) {
	x = d.find(x)
	y = d.find(y)
	if x != y {
		d.parent[x] = y
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(reader, &n, &q)

	rowDSU := NewDSU()
	colDSU := NewDSU()
	rowBound := make(map[int]int)
	colBound := make(map[int]int)

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for i := 0; i < q; i++ {
		var x, y int
		var dir string
		fmt.Fscan(reader, &x, &y, &dir)
		if dir == "U" {
			lim := colBound[x]
			if y <= lim {
				fmt.Fprintln(writer, 0)
				continue
			}
			ans := y - lim
			r := rowDSU.find(y)
			for r > lim {
				rowDSU.union(r, r-1)
				rowBound[r] = x
				r = rowDSU.find(r)
			}
			colBound[x] = y
			fmt.Fprintln(writer, ans)
		} else { // L
			lim := rowBound[y]
			if x <= lim {
				fmt.Fprintln(writer, 0)
				continue
			}
			ans := x - lim
			c := colDSU.find(x)
			for c > lim {
				colDSU.union(c, c-1)
				colBound[c] = y
				c = colDSU.find(c)
			}
			rowBound[y] = x
			fmt.Fprintln(writer, ans)
		}
	}
}
