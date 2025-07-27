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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		possible := true
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'B' {
					for _, d := range dirs {
						nx, ny := i+d[0], j+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < m {
							switch grid[nx][ny] {
							case 'G':
								possible = false
							case '.':
								grid[nx][ny] = '#'
							}
						}
					}
				}
			}
		}

		if !possible {
			fmt.Fprintln(out, "No")
			continue
		}

		visited := make([][]bool, n)
		for i := range visited {
			visited[i] = make([]bool, m)
		}
		queue := make([][2]int, 0)
		if grid[n-1][m-1] != '#' {
			queue = append(queue, [2]int{n - 1, m - 1})
			visited[n-1][m-1] = true
		}
		for head := 0; head < len(queue); head++ {
			x, y := queue[head][0], queue[head][1]
			for _, d := range dirs {
				nx, ny := x+d[0], y+d[1]
				if nx >= 0 && nx < n && ny >= 0 && ny < m && !visited[nx][ny] && grid[nx][ny] != '#' && grid[nx][ny] != 'B' {
					visited[nx][ny] = true
					queue = append(queue, [2]int{nx, ny})
				}
			}
		}

		ok := true
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'G' && !visited[i][j] {
					ok = false
				}
				if grid[i][j] == 'B' && visited[i][j] {
					ok = false
				}
			}
		}

		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
