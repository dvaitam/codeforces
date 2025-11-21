package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int64
	idx int
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int, f func(int) int64, total *int64) {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	*total -= f(d.size[rx])
	*total -= f(d.size[ry])
	d.parent[ry] = rx
	d.size[rx] += d.size[ry]
	*total += f(d.size[rx])
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, d int
		fmt.Fscan(in, &n, &d)
		a := make([]int64, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}

		pairs := make([]Pair, n)
		for i := 0; i < n; i++ {
			pairs[i] = Pair{val: a[i], idx: i}
		}
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].val > pairs[j].val
		})

		total := int64(0)
		active := make([]bool, n)
		dsu := NewDSU(n)
		f := func(sz int) int64 {
			return int64((sz + d - 1) / d)
		}

		ans := int64(0)
		ptr := 0
		var curr int64 = 0
		if n > 0 {
			curr = pairs[0].val
		}
		for ptr < n {
			curr = pairs[ptr].val
			for ptr < n && pairs[ptr].val == curr {
				i := pairs[ptr].idx
				active[i] = true
				dsu.parent[i] = i
				dsu.size[i] = 1
				total += f(1)
				if i > 0 && active[i-1] {
					dsu.union(i, i-1, f, &total)
				}
				if i+1 < n && active[i+1] {
					dsu.union(i, i+1, f, &total)
				}
				ptr++
			}
			next := int64(0)
			if ptr < n {
				next = pairs[ptr].val
			}
			delta := curr - next
			ans += delta * total
		}
		fmt.Fprintln(out, ans)
	}
}
