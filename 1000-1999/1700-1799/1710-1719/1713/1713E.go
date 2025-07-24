package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	xor    []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	xor := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &DSU{parent: parent, xor: xor}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		px := d.find(d.parent[x])
		d.xor[x] ^= d.xor[d.parent[x]]
		d.parent[x] = px
	}
	return d.parent[x]
}

func (d *DSU) union(a, b, v int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	// set xor between a and b to v
	d.parent[ra] = rb
	d.xor[ra] = d.xor[a] ^ d.xor[b] ^ v
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &row[j])
			}
			a[i] = row
		}

		dsu := NewDSU(n)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if a[i][j] == a[j][i] {
					continue
				}
				want := 0
				if a[i][j] > a[j][i] {
					want = 1
				}
				if dsu.find(i) != dsu.find(j) {
					dsu.union(i, j, want)
				}
			}
		}
		orient := make([]int, n)
		for i := 0; i < n; i++ {
			dsu.find(i)
			orient[i] = dsu.xor[i]
		}

		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				val := a[i][j]
				if i < j && (orient[i]^orient[j]) == 1 {
					val = a[j][i]
				} else if j < i && (orient[i]^orient[j]) == 1 {
					val = a[j][i]
				}
				if j > 0 {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, val)
			}
			fmt.Fprintln(writer)
		}
	}
}
