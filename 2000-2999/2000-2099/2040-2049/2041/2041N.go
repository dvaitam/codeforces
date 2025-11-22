package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type dsu struct {
	p    []int
	rank []int
	min  []int
	ptr  []int
}

func newDSU(n int, weights []int64, prison int) *dsu {
	p := make([]int, n+1)
	r := make([]int, n+1)
	min := make([]int, n+1)
	ptr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		r[i] = 0
		min[i] = i
		ptr[i] = 0
	}
	return &dsu{p: p, rank: r, min: min, ptr: ptr}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) unite(x, y int, weights []int64) {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return
	}
	if d.rank[x] < d.rank[y] {
		x, y = y, x
	}
	d.p[y] = x
	if d.rank[x] == d.rank[y] {
		d.rank[x]++
	}
	// keep the node with smaller weight as the representative for min
	if weights[d.min[y]] < weights[d.min[x]] || (weights[d.min[y]] == weights[d.min[x]] && d.min[y] < d.min[x]) {
		d.min[x] = d.min[y]
	}
	if d.ptr[y] < d.ptr[x] {
		d.ptr[x] = d.ptr[y]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	forbid := make([]map[int]struct{}, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if forbid[u] == nil {
			forbid[u] = make(map[int]struct{})
		}
		if forbid[v] == nil {
			forbid[v] = make(map[int]struct{})
		}
		forbid[u][v] = struct{}{}
		forbid[v][u] = struct{}{}
	}

	order := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		order = append(order, i)
	}
	sort.Slice(order, func(i, j int) bool {
		if a[order[i]] == a[order[j]] {
			return order[i] < order[j]
		}
		return a[order[i]] < a[order[j]]
	})

	results := make([]int64, n+1)

	// Helper to check if edge (u,v) is forbidden
	isForbidden := func(u, v int) bool {
		if forbid[u] == nil {
			return false
		}
		_, ok := forbid[u][v]
		return ok
	}

	for prison := 1; prison <= n; prison++ {
		if n == 2 {
			// only one node left, no edges needed
			results[prison] = 0
			continue
		}

		ds := newDSU(n, a, prison)
		comp := n - 1
		var cost int64
		changed := true

		for changed && comp > 1 {
			changed = false
			bestTo := make(map[int]int)
			bestW := make(map[int]int64)

			for i := 1; i <= n; i++ {
				if i == prison {
					continue
				}
				root := ds.find(i)
				if i != root {
					continue
				}
				u := ds.min[root]
				idx := ds.ptr[root]
				for idx < len(order) {
					v := order[idx]
					idx++
					if v == prison {
						continue
					}
					rv := ds.find(v)
					if rv == root {
						continue
					}
					if isForbidden(u, v) {
						continue
					}
					ds.ptr[root] = idx - 1
					w := a[u] + a[v]
					bestTo[root] = v
					bestW[root] = w
					break
				}
				ds.ptr[root] = idx
			}

			for root, v := range bestTo {
				r1 := ds.find(root)
				r2 := ds.find(v)
				if r1 == r2 {
					continue
				}
				cost += bestW[root]
				ds.unite(r1, r2, a)
				comp--
				changed = true
			}
		}

		if comp != 1 {
			results[prison] = -1
		} else {
			results[prison] = cost
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, results[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
