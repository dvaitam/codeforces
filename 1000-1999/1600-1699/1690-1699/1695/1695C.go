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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &grid[i][j])
			}
		}

		if (n+m-1)%2 == 1 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		dpMin := make([][]int, n)
		dpMax := make([][]int, n)
		for i := 0; i < n; i++ {
			dpMin[i] = make([]int, m)
			dpMax[i] = make([]int, m)
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				val := grid[i][j]
				if i == 0 && j == 0 {
					dpMin[i][j] = val
					dpMax[i][j] = val
				} else if i == 0 {
					dpMin[i][j] = dpMin[i][j-1] + val
					dpMax[i][j] = dpMax[i][j-1] + val
				} else if j == 0 {
					dpMin[i][j] = dpMin[i-1][j] + val
					dpMax[i][j] = dpMax[i-1][j] + val
				} else {
					minPrev := dpMin[i-1][j]
					if dpMin[i][j-1] < minPrev {
						minPrev = dpMin[i][j-1]
					}
					maxPrev := dpMax[i-1][j]
					if dpMax[i][j-1] > maxPrev {
						maxPrev = dpMax[i][j-1]
					}
					dpMin[i][j] = minPrev + val
					dpMax[i][j] = maxPrev + val
				}
			}
		}
		if dpMin[n-1][m-1] <= 0 && dpMax[n-1][m-1] >= 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
