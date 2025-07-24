package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grid[i])
	}

	// prefix sums for rows
	badRow := make([][]int, n+1)
	medRow := make([][]int, n+1)
	for i := range badRow {
		badRow[i] = make([]int, m+1)
		medRow[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			badRow[i][j] = badRow[i][j-1]
			medRow[i][j] = medRow[i][j-1]
			ch := grid[i-1][j-1]
			if ch == '#' {
				badRow[i][j]++
			} else if ch == 'm' {
				medRow[i][j]++
			}
		}
	}

	// prefix sums for columns and up/down arrays
	medCol := make([][]int, m+1)
	up := make([][]int, n+2)
	down := make([][]int, n+2)
	for j := 0; j <= m; j++ {
		medCol[j] = make([]int, n+1)
	}
	for i := 0; i <= n+1; i++ {
		up[i] = make([]int, m+1)
		down[i] = make([]int, m+1)
	}
	for j := 1; j <= m; j++ {
		for i := 1; i <= n; i++ {
			ch := grid[i-1][j-1]
			medCol[j][i] = medCol[j][i-1]
			if ch == 'm' {
				medCol[j][i]++
			}
			if ch != '#' {
				up[i][j] = up[i-1][j] + 1
			} else {
				up[i][j] = 0
			}
		}
		for i := n; i >= 1; i-- {
			ch := grid[i-1][j-1]
			if ch != '#' {
				down[i][j] = down[i+1][j] + 1
			} else {
				down[i][j] = 0
			}
		}
	}

	ans := 0
	for r := 2; r <= n-1; r++ {
		for c1 := 1; c1 <= m-2; c1++ {
			for c2 := c1 + 2; c2 <= m; c2++ {
				if badRow[r][c2]-badRow[r][c1-1] > 0 {
					continue
				}
				// number of medium cells on horizontal line
				mhor := medRow[r][c2] - medRow[r][c1-1]
				inter := 0
				if grid[r-1][c1-1] == 'm' {
					inter++
				}
				if grid[r-1][c2-1] == 'm' {
					inter++
				}
				k := 1 - (mhor - inter)
				if k < 0 {
					continue
				}
				top := min(up[r][c1], up[r][c2]) - 1
				bottom := min(down[r][c1], down[r][c2]) - 1
				if top < 1 || bottom < 1 {
					continue
				}
				start := r - top
				end := r + bottom
				for r1 := start; r1 <= r-1; r1++ {
					for r2 := r + 1; r2 <= end; r2++ {
						medV := (medCol[c1][r2] - medCol[c1][r1-1]) + (medCol[c2][r2] - medCol[c2][r1-1])
						if medV <= k {
							area := 2*(r2-r1+1) + (c2 - c1 - 1)
							if area > ans {
								ans = area
							}
						}
					}
				}
			}
		}
	}

	fmt.Fprintln(out, ans)
}
