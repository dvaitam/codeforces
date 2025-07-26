package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// DSU with union by size and path compression
type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
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

func (d *DSU) union(x, y int) int {
	x = d.find(x)
	y = d.find(y)
	if x == y {
		return x
	}
	if d.size[x] < d.size[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.size[x] += d.size[y]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		adj := make([][]int, n)
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			edges[i] = [2]int{u, v}
		}

		dsu := NewDSU(n)
		added := make([]bool, n)
		active := make([]bool, n)
		compSize := make([]int, n)

		// union zero vertices connected directly
		for _, e := range edges {
			u, v := e[0], e[1]
			if a[u] == 0 && a[v] == 0 {
				dsu.union(u, v)
			}
		}
		for i := 0; i < n; i++ {
			if a[i] == 0 {
				added[i] = true
				r := dsu.find(i)
				active[r] = true
				compSize[r]++
			}
		}
		k := 0
		for i := 0; i < n; i++ {
			if active[i] && compSize[i] > k {
				k = compSize[i]
			}
		}
		if k == 0 {
			fmt.Fprintln(out, "NO")
			continue
		}

		order := make([]int, n)
		for i := 0; i < n; i++ {
			order[i] = i
		}
		sort.Slice(order, func(i, j int) bool { return a[order[i]] < a[order[j]] })
		pos := 0
		for pos < n && a[order[pos]] == 0 {
			pos++
		}

		changed := true
		for changed {
			changed = false
			for pos < n && a[order[pos]] <= k {
				v := order[pos]
				pos++
				if added[v] {
					continue
				}
				added[v] = true
				r := dsu.find(v)
				compSize[r]++
				if a[v] == 0 {
					active[r] = true
				}
				for _, to := range adj[v] {
					if added[to] {
						r = dsu.find(r)
						rt := dsu.find(to)
						if r != rt {
							nr := dsu.union(r, rt)
							if nr == r {
								active[r] = active[r] || active[rt]
								compSize[r] += compSize[rt]
								active[rt] = false
								compSize[rt] = 0
							} else {
								active[rt] = active[r] || active[rt]
								compSize[rt] += compSize[r]
								active[r] = false
								compSize[r] = 0
								r = rt
							}
						}
					}
				}
				r = dsu.find(r)
				if active[r] && compSize[r] > k {
					k = compSize[r]
					changed = true
				}
			}
		}
		if k == n {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
