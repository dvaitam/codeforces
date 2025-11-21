package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type cell struct {
	x int
	y int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for tc := 0; tc < t; tc++ {
		var n, m int
		fmt.Fscan(in, &n, &m)

		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(in, &row[j])
			}
			grid[i] = row
		}

		moves := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &moves[i])
		}

		totalTimes := 2*n - 1
		vals := make([]int, 0, n*totalTimes)

		for l := 1; l <= n; l++ {
			path := moves[l-1]
			bodyLen := l
			body := make([]cell, bodyLen)
			for idx := 0; idx < bodyLen; idx++ {
				body[idx] = cell{0, idx}
			}
			headIdx := 0

			for time := 1; time <= totalTimes; time++ {
				maxVal := 0
				for k := 0; k < bodyLen; k++ {
					idx := headIdx + k
					if idx >= bodyLen {
						idx -= bodyLen
					}
					pos := body[idx]
					v := grid[pos.x][pos.y]
					if v > maxVal {
						maxVal = v
					}
				}
				vals = append(vals, maxVal)
				if time == totalTimes {
					break
				}

				move := path[time-1]
				head := body[headIdx]
				var newHead cell
				if move == 'D' {
					newHead = cell{head.x + 1, head.y}
				} else {
					newHead = cell{head.x, head.y + 1}
				}

				headIdx--
				if headIdx < 0 {
					headIdx += bodyLen
				}
				body[headIdx] = newHead
			}
		}

		sort.Ints(vals)
		for i := 0; i < m; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, vals[i])
		}
		if tc+1 < t {
			fmt.Fprint(out, "\n")
		}
	}
}
