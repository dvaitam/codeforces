package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	var sx, sy, ex, ey int
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
		for j := 0; j < m; j++ {
			if grid[i][j] == 'E' {
				sx, sy = i, j
			} else if grid[i][j] == 'T' {
				ex, ey = i, j
			}
		}
	}

	type pt struct{ x, y int }
	q := []pt{{sx, sy}}
	dist := make([][]int, n)
	prev := make([][]pt, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, m)
		prev[i] = make([]pt, m)
		for j := 0; j < m; j++ {
			dist[i][j] = -1
		}
	}
	dist[sx][sy] = 0
	dirs := [][3]int{{-1, 0, 'N'}, {1, 0, 'S'}, {0, -1, 'W'}, {0, 1, 'E'}}

	for head := 0; head < len(q); head++ {
		v := q[head]
		if v.x == ex && v.y == ey {
			break
		}
		for _, d := range dirs {
			nx, ny := v.x+d[0], v.y+d[1]
			if nx >= 0 && nx < n && ny >= 0 && ny < m {
				if grid[nx][ny] != 'X' && dist[nx][ny] == -1 {
					dist[nx][ny] = dist[v.x][v.y] + 1
					prev[nx][ny] = v
					q = append(q, pt{nx, ny})
				}
			}
		}
	}

	if dist[ex][ey] == -1 {
		fmt.Fprintln(out, "No solution")
		return
	}
	// reconstruct path
	res := make([]byte, dist[ex][ey])
	cur := pt{ex, ey}
	for dist[cur.x][cur.y] > 0 {
		p := prev[cur.x][cur.y]
		dx := cur.x - p.x
		dy := cur.y - p.y
		var ch byte
		switch {
		case dx == -1 && dy == 0:
			ch = 'N'
		case dx == 1 && dy == 0:
			ch = 'S'
		case dx == 0 && dy == -1:
			ch = 'W'
		case dx == 0 && dy == 1:
			ch = 'E'
		}
		res[dist[cur.x][cur.y]-1] = ch
		cur = p
	}
	fmt.Fprintln(out, string(res))
}
