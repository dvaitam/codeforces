package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct{ parent []int }

func newDsu(n int) *dsu {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{p}
}
func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}
func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra != rb {
		d.parent[ra] = rb
	}
}

func build(p []int, desc bool) []int {
	n := len(p)
	B := append([]int(nil), p...)
	d := newDsu(n)
	for i := 0; i < n; i++ {
		d.union(i, p[i]-1)
	}
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		if desc {
			return p[idx[i]] > p[idx[j]]
		}
		return p[idx[i]] < p[idx[j]]
	})
	for i := 1; i < n; i++ {
		u, v := idx[i-1], idx[i]
		if d.find(u) != d.find(v) {
			B[u], B[v] = B[v], B[u]
			d.union(u, v)
		}
	}
	succ := make([]int, n)
	for i, v := range B {
		succ[v-1] = i
	}
	res := make([]int, 0, n)
	start := 0
	res = append(res, start+1)
	x := succ[start]
	for x != start {
		res = append(res, x+1)
		x = succ[x]
	}
	return res
}

func less(a, b []int) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func solve(p []int) []int {
	q1 := build(p, false)
	q2 := build(p, true)
	if less(q2, q1) {
		return q2
	}
	return q1
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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		ans := solve(p)
		for i := 0; i < len(ans); i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
