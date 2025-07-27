package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	u, v int
	w    int
}

type DSU struct {
	p []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{p}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) bool {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return false
	}
	d.p[ra] = rb
	return true
}

func buildMST(n int, edges []Edge) ([]Edge, int, []int) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].w < edges[j].w
	})
	dsu := NewDSU(n)
	mst := make([]Edge, 0, n-1)
	deg := make([]int, n)
	cost := 0
	for _, e := range edges {
		if dsu.Union(e.u, e.v) {
			mst = append(mst, e)
			deg[e.u]++
			deg[e.v]++
			cost += e.w
			if len(mst) == n-1 {
				break
			}
		}
	}
	return mst, cost, deg
}

func componentsAfterRemoving(n int, edges []Edge, skip int) []int {
	g := make([][]int, n)
	for idx, e := range edges {
		if idx == skip {
			continue
		}
		g[e.u] = append(g[e.u], e.v)
		g[e.v] = append(g[e.v], e.u)
	}
	comp := make([]int, n)
	queue := []int{edges[skip].u}
	comp[edges[skip].u] = 1
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range g[v] {
			if comp[to] == 0 {
				comp[to] = 1
				queue = append(queue, to)
			}
		}
	}
	for i := 0; i < n; i++ {
		if comp[i] == 0 {
			comp[i] = 2
		}
	}
	return comp
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	d := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &d[i])
	}

	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
	}
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			fmt.Fscan(reader, &w[i][j])
			w[j][i] = w[i][j]
		}
	}

	edges := make([]Edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{u: i, v: j, w: w[i][j]})
		}
	}

	tree, cost, deg := buildMST(n, edges)

	for {
		ok := true
		for i := 0; i < k; i++ {
			if deg[i] > d[i] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, cost)
			return
		}

		bestDiff := int(1e9)
		removeIdx := -1
		var addEdge Edge

		for idx, e := range tree {
			if (e.u < k && deg[e.u] > d[e.u]) || (e.v < k && deg[e.v] > d[e.v]) {
				comp := componentsAfterRemoving(n, tree, idx)
				for a := 0; a < n; a++ {
					for b := a + 1; b < n; b++ {
						if comp[a] == comp[b] {
							continue
						}
						if a == e.u && b == e.v || a == e.v && b == e.u {
							continue
						}
						if a < k && deg[a]+1 > d[a] {
							continue
						}
						if b < k && deg[b]+1 > d[b] {
							continue
						}
						diff := w[a][b] - e.w
						if diff < bestDiff {
							bestDiff = diff
							removeIdx = idx
							addEdge = Edge{u: a, v: b, w: w[a][b]}
						}
					}
				}
			}
		}

		if removeIdx == -1 || bestDiff >= 1e9 {
			// cannot satisfy constraints
			fmt.Fprintln(writer, cost)
			return
		}

		// apply replacement
		rem := tree[removeIdx]
		tree = append(tree[:removeIdx], tree[removeIdx+1:]...)
		deg[rem.u]--
		deg[rem.v]--
		cost -= rem.w

		tree = append(tree, addEdge)
		deg[addEdge.u]++
		deg[addEdge.v]++
		cost += addEdge.w
	}
}
