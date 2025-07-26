package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ x, y int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	dx := []int{1, -1, 0, 0}
	dy := []int{0, 0, 1, -1}

	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &grid[i][j])
			}
		}
		visited := make([][]bool, n)
		for i := range visited {
			visited[i] = make([]bool, m)
		}
		best := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 0 || visited[i][j] {
					continue
				}
				sum := 0
				queue := []pair{{i, j}}
				visited[i][j] = true
				for head := 0; head < len(queue); head++ {
					p := queue[head]
					sum += grid[p.x][p.y]
					for dir := 0; dir < 4; dir++ {
						nx := p.x + dx[dir]
						ny := p.y + dy[dir]
						if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] && grid[nx][ny] > 0 {
							visited[nx][ny] = true
							queue = append(queue, pair{nx, ny})
						}
					}
				}
				if sum > best {
					best = sum
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
