package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ r, c int }

func check(grid [][]byte, n, m int) bool {
	vis := make([][]int, n)
	for i := range vis {
		vis[i] = make([]int, m)
	}
	dirs4 := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	dirs8 := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	id := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' && vis[i][j] == 0 {
				id++
				q := []pair{{i, j}}
				vis[i][j] = id
				for k := 0; k < len(q); k++ {
					p := q[k]
					for _, d := range dirs4 {
						nr, nc := p.r+d[0], p.c+d[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '*' && vis[nr][nc] == 0 {
							vis[nr][nc] = id
							q = append(q, pair{nr, nc})
						}
					}
				}
				if len(q) != 3 {
					return false
				}
				minr, maxr := q[0].r, q[0].r
				minc, maxc := q[0].c, q[0].c
				for _, p := range q[1:] {
					if p.r < minr {
						minr = p.r
					}
					if p.r > maxr {
						maxr = p.r
					}
					if p.c < minc {
						minc = p.c
					}
					if p.c > maxc {
						maxc = p.c
					}
				}
				if maxr-minr != 1 || maxc-minc != 1 {
					return false
				}
				count := 0
				for r := minr; r <= maxr; r++ {
					for c := minc; c <= maxc; c++ {
						if grid[r][c] == '*' {
							if vis[r][c] != id {
								return false
							}
							count++
						}
					}
				}
				if count != 3 {
					return false
				}
				for _, p := range q {
					for _, d := range dirs8 {
						nr, nc := p.r+d[0], p.c+d[1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '*' && vis[nr][nc] != id {
							return false
						}
					}
				}
			}
		}
	}
	return true
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
		if check(grid, n, m) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
