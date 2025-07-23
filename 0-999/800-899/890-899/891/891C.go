package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v, w int
}

type DSU struct {
	parent []int
	rank   []int
}

func newDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) union(x, y int) bool {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return false
	}
	if d.rank[rx] < d.rank[ry] {
		rx, ry = ry, rx
	}
	d.parent[ry] = rx
	if d.rank[rx] == d.rank[ry] {
		d.rank[rx]++
	}
	return true
}

// QueryInfo stores edges of a query with the same weight.
type QueryInfo struct {
	q     int
	edges []int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([]Edge, m)
	maxW := 0
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
		edges[i].u--
		edges[i].v--
		if edges[i].w > maxW {
			maxW = edges[i].w
		}
	}
	var q int
	fmt.Fscan(reader, &q)

	ans := make([]bool, q)
	for i := range ans {
		ans[i] = true
	}

	// queriesByWeight groups queries by edge weight
	queriesByWeight := make(map[int][]QueryInfo)

	for qi := 0; qi < q; qi++ {
		var k int
		fmt.Fscan(reader, &k)
		temp := make(map[int][]int)
		for j := 0; j < k; j++ {
			var idx int
			fmt.Fscan(reader, &idx)
			idx--
			w := edges[idx].w
			temp[w] = append(temp[w], idx)
		}
		for w, list := range temp {
			copied := make([]int, len(list))
			copy(copied, list)
			queriesByWeight[w] = append(queriesByWeight[w], QueryInfo{q: qi, edges: copied})
		}
	}

	// edges grouped by weight for MST
	edgesByWeight := make(map[int][]int)
	for i, e := range edges {
		edgesByWeight[e.w] = append(edgesByWeight[e.w], i)
	}

	// collect all unique weights
	weights := make([]int, 0, len(edgesByWeight))
	for w := range edgesByWeight {
		weights = append(weights, w)
	}
	sort.Ints(weights)

	dsu := newDSU(n)

	for _, w := range weights {
		// process queries involving edges of weight w
		for _, info := range queriesByWeight[w] {
			if !ans[info.q] {
				continue
			}
			compID := make(map[int]int)
			pairs := make([][2]int, len(info.edges))
			for idx, eIdx := range info.edges {
				e := edges[eIdx]
				u := dsu.find(e.u)
				v := dsu.find(e.v)
				if _, ok := compID[u]; !ok {
					compID[u] = len(compID)
				}
				if _, ok := compID[v]; !ok {
					compID[v] = len(compID)
				}
				pairs[idx] = [2]int{compID[u], compID[v]}
			}
			tmp := newDSU(len(compID))
			ok := true
			for _, p := range pairs {
				if !tmp.union(p[0], p[1]) {
					ok = false
					break
				}
			}
			if !ok {
				ans[info.q] = false
			}
		}
		// merge edges of weight w into global DSU
		for _, idx := range edgesByWeight[w] {
			e := edges[idx]
			dsu.union(e.u, e.v)
		}
	}

	for i := 0; i < q; i++ {
		if ans[i] {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
