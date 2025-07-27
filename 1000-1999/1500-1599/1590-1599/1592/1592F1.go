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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		row := make([]int, m)
		for j := 0; j < m; j++ {
			if s[j] == 'B' {
				row[j] = 1
			}
		}
		grid[i] = row
	}

	diff := make([][]int, n)
	for i := range diff {
		diff[i] = make([]int, m)
	}
	base := 0
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			v := grid[i][j]
			if i+1 < n {
				v ^= grid[i+1][j]
			}
			if j+1 < m {
				v ^= grid[i][j+1]
			}
			if i+1 < n && j+1 < m {
				v ^= grid[i+1][j+1]
			}
			diff[i][j] = v
			if v == 1 {
				base++
			}
		}
	}

	best := base
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			change := 0
			if x > 0 && y > 0 {
				if diff[x-1][y-1] == 1 {
					change--
				} else {
					change++
				}
			}
			if x > 0 {
				if diff[x-1][m-1] == 1 {
					change--
				} else {
					change++
				}
			}
			if y > 0 {
				if diff[n-1][y-1] == 1 {
					change--
				} else {
					change++
				}
			}
			if diff[n-1][m-1] == 1 {
				change--
			} else {
				change++
			}
			cost := base + change + 3
			if cost < best {
				best = cost
			}
		}
	}

	fmt.Fprintln(out, best)
}
