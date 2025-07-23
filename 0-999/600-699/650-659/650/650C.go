package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(x, y int) {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx != ry {
		d.parent[rx] = ry
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	cells := make([][]int, n)
	for i := 0; i < n; i++ {
		cells[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &cells[i][j])
		}
	}

	total := n * m
	dsu := NewDSU(total)

	id := func(i, j int) int { return i*m + j }

	// union equal values in rows
	for i := 0; i < n; i++ {
		arr := make([]pair, m)
		for j := 0; j < m; j++ {
			arr[j] = pair{cells[i][j], id(i, j)}
		}
		sort.Slice(arr, func(a, b int) bool { return arr[a].val < arr[b].val })
		for k := 1; k < m; k++ {
			if arr[k].val == arr[k-1].val {
				dsu.Union(arr[k].idx, arr[k-1].idx)
			}
		}
	}

	// union equal values in columns
	for j := 0; j < m; j++ {
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{cells[i][j], id(i, j)}
		}
		sort.Slice(arr, func(a, b int) bool { return arr[a].val < arr[b].val })
		for k := 1; k < n; k++ {
			if arr[k].val == arr[k-1].val {
				dsu.Union(arr[k].idx, arr[k-1].idx)
			}
		}
	}

	edges := make([][]int, total)
	indeg := make([]int, total)
	edgeSet := make(map[[2]int]struct{})

	// build edges from rows
	for i := 0; i < n; i++ {
		arr := make([]pair, m)
		for j := 0; j < m; j++ {
			arr[j] = pair{cells[i][j], dsu.Find(id(i, j))}
		}
		sort.Slice(arr, func(a, b int) bool {
			if arr[a].val == arr[b].val {
				return arr[a].idx < arr[b].idx
			}
			return arr[a].val < arr[b].val
		})
		for k := 1; k < m; k++ {
			if arr[k].val > arr[k-1].val {
				u := arr[k-1].idx
				v := arr[k].idx
				if u != v {
					key := [2]int{u, v}
					if _, ok := edgeSet[key]; !ok {
						edgeSet[key] = struct{}{}
						edges[u] = append(edges[u], v)
						indeg[v]++
					}
				}
			}
		}
	}

	// build edges from columns
	for j := 0; j < m; j++ {
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{cells[i][j], dsu.Find(id(i, j))}
		}
		sort.Slice(arr, func(a, b int) bool {
			if arr[a].val == arr[b].val {
				return arr[a].idx < arr[b].idx
			}
			return arr[a].val < arr[b].val
		})
		for k := 1; k < n; k++ {
			if arr[k].val > arr[k-1].val {
				u := arr[k-1].idx
				v := arr[k].idx
				if u != v {
					key := [2]int{u, v}
					if _, ok := edgeSet[key]; !ok {
						edgeSet[key] = struct{}{}
						edges[u] = append(edges[u], v)
						indeg[v]++
					}
				}
			}
		}
	}

	// gather unique groups
	repMap := make(map[int]struct{})
	reps := make([]int, 0)
	group := make([]int, total)
	for idx := 0; idx < total; idx++ {
		g := dsu.Find(idx)
		group[idx] = g
		if _, ok := repMap[g]; !ok {
			repMap[g] = struct{}{}
			reps = append(reps, g)
		}
	}

	value := make([]int, total)
	for _, g := range reps {
		value[g] = 1
	}
	queue := make([]int, 0, len(reps))
	for _, g := range reps {
		if indeg[g] == 0 {
			queue = append(queue, g)
		}
	}

	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, v := range edges[u] {
			if value[v] < value[u]+1 {
				value[v] = value[u] + 1
			}
			indeg[v]--
			if indeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				fmt.Fprint(writer, " ")
			}
			idx := id(i, j)
			fmt.Fprint(writer, value[group[idx]])
		}
		if i+1 < n {
			fmt.Fprintln(writer)
		}
	}
}
