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
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}
		minRow, minCol := n, m
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'R' {
					if i < minRow {
						minRow = i
					}
					if j < minCol {
						minCol = j
					}
				}
			}
		}
		if grid[minRow][minCol] == 'R' {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
