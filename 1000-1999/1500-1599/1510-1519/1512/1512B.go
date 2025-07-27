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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		grid := make([][]byte, n)
		var stars [][2]int
		for i := 0; i < n; i++ {
			var row string
			fmt.Fscan(reader, &row)
			grid[i] = []byte(row)
			for j := 0; j < n; j++ {
				if grid[i][j] == '*' {
					stars = append(stars, [2]int{i, j})
				}
			}
		}
		r1, c1 := stars[0][0], stars[0][1]
		r2, c2 := stars[1][0], stars[1][1]
		if r1 == r2 {
			r3 := r1 + 1
			if r3 >= n {
				r3 = r1 - 1
			}
			grid[r3][c1] = '*'
			grid[r3][c2] = '*'
		} else if c1 == c2 {
			c3 := c1 + 1
			if c3 >= n {
				c3 = c1 - 1
			}
			grid[r1][c3] = '*'
			grid[r2][c3] = '*'
		} else {
			grid[r1][c2] = '*'
			grid[r2][c1] = '*'
		}
		for i := 0; i < n; i++ {
			writer.Write(grid[i])
			writer.WriteByte('\n')
		}
	}
}
