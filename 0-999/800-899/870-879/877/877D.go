package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pos struct {
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
		var line string
		fmt.Fscan(reader, &line)
		grid[i] = []byte(line)
	}
	var x1, y1, x2, y2 int
	fmt.Fscan(reader, &x1, &y1, &x2, &y2)
	x1--
	y1--
	x2--
	y2--

	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = -1
		}
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	queue := make([]Pos, 0)
	queue = append(queue, Pos{x1, y1})
	dist[x1][y1] = 0
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		d := dist[cur.x][cur.y]
		if cur.x == x2 && cur.y == y2 {
			break
		}
		for _, dir := range dirs {
			nx := cur.x
			ny := cur.y
			for step := 1; step <= k; step++ {
				nx += dir[0]
				ny += dir[1]
				if nx < 0 || nx >= n || ny < 0 || ny >= m {
					break
				}
				if grid[nx][ny] == '#' {
					break
				}
				if dist[nx][ny] != -1 {
					if dist[nx][ny] < d+1 {
						break
					}
					if dist[nx][ny] == d+1 {
						continue
					}
				}
				dist[nx][ny] = d + 1
				queue = append(queue, Pos{nx, ny})
			}
		}
	}

	fmt.Fprintln(writer, dist[x2][y2])
}
