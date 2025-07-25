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
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		compID := make([][]int, n)
		for i := range compID {
			compID[i] = make([]int, m)
			for j := range compID[i] {
				compID[i][j] = -1
			}
		}
		sizes := []int{}
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '#' && compID[i][j] == -1 {
					idx := len(sizes)
					queue := [][2]int{{i, j}}
					compID[i][j] = idx
					sz := 0
					for len(queue) > 0 {
						x, y := queue[0][0], queue[0][1]
						queue = queue[1:]
						sz++
						for _, d := range dirs {
							nx, ny := x+d[0], y+d[1]
							if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] == '#' && compID[nx][ny] == -1 {
								compID[nx][ny] = idx
								queue = append(queue, [2]int{nx, ny})
							}
						}
					}
					sizes = append(sizes, sz)
				}
			}
		}
		maxComp := 0
		for _, s := range sizes {
			if s > maxComp {
				maxComp = s
			}
		}

		mark := make([]int, len(sizes))
		cur := 1
		ans := maxComp

		// Row operations
		for r := 0; r < n; r++ {
			cur++
			rowSize := 0
			for c := 0; c < m; c++ {
				if grid[r][c] == '.' {
					rowSize++
				} else {
					id := compID[r][c]
					if mark[id] != cur {
						mark[id] = cur
						rowSize += sizes[id]
					}
				}
				if r > 0 && grid[r-1][c] == '#' {
					id := compID[r-1][c]
					if mark[id] != cur {
						mark[id] = cur
						rowSize += sizes[id]
					}
				}
				if r+1 < n && grid[r+1][c] == '#' {
					id := compID[r+1][c]
					if mark[id] != cur {
						mark[id] = cur
						rowSize += sizes[id]
					}
				}
			}
			if rowSize > ans {
				ans = rowSize
			}
		}

		// Column operations
		for c := 0; c < m; c++ {
			cur++
			colSize := 0
			for r := 0; r < n; r++ {
				if grid[r][c] == '.' {
					colSize++
				} else {
					id := compID[r][c]
					if mark[id] != cur {
						mark[id] = cur
						colSize += sizes[id]
					}
				}
				if c > 0 && grid[r][c-1] == '#' {
					id := compID[r][c-1]
					if mark[id] != cur {
						mark[id] = cur
						colSize += sizes[id]
					}
				}
				if c+1 < m && grid[r][c+1] == '#' {
					id := compID[r][c+1]
					if mark[id] != cur {
						mark[id] = cur
						colSize += sizes[id]
					}
				}
			}
			if colSize > ans {
				ans = colSize
			}
		}

		fmt.Fprintln(out, ans)
	}
}
