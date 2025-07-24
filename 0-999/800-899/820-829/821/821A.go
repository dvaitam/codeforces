package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program determines whether a given n x n grid is "good" according to the
// rules described in problemA.txt. Every cell value not equal to 1 must be
// representable as the sum of some value in the same row and some value in the
// same column.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &grid[i][j])
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				continue
			}
			ok := false
			for r := 0; r < n && !ok; r++ {
				for c := 0; c < n; c++ {
					if grid[i][r]+grid[c][j] == grid[i][j] {
						ok = true
						break
					}
				}
			}
			if !ok {
				fmt.Fprintln(writer, "No")
				return
			}
		}
	}

	fmt.Fprintln(writer, "Yes")
}
