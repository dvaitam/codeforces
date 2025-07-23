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
	parent := make([]int, n)
	size := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		size[i] = 1
	}
	return &DSU{parent: parent, size: size}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return
	}
	if d.size[x] < d.size[y] {
		x, y = y, x
	}
	d.parent[y] = x
	d.size[x] += d.size[y]
}

func (d *DSU) Size(x int) int {
	return d.size[d.Find(x)]
}

type Cell struct {
	val int
	r   int
	c   int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	cells := make([]Cell, 0, n*m)
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &grid[i][j])
			cells = append(cells, Cell{val: grid[i][j], r: i, c: j})
		}
	}

	sort.Slice(cells, func(i, j int) bool {
		return cells[i].val > cells[j].val
	})

	dsu := NewDSU(n * m)
	active := make([]bool, n*m)
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	idx := 0
	for idx < len(cells) {
		h := cells[idx].val
		start := idx
		for idx < len(cells) && cells[idx].val == h {
			r := cells[idx].r
			c := cells[idx].c
			id := r*m + c
			active[id] = true
			for _, d := range dirs {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < n && nc >= 0 && nc < m {
					nid := nr*m + nc
					if active[nid] {
						dsu.Union(id, nid)
					}
				}
			}
			idx++
		}
		if k%int64(h) == 0 {
			need := int(k / int64(h))
			if need <= n*m {
				for i := start; i < idx; i++ {
					r := cells[i].r
					c := cells[i].c
					id := r*m + c
					root := dsu.Find(id)
					if dsu.Size(root) >= need {
						result := make([][]int, n)
						for p := range result {
							result[p] = make([]int, m)
						}
						vis := make([][]bool, n)
						for p := range vis {
							vis[p] = make([]bool, m)
						}
						queue := make([][2]int, 0)
						queue = append(queue, [2]int{r, c})
						vis[r][c] = true
						result[r][c] = h
						cnt := 1
						for head := 0; head < len(queue) && cnt < need; head++ {
							cr := queue[head][0]
							cc := queue[head][1]
							for _, d := range dirs {
								nr, nc := cr+d[0], cc+d[1]
								if nr >= 0 && nr < n && nc >= 0 && nc < m && !vis[nr][nc] {
									nid := nr*m + nc
									if active[nid] && dsu.Find(nid) == root {
										vis[nr][nc] = true
										result[nr][nc] = h
										queue = append(queue, [2]int{nr, nc})
										cnt++
										if cnt == need {
											break
										}
									}
								}
							}
						}
						if cnt == need {
							fmt.Fprintln(writer, "YES")
							for x := 0; x < n; x++ {
								for y := 0; y < m; y++ {
									if y > 0 {
										fmt.Fprint(writer, " ")
									}
									fmt.Fprint(writer, result[x][y])
								}
								fmt.Fprintln(writer)
							}
							return
						}
					}
				}
			}
		}
	}
	fmt.Fprintln(writer, "NO")
}
