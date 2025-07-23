package main

import (
	"bufio"
	"fmt"
	"os"
)

// pair represents a cell position in the grid.
type pair struct {
	x int
	y int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}

	comp := make([][]int, n)
	for i := 0; i < n; i++ {
		comp[i] = make([]int, m)
	}

	// stores number of pictures for each component (1-based index)
	pictures := []int{0}
	compID := 0

	dirs := []int{-1, 0, 1, 0, -1}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '.' || comp[i][j] != 0 {
				continue
			}
			compID++
			queue := []pair{{i, j}}
			comp[i][j] = compID
			pics := 0
			for head := 0; head < len(queue); head++ {
				cur := queue[head]
				x, y := cur.x, cur.y
				for d := 0; d < 4; d++ {
					nx := x + dirs[d]
					ny := y + dirs[d+1]
					if nx < 0 || nx >= n || ny < 0 || ny >= m {
						continue
					}
					if grid[nx][ny] == '*' {
						pics++
					} else if comp[nx][ny] == 0 {
						comp[nx][ny] = compID
						queue = append(queue, pair{nx, ny})
					}
				}
			}
			pictures = append(pictures, pics)
		}
	}

	for q := 0; q < k; q++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		id := comp[x][y]
		fmt.Fprintln(writer, pictures[id])
	}
}
