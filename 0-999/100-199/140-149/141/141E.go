package main

import (
	"bufio"
	"fmt"
	"os"
)

// DSU implements disjoint set union
type DSU struct {
	p []int
}

// NewDSU creates a DSU for 1..n
func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p}
}

// Find with path compression
func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

// Union merges x and y
func (d *DSU) Union(x, y int) {
	fx := d.Find(x)
	fy := d.Find(y)
	if fx != fy {
		d.p[fx] = fy
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
	if n%2 == 0 {
		fmt.Fprintln(writer, -1)
		return
	}

	type Edge struct {
		u, v int
		s    bool
	}
	edges := make([]Edge, m)
	initialUsd := make([]bool, m)
	usd := make([]bool, m)
	dsu := NewDSU(n)
	num := 0

	// read and initial spanning forest
	for i := 0; i < m; i++ {
		var u, v int
		var str string
		fmt.Fscan(reader, &u, &v, &str)
		isS := (str[0] == 'S')
		edges[i] = Edge{u, v, isS}
		if dsu.Find(u) != dsu.Find(v) {
			dsu.Union(u, v)
			initialUsd[i] = true
			if isS {
				num++
			}
		}
	}
	copy(usd, initialUsd)

	// adjust to get exactly (n-1)/2 edges of type 'S'
	if num*2+1 != n {
		// rebuild DSU
		dsu = NewDSU(n)
		if num*2+1 < n {
			// need more 'S'
			newUsd := make([]bool, m)
			newNum := 0
			// keep initial 'S' edges
			for i := 0; i < m; i++ {
				if edges[i].s && initialUsd[i] {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			// add more 'S'
			for i := 0; i < m && newNum*2+1 != n; i++ {
				if edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			// add other edges
			for i := 0; i < m; i++ {
				if !edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		} else {
			// need fewer 'S'
			newUsd := make([]bool, m)
			newNum := num
			// keep initial non-'S'
			for i := 0; i < m; i++ {
				if !edges[i].s && initialUsd[i] {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			// add more non-'S'
			for i := 0; i < m && newNum*2+1 != n; i++ {
				if !edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum--
				}
			}
			// add 'S' edges
			for i := 0; i < m; i++ {
				if edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		}
	}

	// output result
	if num*2+1 == n {
		fmt.Fprintln(writer, n-1)
		cnt := 0
		for i := 0; i < m && cnt < n-1; i++ {
			if usd[i] {
				fmt.Fprint(writer, i+1)
				fmt.Fprint(writer, " ")
				cnt++
			}
		}
		fmt.Fprintln(writer)
	} else {
		fmt.Fprintln(writer, -1)
	}
}
