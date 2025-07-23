package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	w int
	u int
	v int
}

type DSU struct {
	parent  []int
	size    []int
	members [][]int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	mem := make([][]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		sz[i] = 1
		mem[i] = []int{i}
	}
	return &DSU{parent: p, size: sz, members: mem}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) int {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return x
	}
	if d.size[x] < d.size[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.size[x] += d.size[y]
	d.members[x] = append(d.members[x], d.members[y]...)
	d.members[y] = nil
	return x
}

func isMagic(n int, a [][]int) bool {
	for i := 0; i < n; i++ {
		if a[i][i] != 0 {
			return false
		}
		for j := i + 1; j < n; j++ {
			if a[i][j] != a[j][i] {
				return false
			}
		}
	}

	edges := make([]Edge, 0, n*(n-1)/2)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{a[i][j], i, j})
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

	dsu := NewDSU(n)
	for idx := 0; idx < len(edges); {
		w := edges[idx].w
		j := idx
		for j < len(edges) && edges[j].w == w {
			j++
		}
		group := edges[idx:j]
		idx = j

		adj := make(map[int][]int)
		roots := make(map[int]struct{})
		for _, e := range group {
			ru := dsu.Find(e.u)
			rv := dsu.Find(e.v)
			if ru == rv {
				continue
			}
			adj[ru] = append(adj[ru], rv)
			adj[rv] = append(adj[rv], ru)
			roots[ru] = struct{}{}
			roots[rv] = struct{}{}
		}
		visited := make(map[int]bool)
		for r := range roots {
			if visited[r] {
				continue
			}
			stack := []int{r}
			visited[r] = true
			comps := []int{}
			for len(stack) > 0 {
				x := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				comps = append(comps, x)
				for _, y := range adj[x] {
					if !visited[y] {
						visited[y] = true
						stack = append(stack, y)
					}
				}
			}
			if len(comps) <= 1 {
				continue
			}
			verts := []int{}
			for _, c := range comps {
				verts = append(verts, dsu.members[c]...)
			}
			for i := 0; i < len(verts); i++ {
				for j := i + 1; j < len(verts); j++ {
					if a[verts[i]][verts[j]] > w {
						return false
					}
				}
			}
			base := comps[0]
			for _, c := range comps[1:] {
				base = dsu.Union(base, c)
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}
	if isMagic(n, a) {
		fmt.Fprintln(out, "MAGIC")
	} else {
		fmt.Fprintln(out, "NOT MAGIC")
	}
}
