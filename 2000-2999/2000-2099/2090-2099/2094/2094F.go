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
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		cnt := (n * m) / k

		black := make([][2]int, 0)
		white := make([][2]int, 0)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if (i+j)%2 == 0 {
					black = append(black, [2]int{i, j})
				} else {
					white = append(white, [2]int{i, j})
				}
			}
		}

		grid := make([][]int, n)
		for i := range grid {
			grid[i] = make([]int, m)
		}

		posBlack := 0
		posWhite := 0

		for color := 1; color <= k; color++ {
			if len(black)-posBlack >= cnt {
				for t := 0; t < cnt; t++ {
					cell := black[posBlack]
					posBlack++
					grid[cell[0]][cell[1]] = color
				}
			} else {
				for t := 0; t < cnt; t++ {
					cell := white[posWhite]
					posWhite++
					grid[cell[0]][cell[1]] = color
				}
			}
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, grid[i][j])
			}
			fmt.Fprintln(out)
		}
	}
}
