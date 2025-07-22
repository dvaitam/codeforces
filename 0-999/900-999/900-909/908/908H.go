package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU structure for union-find
type DSU struct {
	parent []int
}

// NewDSU initializes a DSU for elements 0..n-1
func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{parent: p}
}

// Find returns representative of x with path compression
func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

// Union merges sets containing x and y
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

	matrix := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &matrix[i])
	}

	dsu := NewDSU(n)

	// merge nodes which must be strongly connected
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if matrix[i][j] == 'A' {
				dsu.Union(i, j)
			}
		}
	}

	// check for contradictions: XOR inside same component
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if matrix[i][j] == 'X' && dsu.Find(i) == dsu.Find(j) {
				fmt.Fprintln(writer, -1)
				return
			}
		}
	}

	// compute component sizes
	compSize := make(map[int]int)
	for i := 0; i < n; i++ {
		root := dsu.Find(i)
		compSize[root]++
	}

	components := len(compSize)
	edges := 0
	for _, sz := range compSize {
		if sz > 1 {
			edges += sz
		}
	}
	if components > 0 {
		edges += components - 1
	}

	fmt.Fprintln(writer, edges)
}
