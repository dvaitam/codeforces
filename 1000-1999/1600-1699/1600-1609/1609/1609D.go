package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type DSU struct {
	parent []int
	size   []int
}

func newDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n+1),
		size:   make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int) bool {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return false
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	d.parent[ry] = rx
	d.size[rx] += d.size[ry]
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, d int
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}

	dsu := newDSU(n)
	extra := 0
	for i := 0; i < d; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if !dsu.union(x, y) {
			extra++
		}
		seen := make(map[int]struct{})
		comps := make([]int, 0)
		for v := 1; v <= n; v++ {
			r := dsu.find(v)
			if _, ok := seen[r]; !ok {
				seen[r] = struct{}{}
				comps = append(comps, dsu.size[r])
			}
		}
		sort.Sort(sort.Reverse(sort.IntSlice(comps)))
		k := extra + 1
		if k > len(comps) {
			k = len(comps)
		}
		sum := 0
		for j := 0; j < k; j++ {
			sum += comps[j]
		}
		fmt.Fprintln(out, sum-1)
	}
}
