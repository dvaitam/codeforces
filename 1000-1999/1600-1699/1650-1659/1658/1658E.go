package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	grid := make([][]int, n)
	maxVal := -1
	maxX, maxY := 0, 0
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
			if grid[i][j] > maxVal {
				maxVal = grid[i][j]
				maxX, maxY = i, j
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == maxX && j == maxY {
				fmt.Fprint(out, "M")
				continue
			}
			d := abs(i-maxX) + abs(j-maxY)
			if d > k {
				fmt.Fprint(out, "G")
			} else {
				fmt.Fprint(out, "D")
			}
		}
		fmt.Fprintln(out)
	}
}
