package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}

	rowEmpty := make([]bool, n)
	colEmpty := make([]bool, m)

	for i := 0; i < n; i++ {
		first, last := -1, -1
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				if first == -1 {
					first = j
				}
				last = j
			}
		}
		if first == -1 {
			rowEmpty[i] = true
		} else {
			for j := first; j <= last; j++ {
				if grid[i][j] != '#' {
					fmt.Fprintln(writer, -1)
					return
				}
			}
		}
	}

	for j := 0; j < m; j++ {
		first, last := -1, -1
		for i := 0; i < n; i++ {
			if grid[i][j] == '#' {
				if first == -1 {
					first = i
				}
				last = i
			}
		}
		if first == -1 {
			colEmpty[j] = true
		} else {
			for i := first; i <= last; i++ {
				if grid[i][j] != '#' {
					fmt.Fprintln(writer, -1)
					return
				}
			}
		}
	}

	rowZero := false
	colZero := false
	for i := 0; i < n; i++ {
		if rowEmpty[i] {
			rowZero = true
			break
		}
	}
	for j := 0; j < m; j++ {
		if colEmpty[j] {
			colZero = true
			break
		}
	}
	if rowZero != colZero {
		fmt.Fprintln(writer, -1)
		return
	}

	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	components := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' && !visited[i][j] {
				components++
				queue := [][2]int{{i, j}}
				visited[i][j] = true
				for qi := 0; qi < len(queue); qi++ {
					x, y := queue[qi][0], queue[qi][1]
					for _, d := range dirs {
						nx, ny := x+d[0], y+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < m {
							if grid[nx][ny] == '#' && !visited[nx][ny] {
								visited[nx][ny] = true
								queue = append(queue, [2]int{nx, ny})
							}
						}
					}
				}
			}
		}
	}
	fmt.Fprintln(writer, components)
}
