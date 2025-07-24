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
	for i := range d.p {
		d.p[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) union(a, b int) {
	fa := d.find(a)
	fb := d.find(b)
	if fa != fb {
		d.p[fa] = fb
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		m := 2 * n
		pref := make([]int, m+1)
		maxh := 0
		for i := 1; i <= m; i++ {
			if s[i-1] == '(' {
				pref[i] = pref[i-1] + 1
			} else {
				pref[i] = pref[i-1] - 1
			}
			if pref[i] > maxh {
				maxh = pref[i]
			}
		}

		dsu := NewDSU(m + 1)
		first := make([]int, maxh+2)
		last := make([]int, maxh+2)
		low := make([]int, maxh+2)
		for i := 0; i < len(first); i++ {
			first[i] = -1
			last[i] = -1
			low[i] = -1
		}
		first[0] = 0
		last[0] = 0

		for i := 1; i <= m; i++ {
			prev := pref[i-1]
			d := pref[i]
			if d < prev {
				low[prev] = i - 1
			}
			if last[d] > low[d] {
				dsu.union(last[d]+1, i)
				dsu.union(first[d]+1, i)
			} else {
				first[d] = i
			}
			last[d] = i
		}

		comp := make(map[int]struct{})
		for i := 1; i <= m; i++ {
			comp[dsu.find(i)] = struct{}{}
		}
		fmt.Fprintln(out, len(comp))
	}
}
