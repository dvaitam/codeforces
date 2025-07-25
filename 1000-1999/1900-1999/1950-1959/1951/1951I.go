package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type edge struct {
	u, v int
	a, b int64
	idx  int
}

type dsu struct {
	parent []int
	rank   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
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
		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			var u, v int
			var a, b int64
			fmt.Fscan(reader, &u, &v, &a, &b)
			edges[i] = edge{u - 1, v - 1, a, b, i}
		}
		used := make([]int64, m)
		var total int64
		for step := int64(0); step < k; step++ {
			weights := make([]struct {
				w int64
				e edge
			}, m)
			for i, e := range edges {
				w := e.a*(2*used[i]+1) + e.b
				weights[i] = struct {
					w int64
					e edge
				}{w, e}
			}
			sort.Slice(weights, func(i, j int) bool {
				return weights[i].w < weights[j].w
			})
			d := newDSU(n)
			cnt := 0
			for _, item := range weights {
				if d.union(item.e.u, item.e.v) {
					used[item.e.idx]++
					total += item.w
					cnt++
					if cnt == n-1 {
						break
					}
				}
			}
		}
		fmt.Fprintln(writer, total)
	}
}
