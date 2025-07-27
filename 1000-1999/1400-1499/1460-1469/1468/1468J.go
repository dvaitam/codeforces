package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int64
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p, sz}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) bool {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		var k int64
		fmt.Fscan(reader, &n, &m, &k)
		edges := make([]Edge, m)
		maxLess := int64(-1)
		minGreater := int64(1 << 62)
		for i := 0; i < m; i++ {
			var x, y int
			var w int64
			fmt.Fscan(reader, &x, &y, &w)
			edges[i] = Edge{x - 1, y - 1, w}
			if w < k {
				if w > maxLess {
					maxLess = w
				}
			} else {
				if w-k < minGreater {
					minGreater = w - k
				}
			}
		}
		less := make([]Edge, 0)
		greater := make([]Edge, 0)
		for _, e := range edges {
			if e.w < k {
				less = append(less, e)
			} else {
				greater = append(greater, e)
			}
		}
		sort.Slice(less, func(i, j int) bool { return less[i].w < less[j].w })
		sort.Slice(greater, func(i, j int) bool { return greater[i].w < greater[j].w })
		dsu := NewDSU(n)
		cnt := 0
		cost := int64(0)
		for _, e := range less {
			if dsu.Union(e.u, e.v) {
				cnt++
				if cnt == n-1 {
					break
				}
			}
		}
		if cnt < n-1 {
			for _, e := range greater {
				if dsu.Union(e.u, e.v) {
					cost += e.w - k
					cnt++
					if cnt == n-1 {
						break
					}
				}
			}
		}
		if cost > 0 {
			fmt.Fprintln(writer, cost)
		} else {
			// all edges used have weight < k
			addLess := int64(1 << 62)
			if maxLess != -1 {
				addLess = k - maxLess
			}
			ans := addLess
			if minGreater < ans {
				ans = minGreater
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
