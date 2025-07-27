package main

import (
	"bufio"
	"fmt"
	"os"
)

// Note: This is a partial implementation for problem F.
// The algorithm maintains connectivity using a weighted DSU
// and ensures each added cycle has xor weight 1 while each
// tree edge participates in at most one cycle. The correctness
// is not fully proven.

type DSU struct {
	p    []int
	rank []int
	xor  []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	x := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p: p, rank: r, xor: x}
}

func (d *DSU) Find(x int) (int, int) {
	if d.p[x] != x {
		root, xr := d.Find(d.p[x])
		d.xor[x] ^= xr
		d.p[x] = root
	}
	return d.p[x], d.xor[x]
}

func (d *DSU) Union(u, v, w int) {
	ru, xu := d.Find(u)
	rv, xv := d.Find(v)
	if ru == rv {
		return
	}
	if d.rank[ru] < d.rank[rv] {
		ru, rv = rv, ru
		xu, xv = xv, xu
	}
	d.p[rv] = ru
	d.xor[rv] = xu ^ xv ^ w
	if d.rank[ru] == d.rank[rv] {
		d.rank[ru]++
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	d := NewDSU(n)
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	used := make([]bool, n+1)
	skip := make([]int, n+1)
	for i := 1; i <= n; i++ {
		skip[i] = i
	}

	var get func(int) int
	get = func(x int) int {
		if skip[x] == x {
			return x
		}
		skip[x] = get(skip[x])
		return skip[x]
	}

	for ; q > 0; q-- {
		var u, v, x int
		fmt.Fscan(in, &u, &v, &x)
		ru, xu := d.Find(u)
		rv, xv := d.Find(v)
		if ru != rv {
			fmt.Fprintln(out, "YES")
			d.Union(u, v, x)
			// attach tree edge between component roots
			if depth[ru] < depth[rv] {
				ru, rv = rv, ru
			}
			parent[rv] = ru
			depth[rv] = depth[ru] + 1
			skip[rv] = rv
		} else {
			pathXor := xu ^ xv
			if pathXor^x != 1 {
				fmt.Fprintln(out, "NO")
				continue
			}
			// check path from u to v for unused edges
			x1 := u
			x2 := v
			ok := true
			for get(x1) != get(x2) {
				if depth[get(x1)] < depth[get(x2)] {
					x1, x2 = x2, x1
				}
				cur := get(x1)
				if used[cur] {
					ok = false
					break
				}
				used[cur] = true
				skip[cur] = parent[cur]
				x1 = get(cur)
			}
			if ok {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
