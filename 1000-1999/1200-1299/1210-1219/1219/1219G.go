package main

import (
	"bufio"
	"fmt"
	"os"
)

func updateTop3(best *[3]int64, val int64) {
	if val > best[0] {
		best[2] = best[1]
		best[1] = best[0]
		best[0] = val
	} else if val > best[1] {
		best[2] = best[1]
		best[1] = val
	} else if val > best[2] {
		best[2] = val
	}
}

func updateTop2(best *[2]int64, val int64) {
	if val > best[0] {
		best[1] = best[0]
		best[0] = val
	} else if val > best[1] {
		best[1] = val
	}
}

func top4Sum(arr []int64) int64 {
	best := [4]int64{}
	for _, v := range arr {
		if v > best[0] {
			best[3] = best[2]
			best[2] = best[1]
			best[1] = best[0]
			best[0] = v
		} else if v > best[1] {
			best[3] = best[2]
			best[2] = best[1]
			best[1] = v
		} else if v > best[2] {
			best[3] = best[2]
			best[2] = v
		} else if v > best[3] {
			best[3] = v
		}
	}
	limit := 4
	if len(arr) < 4 {
		limit = len(arr)
	}
	sum := int64(0)
	for i := 0; i < limit; i++ {
		sum += best[i]
	}
	return sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	grid := make([][]int64, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int64, m)
	}
	rowSum := make([]int64, n)
	colSum := make([]int64, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var x int64
			fmt.Fscan(in, &x)
			grid[i][j] = x
			rowSum[i] += x
			colSum[j] += x
		}
	}

	var best int64

	// 4 rows
	if n > 0 {
		v := top4Sum(rowSum)
		if v > best {
			best = v
		}
	}
	// 4 columns
	if m > 0 {
		v := top4Sum(colSum)
		if v > best {
			best = v
		}
	}

	// 3 rows + 1 column
	if n >= 3 {
		for c := 0; c < m; c++ {
			var bestRows [3]int64
			for r := 0; r < n; r++ {
				val := rowSum[r] - grid[r][c]
				updateTop3(&bestRows, val)
			}
			v := colSum[c] + bestRows[0] + bestRows[1] + bestRows[2]
			if v > best {
				best = v
			}
		}
	}

	// 1 row + 3 columns
	if m >= 3 {
		for r := 0; r < n; r++ {
			var bestCols [3]int64
			for c := 0; c < m; c++ {
				val := colSum[c] - grid[r][c]
				updateTop3(&bestCols, val)
			}
			v := rowSum[r] + bestCols[0] + bestCols[1] + bestCols[2]
			if v > best {
				best = v
			}
		}
	}

	// 2 rows + 2 columns
	if n >= 2 && m >= 2 {
		if n <= m {
			for r1 := 0; r1 < n; r1++ {
				for r2 := r1 + 1; r2 < n; r2++ {
					var bestCols [2]int64
					for c := 0; c < m; c++ {
						val := colSum[c] - grid[r1][c] - grid[r2][c]
						updateTop2(&bestCols, val)
					}
					v := rowSum[r1] + rowSum[r2] + bestCols[0] + bestCols[1]
					if v > best {
						best = v
					}
				}
			}
		} else {
			for c1 := 0; c1 < m; c1++ {
				for c2 := c1 + 1; c2 < m; c2++ {
					var bestRows [2]int64
					for r := 0; r < n; r++ {
						val := rowSum[r] - grid[r][c1] - grid[r][c2]
						updateTop2(&bestRows, val)
					}
					v := colSum[c1] + colSum[c2] + bestRows[0] + bestRows[1]
					if v > best {
						best = v
					}
				}
			}
		}
	}

	fmt.Fprintln(out, best)
}
