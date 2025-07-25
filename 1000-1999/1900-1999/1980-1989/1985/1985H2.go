package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

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
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func uniqueSorted(a []int) []int {
	if len(a) == 0 {
		return a
	}
	sort.Ints(a)
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		N := n * m
		dsu := NewDSU(N)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] != '#' {
					continue
				}
				idx := i*m + j
				if i > 0 && grid[i-1][j] == '#' {
					dsu.Union(idx, (i-1)*m+j)
				}
				if j > 0 && grid[i][j-1] == '#' {
					dsu.Union(idx, i*m+j-1)
				}
			}
		}

		rowEmpty := make([]int, n)
		colEmpty := make([]int, m)
		rowsMap := make(map[int][]int)
		colsMap := make(map[int][]int)
		seen := make(map[int]bool)
		best := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '#' {
					idx := i*m + j
					root := dsu.Find(idx)
					rowsMap[root] = append(rowsMap[root], i)
					colsMap[root] = append(colsMap[root], j)
					if !seen[root] {
						seen[root] = true
						if dsu.size[root] > best {
							best = dsu.size[root]
						}
					}
				} else {
					rowEmpty[i]++
					colEmpty[j]++
				}
			}
		}

		ans := make([]int, N)
		for root, rows := range rowsMap {
			cols := colsMap[root]
			rows = uniqueSorted(rows)
			cols = uniqueSorted(cols)

			rowReach := make([]int, 0, len(rows)*3)
			for _, r := range rows {
				if r-1 >= 0 {
					rowReach = append(rowReach, r-1)
				}
				rowReach = append(rowReach, r)
				if r+1 < n {
					rowReach = append(rowReach, r+1)
				}
			}
			rowReach = uniqueSorted(rowReach)

			colReach := make([]int, 0, len(cols)*3)
			for _, c := range cols {
				if c-1 >= 0 {
					colReach = append(colReach, c-1)
				}
				colReach = append(colReach, c)
				if c+1 < m {
					colReach = append(colReach, c+1)
				}
			}
			colReach = uniqueSorted(colReach)

			sz := dsu.size[dsu.Find(root)]
			for _, r := range rowReach {
				base := r * m
				for _, c := range colReach {
					ans[base+c] += sz
				}
			}
		}

		for r := 0; r < n; r++ {
			for c := 0; c < m; c++ {
				total := ans[r*m+c] + rowEmpty[r] + colEmpty[c]
				if grid[r][c] == '.' {
					total--
				}
				if total > best {
					best = total
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}
