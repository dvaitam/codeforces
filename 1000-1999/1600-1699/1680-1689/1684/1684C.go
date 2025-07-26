package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &row[j])
			}
			grid[i] = row
		}
		l, r := -1, -1
		for i := 0; i < n; i++ {
			sortedRow := make([]int, m)
			copy(sortedRow, grid[i])
			sort.Ints(sortedRow)
			diff := make([]int, 0)
			for j := 0; j < m; j++ {
				if grid[i][j] != sortedRow[j] {
					diff = append(diff, j)
				}
			}
			if len(diff) > 0 && l == -1 {
				l = diff[0]
				r = diff[len(diff)-1]
			}
		}
		if l == -1 {
			fmt.Fprintln(writer, 1, 1)
			continue
		}
		ok := true
		for i := 0; i < n && ok; i++ {
			grid[i][l], grid[i][r] = grid[i][r], grid[i][l]
			for j := 1; j < m; j++ {
				if grid[i][j] < grid[i][j-1] {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, l+1, r+1)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
