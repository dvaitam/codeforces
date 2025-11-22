package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

type dsu struct {
	parent []int
	sz     []int
	edges  []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	e := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{parent: p, sz: s, edges: e}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.sz[ra] += d.sz[rb]
	d.edges[ra] += d.edges[rb]
}

type variable struct {
	a, b     int
	len      int
	resolved bool
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)

		odds := make([][2]int, k+1)
		oddSet := make(map[int]struct{}, k+1)
		for i := 0; i <= k; i++ {
			fmt.Fscan(in, &odds[i][0], &odds[i][1])
			id := (odds[i][0]-1)*m + (odds[i][1] - 1)
			oddSet[id] = struct{}{}
		}

		vars := make([]variable, k)
		cellToVars := make(map[int][]int)
		ok := true

		for i := 0; i < k; i++ {
			a := odds[i]
			b := odds[i+1]
			list := make([]int, 0, 2)
			for _, d := range dirs {
				x := a[0] + d[0]
				y := a[1] + d[1]
				if x < 1 || x > n || y < 1 || y > m {
					continue
				}
				if abs(x-b[0])+abs(y-b[1]) != 1 {
					continue
				}
				id := (x-1)*m + (y - 1)
				if _, exists := oddSet[id]; exists {
					continue
				}
				list = append(list, id)
			}
			if len(list) == 0 {
				ok = false
				break
			}
			if len(list) == 1 {
				vars[i] = variable{a: list[0], len: 1}
			} else {
				vars[i] = variable{a: list[0], b: list[1], len: 2}
			}
			for _, id := range list {
				cellToVars[id] = append(cellToVars[id], i)
			}
		}

		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		used := make(map[int]bool)
		queue := make([]int, 0)
		for i, v := range vars {
			if v.len == 1 {
				queue = append(queue, i)
			}
		}

		for head := 0; head < len(queue) && ok; head++ {
			idx := queue[head]
			v := &vars[idx]
			if v.resolved || v.len != 1 {
				continue
			}
			cell := v.a
			if used[cell] {
				ok = false
				break
			}
			used[cell] = true
			v.resolved = true

			for _, other := range cellToVars[cell] {
				if vars[other].resolved {
					continue
				}
				if vars[other].len == 1 {
					if vars[other].a == cell {
						ok = false
						break
					}
				} else if vars[other].len == 2 {
					if vars[other].a == cell {
						vars[other].a = vars[other].b
						vars[other].len = 1
						queue = append(queue, other)
					} else if vars[other].b == cell {
						vars[other].len = 1
						queue = append(queue, other)
					}
				}
			}
		}

		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		for i := 0; i < k; i++ {
			if vars[i].len == 0 {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		cellID := make(map[int]int)
		getID := func(id int) int {
			if v, exists := cellID[id]; exists {
				return v
			}
			idx := len(cellID)
			cellID[id] = idx
			return idx
		}

		edges := make([][2]int, 0)
		for i := 0; i < k; i++ {
			if vars[i].resolved {
				continue
			}
			if vars[i].len != 2 {
				// Should not happen after propagation; if it does, resolve it now.
				if vars[i].len == 1 {
					if used[vars[i].a] {
						ok = false
						break
					}
					used[vars[i].a] = true
					vars[i].resolved = true
					continue
				}
				ok = false
				break
			}
			edges = append(edges, [2]int{vars[i].a, vars[i].b})
			getID(vars[i].a)
			getID(vars[i].b)
		}

		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		if len(edges) == 0 {
			fmt.Fprintln(out, 1)
			continue
		}

		d := newDSU(len(cellID))
		for _, e := range edges {
			u := getID(e[0])
			v := getID(e[1])
			d.union(u, v)
		}

		for _, e := range edges {
			u := getID(e[0])
			root := d.find(u)
			d.edges[root]++
		}

		ans := int64(1)
		for i := 0; i < len(d.parent); i++ {
			if d.find(i) != i {
				continue
			}
			e := d.edges[i]
			vCount := d.sz[i]
			if e > vCount {
				ok = false
				break
			}
			if e == vCount {
				ans = ans * 2 % mod
			} else if e == vCount-1 {
				ans = ans * int64(vCount) % mod
			} else {
				ok = false
				break
			}
		}

		if !ok {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, ans%mod)
		}
	}
}
