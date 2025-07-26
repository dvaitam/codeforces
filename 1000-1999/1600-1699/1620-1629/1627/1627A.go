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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, r, c int
		fmt.Fscan(reader, &n, &m, &r, &c)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		r--
		c--
		if grid[r][c] == 'B' {
			fmt.Fprintln(writer, 0)
			continue
		}
		rowBlack := false
		for j := 0; j < m; j++ {
			if grid[r][j] == 'B' {
				rowBlack = true
				break
			}
		}
		colBlack := false
		for i := 0; i < n; i++ {
			if grid[i][c] == 'B' {
				colBlack = true
				break
			}
		}
		if rowBlack || colBlack {
			fmt.Fprintln(writer, 1)
			continue
		}
		anyBlack := false
		for i := 0; i < n && !anyBlack; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'B' {
					anyBlack = true
					break
				}
			}
		}
		if anyBlack {
			fmt.Fprintln(writer, 2)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
