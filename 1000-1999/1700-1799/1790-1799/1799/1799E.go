package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ r, c int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}
		solve(n, m, grid, out)
	}
}

func solve(n, m int, grid [][]byte, out *bufio.Writer) {
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	comp := make([][]Point, 0, 2)
	vis := make([][]bool, n)
	for i := 0; i < n; i++ {
		vis[i] = make([]bool, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' && !vis[i][j] {
				q := []Point{{i, j}}
				vis[i][j] = true
				cells := []Point{{i, j}}
				for len(q) > 0 {
					p := q[0]
					q = q[1:]
					for _, d := range dirs {
						ni, nj := p.r+d[0], p.c+d[1]
						if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '#' && !vis[ni][nj] {
							vis[ni][nj] = true
							q = append(q, Point{ni, nj})
							cells = append(cells, Point{ni, nj})
						}
					}
				}
				comp = append(comp, cells)
			}
		}
	}
	if len(comp) != 2 {
		// nothing to do
		for i := 0; i < n; i++ {
			fmt.Fprintln(out, string(grid[i]))
		}
		return
	}
	bestCount := n*m + 1
	var bestPath []Point
	// iterate over cell pairs
	for _, a := range comp[0] {
		for _, b := range comp[1] {
			// path horizontal then vertical
			path := make([]Point, 0)
			// horizontal part
			dc := 1
			if b.c < a.c {
				dc = -1
			}
			for j := a.c; j != b.c; j += dc {
				path = append(path, Point{a.r, j})
			}
			path = append(path, Point{a.r, b.c})
			dr := 1
			if b.r < a.r {
				dr = -1
			}
			for i := a.r; i != b.r; i += dr {
				if i == a.r && len(path) > 0 && path[len(path)-1].c == b.c {
					continue
				}
				path = append(path, Point{i, b.c})
			}
			if b.r != a.r {
				path = append(path, Point{b.r, b.c})
			}
			cnt := 0
			mp := make(map[Point]struct{})
			for _, p := range path {
				if grid[p.r][p.c] != '#' {
					if _, ok := mp[p]; !ok {
						cnt++
						mp[p] = struct{}{}
					}
				}
			}
			if cnt < bestCount {
				bestCount = cnt
				bestPath = path
			}
			// path vertical then horizontal
			path = make([]Point, 0)
			dr = 1
			if b.r < a.r {
				dr = -1
			}
			for i := a.r; i != b.r; i += dr {
				path = append(path, Point{i, a.c})
			}
			path = append(path, Point{b.r, a.c})
			dc = 1
			if b.c < a.c {
				dc = -1
			}
			for j := a.c; j != b.c; j += dc {
				if j == a.c && len(path) > 0 && path[len(path)-1].r == b.r {
					continue
				}
				path = append(path, Point{b.r, j})
			}
			if b.c != a.c {
				path = append(path, Point{b.r, b.c})
			}
			cnt = 0
			mp = make(map[Point]struct{})
			for _, p := range path {
				if grid[p.r][p.c] != '#' {
					if _, ok := mp[p]; !ok {
						cnt++
						mp[p] = struct{}{}
					}
				}
			}
			if cnt < bestCount {
				bestCount = cnt
				bestPath = path
			}
		}
	}
	for _, p := range bestPath {
		grid[p.r][p.c] = '#'
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, string(grid[i]))
	}
}
