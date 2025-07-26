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

	var n, m, t int
	if _, err := fmt.Fscan(in, &n, &m, &t); err != nil {
		return
	}

	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		row := make([]int, m)
		for j := 0; j < m; j++ {
			if s[j] == '1' {
				row[j] = 1
			}
		}
		grid[i] = row
	}

	type pair struct{ x, y int }
	dist := make([][]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			dist[i][j] = -1
		}
	}

	q := make([]pair, 0)
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == grid[i][j] {
					dist[i][j] = 0
					q = append(q, pair{i, j})
					break
				}
			}
		}
	}

	for head := 0; head < len(q); head++ {
		cur := q[head]
		for _, d := range dirs {
			ni, nj := cur.x+d[0], cur.y+d[1]
			if ni >= 0 && ni < n && nj >= 0 && nj < m && dist[ni][nj] == -1 {
				dist[ni][nj] = dist[cur.x][cur.y] + 1
				q = append(q, pair{ni, nj})
			}
		}
	}

	for ; t > 0; t-- {
		var i, j int
		var p int64
		fmt.Fscan(in, &i, &j, &p)
		i--
		j--
		cell := grid[i][j]
		d := dist[i][j]
		if d == -1 || p <= d {
			fmt.Fprintln(out, cell)
		} else {
			if (p-d)%2 == 1 {
				fmt.Fprintln(out, 1-cell)
			} else {
				fmt.Fprintln(out, cell)
			}
		}
	}
}
