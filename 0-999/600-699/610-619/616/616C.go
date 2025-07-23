package main

import (
	"bufio"
	"fmt"
	"os"
)

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

	id := make([][]int, n)
	for i := range id {
		id[i] = make([]int, m)
	}
	sizes := []int{0}
	type pair struct{ r, c int }
	dirs := []pair{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queue := make([]pair, n*m)

	curID := 1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' && id[i][j] == 0 {
				front, back := 0, 0
				id[i][j] = curID
				queue[back] = pair{i, j}
				back++
				for front < back {
					p := queue[front]
					front++
					for _, d := range dirs {
						nr, nc := p.r+d.r, p.c+d.c
						if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '.' && id[nr][nc] == 0 {
							id[nr][nc] = curID
							queue[back] = pair{nr, nc}
							back++
						}
					}
				}
				sizes = append(sizes, back)
				curID++
			}
		}
	}

	outLines := make([][]byte, n)
	for i := 0; i < n; i++ {
		line := make([]byte, m)
		for j := 0; j < m; j++ {
			if grid[i][j] == '.' {
				line[j] = '.'
			} else {
				sum := 1
				seen := make(map[int]bool, 4)
				for _, d := range dirs {
					nr, nc := i+d.r, j+d.c
					if nr >= 0 && nr < n && nc >= 0 && nc < m && grid[nr][nc] == '.' {
						cid := id[nr][nc]
						if !seen[cid] {
							sum += sizes[cid]
							seen[cid] = true
						}
					}
				}
				line[j] = byte('0' + sum%10)
			}
		}
		outLines[i] = line
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		out.Write(outLines[i])
		out.WriteByte('\n')
	}
	out.Flush()
}
