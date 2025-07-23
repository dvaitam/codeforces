package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const mod int64 = 1000000007

func powmod(base int64, exp int) int64 {
	result := int64(1)
	b := base % mod
	e := exp
	for e > 0 {
		if e&1 == 1 {
			result = result * b % mod
		}
		b = b * b % mod
		e >>= 1
	}
	return result
}

func min4(a, b, c, d int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	if d < m {
		m = d
	}
	return m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		row := make([]int, n)
		for j := 0; j < n; j++ {
			row[j] = int(s[j] - '0')
		}
		grid[i] = row
	}
	// prefix sums for counts of 2 and 3
	row2 := make([][]int, n)
	row3 := make([][]int, n)
	for i := 0; i < n; i++ {
		row2[i] = make([]int, n+1)
		row3[i] = make([]int, n+1)
		for j := 0; j < n; j++ {
			val := grid[i][j]
			row2[i][j+1] = row2[i][j]
			row3[i][j+1] = row3[i][j]
			if val == 2 {
				row2[i][j+1]++
			} else if val == 3 {
				row3[i][j+1]++
			}
		}
	}
	col2 := make([][]int, n+1)
	col3 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		col2[i] = make([]int, n)
		col3[i] = make([]int, n)
	}
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			val := grid[i][j]
			col2[i+1][j] = col2[i][j]
			col3[i+1][j] = col3[i][j]
			if val == 2 {
				col2[i+1][j]++
			} else if val == 3 {
				col3[i+1][j]++
			}
		}
	}
	diag12 := make([][]int, n+1)
	diag13 := make([][]int, n+1)
	diag22 := make([][]int, n+1)
	diag23 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		diag12[i] = make([]int, n+1)
		diag13[i] = make([]int, n+1)
		diag22[i] = make([]int, n+1)
		diag23[i] = make([]int, n+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := grid[i][j]
			diag12[i+1][j+1] = diag12[i][j]
			diag13[i+1][j+1] = diag13[i][j]
			diag22[i+1][j] = diag22[i][j+1]
			diag23[i+1][j] = diag23[i][j+1]
			if val == 2 {
				diag12[i+1][j+1]++
				diag22[i+1][j]++
			} else if val == 3 {
				diag13[i+1][j+1]++
				diag23[i+1][j]++
			}
		}
	}
	// lengths of non-zero segments in 8 directions
	left := make([][]int, n)
	right := make([][]int, n)
	up := make([][]int, n)
	down := make([][]int, n)
	ul := make([][]int, n)
	ur := make([][]int, n)
	dl := make([][]int, n)
	dr := make([][]int, n)
	for i := 0; i < n; i++ {
		left[i] = make([]int, n)
		right[i] = make([]int, n)
		up[i] = make([]int, n)
		down[i] = make([]int, n)
		ul[i] = make([]int, n)
		ur[i] = make([]int, n)
		dl[i] = make([]int, n)
		dr[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] != 0 {
				left[i][j] = 1
				up[i][j] = 1
				ul[i][j] = 1
				ur[i][j] = 1
				if j > 0 {
					left[i][j] += left[i][j-1]
				}
				if i > 0 {
					up[i][j] += up[i-1][j]
				}
				if i > 0 && j > 0 {
					ul[i][j] += ul[i-1][j-1]
				}
				if i > 0 && j+1 < n {
					ur[i][j] += ur[i-1][j+1]
				}
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if grid[i][j] != 0 {
				right[i][j] = 1
				down[i][j] = 1
				dr[i][j] = 1
				dl[i][j] = 1
				if j+1 < n {
					right[i][j] += right[i][j+1]
				}
				if i+1 < n {
					down[i][j] += down[i+1][j]
				}
				if i+1 < n && j+1 < n {
					dr[i][j] += dr[i+1][j+1]
				}
				if i+1 < n && j > 0 {
					dl[i][j] += dl[i+1][j-1]
				}
			}
		}
	}

	ln2 := math.Log(2)
	ln3 := math.Log(3)
	bestLog := -1.0
	bestA, bestB := 0, 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 0 {
				continue
			}
			center2 := 0
			center3 := 0
			if grid[i][j] == 2 {
				center2 = 1
			} else if grid[i][j] == 3 {
				center3 = 1
			}
			r := min4(left[i][j], right[i][j], up[i][j], down[i][j]) - 1
			if r >= 0 {
				a := row2[i][j+r+1] - row2[i][j-r]
				a += col2[i+r+1][j] - col2[i-r][j]
				a -= center2
				b := row3[i][j+r+1] - row3[i][j-r]
				b += col3[i+r+1][j] - col3[i-r][j]
				b -= center3
				logv := float64(a)*ln2 + float64(b)*ln3
				if logv > bestLog {
					bestLog = logv
					bestA = a
					bestB = b
				}
			}
			r = min4(ul[i][j], ur[i][j], dl[i][j], dr[i][j]) - 1
			if r >= 0 {
				a := diag12[i+r+1][j+r+1] - diag12[i-r][j-r]
				a += diag22[i+r+1][j-r] - diag22[i-r][j+r+1]
				a -= center2
				b := diag13[i+r+1][j+r+1] - diag13[i-r][j-r]
				b += diag23[i+r+1][j-r] - diag23[i-r][j+r+1]
				b -= center3
				logv := float64(a)*ln2 + float64(b)*ln3
				if logv > bestLog {
					bestLog = logv
					bestA = a
					bestB = b
				}
			}
		}
	}
	if bestLog < 0 {
		fmt.Println(0)
		return
	}
	ans := powmod(2, bestA) * powmod(3, bestB) % mod
	fmt.Println(ans)
}
