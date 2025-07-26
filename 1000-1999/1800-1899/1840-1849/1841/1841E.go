package main

import (
	"bufio"
	"fmt"
	"os"
)

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{make([]int, n), make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int, cur *int64) {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	*cur -= int64(d.size[rx]/2 + d.size[ry]/2)
	d.size[rx] += d.size[ry]
	d.parent[ry] = rx
	*cur += int64(d.size[rx] / 2)
}

func solveCase(in *bufio.Reader, out *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	var m int64
	fmt.Fscan(in, &m)

	events := make([][]int, n+2)
	for i, v := range a {
		if v < n {
			events[v+1] = append(events[v+1], i)
		}
	}

	d := NewDSU(n)
	active := make([]bool, n)
	var curPairs int64
	var totalPairs int64

	for r := 1; r <= n; r++ {
		for _, c := range events[r] {
			active[c] = true
			d.parent[c] = c
			d.size[c] = 1
			if c > 0 && active[c-1] {
				d.union(c, c-1, &curPairs)
			}
			if c+1 < n && active[c+1] {
				d.union(c, c+1, &curPairs)
			}
		}
		totalPairs += curPairs
	}

	if totalPairs > m/2 {
		totalPairs = m / 2
	}
	fmt.Fprintln(out, totalPairs)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		solveCase(in, out)
	}
}
