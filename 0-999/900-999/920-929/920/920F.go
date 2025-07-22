package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxA = 1000000

// Fenwick tree for int64 prefix sums
type Fenwick struct {
	n    int
	tree []int64
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int64, n+2)}
}

func (f *Fenwick) Add(i int, v int64) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] += v
	}
}

func (f *Fenwick) Sum(i int) int64 {
	var s int64
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}

func (f *Fenwick) RangeSum(l, r int) int64 {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

// DSU structure to skip indices that are already stable
type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+2)
	for i := range p {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	fx := d.Find(x)
	fy := d.Find(y)
	if fx != fy {
		d.parent[fx] = fy
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	// precompute divisor counts
	divCnt := make([]int, MaxA+1)
	for i := 1; i <= MaxA; i++ {
		for j := i; j <= MaxA; j += i {
			divCnt[j]++
		}
	}

	fenw := NewFenwick(n)
	for i := 1; i <= n; i++ {
		fenw.Add(i, int64(arr[i]))
	}

	dsu := NewDSU(n)
	for i := 1; i <= n; i++ {
		if arr[i] <= 2 {
			dsu.Union(i, i+1)
		}
	}

	for ; m > 0; m-- {
		var t, l, r int
		fmt.Fscan(in, &t, &l, &r)
		if t == 1 {
			pos := dsu.Find(l)
			for pos <= r {
				newVal := divCnt[arr[pos]]
				if newVal != arr[pos] {
					fenw.Add(pos, int64(newVal-arr[pos]))
					arr[pos] = newVal
				}
				if arr[pos] <= 2 {
					dsu.Union(pos, pos+1)
				}
				pos = dsu.Find(pos + 1)
			}
		} else {
			fmt.Fprintln(out, fenw.RangeSum(l, r))
		}
	}
}
