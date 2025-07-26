package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}

	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	var sizes []int

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			// BFS from (i,j)
			q := make([][2]int, 0, 1)
			q = append(q, [2]int{i, j})
			visited[i][j] = true
			size := 0
			for len(q) > 0 {
				cellPos := q[0]
				q = q[1:]
				r, c := cellPos[0], cellPos[1]
				size++
				cell := grid[r][c]
				for d := 0; d < 4; d++ {
					if cell&(1<<(3-d)) == 0 { // no wall in this direction
						nr, nc := r+dirs[d][0], c+dirs[d][1]
						if nr >= 0 && nr < n && nc >= 0 && nc < m && !visited[nr][nc] {
							visited[nr][nc] = true
							q = append(q, [2]int{nr, nc})
						}
					}
				}
			}
			sizes = append(sizes, size)
		}
	}

	sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })

	out := bufio.NewWriter(os.Stdout)
	for i, v := range sizes {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
	out.Flush()
}
