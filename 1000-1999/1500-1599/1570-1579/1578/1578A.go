package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cell struct{ r, c int }

type comp struct {
	letter byte
	cells  []cell
	minR   int
	spawn  []cell
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	comps := make([]*comp, 0)
	tileAt := make([][]int, n)
	for i := range tileAt {
		tileAt[i] = make([]int, m)
	}
	for r := 0; r < n; r++ {
		for c := 0; c < m; c++ {
			if grid[r][c] == '.' || visited[r][c] {
				continue
			}
			letter := grid[r][c]
			queue := []cell{{r, c}}
			visited[r][c] = true
			cc := &comp{letter: letter, minR: n + 1}
			for len(queue) > 0 {
				p := queue[0]
				queue = queue[1:]
				cc.cells = append(cc.cells, cell{p.r + 1, p.c + 1})
				if p.r+1 < cc.minR {
					cc.minR = p.r + 1
				}
				tileAt[p.r][p.c] = len(comps)
				for _, d := range []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					nr, nc := p.r+d.r, p.c+d.c
					if nr >= 0 && nr < n && nc >= 0 && nc < m && !visited[nr][nc] && grid[nr][nc] == letter {
						visited[nr][nc] = true
						queue = append(queue, cell{nr, nc})
					}
				}
			}
			comps = append(comps, cc)
		}
	}
	// compute spawn cells
	for _, cc := range comps {
		for _, p := range cc.cells {
			cc.spawn = append(cc.spawn, cell{p.r - cc.minR + 1, p.c})
		}
	}
	// build graph
	g := make([]map[int]struct{}, len(comps))
	indeg := make([]int, len(comps))
	for i := range g {
		g[i] = make(map[int]struct{})
	}
	for i, cc := range comps {
		for _, p := range cc.spawn {
			r := p.r - 1
			c := p.c - 1
			if r >= 0 && r < n && c >= 0 && c < m {
				j := tileAt[r][c]
				if j >= 0 && j != i {
					if _, ok := g[i][j]; !ok {
						g[i][j] = struct{}{}
						indeg[j]++
					}
				}
			}
		}
	}
	// topological sort
	order := make([]int, 0, len(comps))
	q := make([]int, 0)
	for i := 0; i < len(comps); i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for len(q) > 0 {
		i := q[0]
		q = q[1:]
		order = append(order, i)
		for j := range g[i] {
			indeg[j]--
			if indeg[j] == 0 {
				q = append(q, j)
			}
		}
	}
	if len(order) != len(comps) {
		fmt.Println(-1)
		return
	}
	fmt.Println(len(comps))
	for _, idx := range order {
		cc := comps[idx]
		x := m + 1
		for _, p := range cc.spawn {
			if p.r == 1 && p.c < x {
				x = p.c
			}
		}
		if x == m+1 {
			for _, p := range cc.spawn {
				if p.c < x {
					x = p.c
				}
			}
		}
		steps := strings.Repeat("D", cc.minR-1) + "S"
		fmt.Printf("%d %s\n", x, steps)
	}
}
