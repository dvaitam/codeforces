package main

import (
	"bufio"
	"fmt"
	"os"
)

type point struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]byte, n)
		var sx, sy int
		for i := 0; i < n; i++ {
			var line string
			fmt.Fscan(reader, &line)
			grid[i] = []byte(line)
			if idx := indexByte(grid[i], 'L'); idx >= 0 {
				sx, sy = i, idx
			}
		}
		deg := make([][]int, n)
		for i := 0; i < n; i++ {
			deg[i] = make([]int, m)
		}
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '#' {
					continue
				}
				cnt := 0
				for _, d := range dirs {
					ni, nj := i+d[0], j+d[1]
					if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] != '#' {
						cnt++
					}
				}
				deg[i][j] = cnt
			}
		}
		q := []point{{sx, sy}}
		for head := 0; head < len(q); head++ {
			p := q[head]
			for _, d := range dirs {
				ni, nj := p.x+d[0], p.y+d[1]
				if ni < 0 || ni >= n || nj < 0 || nj >= m {
					continue
				}
				if grid[ni][nj] != '.' {
					continue
				}
				deg[ni][nj]--
				if deg[ni][nj] <= 1 {
					grid[ni][nj] = '+'
					q = append(q, point{ni, nj})
				}
			}
		}
		for i := 0; i < n; i++ {
			writer.WriteString(string(grid[i]))
			writer.WriteByte('\n')
		}
	}
}

func indexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}
