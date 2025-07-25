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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		// prefix sums per row for quick segment queries
		pref := make([][]int, n)
		for i := 0; i < n; i++ {
			pref[i] = make([]int, m+1)
			for j := 0; j < m; j++ {
				val := 0
				if grid[i][j] == '#' {
					val = 1
				}
				pref[i][j+1] = pref[i][j] + val
			}
		}

		getRow := func(r, l, r2 int) int {
			if l < 0 {
				l = 0
			}
			if r2 >= m {
				r2 = m - 1
			}
			if l > r2 || r < 0 || r >= n {
				return 0
			}
			return pref[r][r2+1] - pref[r][l]
		}

		best := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				// right-down
				cnt := 0
				for d := 0; d <= k; d++ {
					x := i + d
					if x >= n {
						break
					}
					y1 := j
					y2 := j + (k - d)
					if y1 >= m {
						break
					}
					if y2 < 0 {
						continue
					}
					cnt += getRow(x, y1, y2)
				}
				if cnt > best {
					best = cnt
				}

				// left-down
				cnt = 0
				for d := 0; d <= k; d++ {
					x := i + d
					if x >= n {
						break
					}
					y2 := j
					y1 := j - (k - d)
					if y2 < 0 {
						break
					}
					if y1 >= m {
						continue
					}
					cnt += getRow(x, y1, y2)
				}
				if cnt > best {
					best = cnt
				}

				// left-up
				cnt = 0
				for d := 0; d <= k; d++ {
					x := i - d
					if x < 0 {
						break
					}
					y2 := j
					y1 := j - (k - d)
					if y2 < 0 {
						break
					}
					if y1 >= m {
						continue
					}
					cnt += getRow(x, y1, y2)
				}
				if cnt > best {
					best = cnt
				}

				// right-up
				cnt = 0
				for d := 0; d <= k; d++ {
					x := i - d
					if x < 0 {
						break
					}
					y1 := j
					y2 := j + (k - d)
					if y1 >= m {
						break
					}
					if y2 < 0 {
						continue
					}
					cnt += getRow(x, y1, y2)
				}
				if cnt > best {
					best = cnt
				}
			}
		}

		fmt.Fprintln(out, best)
	}
}
