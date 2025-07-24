package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemA.txt of contest 1292.
// Implements a dynamic check for path existence in a 2xN grid
// with toggling obstacles.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	grid := make([][]bool, 2)
	grid[0] = make([]bool, n)
	grid[1] = make([]bool, n)
	bad := 0

	for ; q > 0; q-- {
		var r, c int
		fmt.Fscan(in, &r, &c)
		r--
		c--
		if grid[r][c] {
			for dc := -1; dc <= 1; dc++ {
				nc := c + dc
				if nc >= 0 && nc < n && grid[1-r][nc] {
					bad--
				}
			}
			grid[r][c] = false
		} else {
			for dc := -1; dc <= 1; dc++ {
				nc := c + dc
				if nc >= 0 && nc < n && grid[1-r][nc] {
					bad++
				}
			}
			grid[r][c] = true
		}
		if bad == 0 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
