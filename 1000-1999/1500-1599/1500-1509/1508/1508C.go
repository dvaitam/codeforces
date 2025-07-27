package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct{ u, v, w int }

type DSU struct{ p, r []int }

func NewDSU(n int) *DSU {
	d := &DSU{p: make([]int, n), r: make([]int, n)}
	for i := 0; i < n; i++ {
		d.p[i] = i
	}
	return d
}
func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}
func (d *DSU) Union(x, y int) bool {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return false
	}
	if d.r[x] < d.r[y] {
		x, y = y, x
	}
	d.p[y] = x
	if d.r[x] == d.r[y] {
		d.r[x]++
	}
	return true
}

func complement(n int, adj []map[int]struct{}) (*DSU, int) {
	d := NewDSU(n)
	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	next[n+1] = n + 1
	remove := func(x int) { next[prev[x]] = next[x]; prev[next[x]] = prev[x] }
	queue := make([]int, 0)
	comp := 0
	for next[0] != n+1 {
		v := next[0]
		remove(v)
		queue = append(queue, v)
		comp++
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for u := next[0]; u != n+1; {
				if _, ok := adj[cur-1][u-1]; ok {
					u = next[u]
					continue
				}
				d.Union(cur-1, u-1)
				queue = append(queue, u)
				tmp := next[u]
				remove(u)
				u = tmp
			}
		}
	}
	return d, comp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	adj := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		adj[i] = make(map[int]struct{})
	}
	edges := make([]Edge, m)
	xor := 0
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		edges[i] = Edge{u, v, w}
		adj[u][v] = struct{}{}
		adj[v][u] = struct{}{}
		xor ^= w
	}
	dsu, comp := complement(n, adj)
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	base := int64(0)
	minCycle := int(^uint(0) >> 1)
	for _, e := range edges {
		if dsu.Union(e.u, e.v) {
			base += int64(e.w)
		} else if e.w < minCycle {
			minCycle = e.w
		}
	}
	totalMissing := int64(n*(n-1)/2 - m)
	usedZero := int64(n - comp)
	if totalMissing > usedZero {
		fmt.Println(base)
	} else {
		if xor < minCycle {
			minCycle = xor
		}
		fmt.Println(base + int64(minCycle))
	}
}
