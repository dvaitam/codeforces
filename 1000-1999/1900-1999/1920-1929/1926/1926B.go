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
		var n int
		fmt.Fscan(in, &n)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}

		minR, maxR := n, -1
		minC, maxC := n, -1
		ones := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if grid[i][j] == '1' {
					ones++
					if i < minR {
						minR = i
					}
					if i > maxR {
						maxR = i
					}
					if j < minC {
						minC = j
					}
					if j > maxC {
						maxC = j
					}
				}
			}
		}

		height := maxR - minR + 1
		width := maxC - minC + 1
		if height == width {
			ok := true
			for i := minR; i <= maxR && ok; i++ {
				for j := minC; j <= maxC; j++ {
					if grid[i][j] != '1' {
						ok = false
						break
					}
				}
			}
			if ok && ones == height*width {
				fmt.Fprintln(out, "SQUARE")
				continue
			}
		}
		fmt.Fprintln(out, "TRIANGLE")
	}
}
