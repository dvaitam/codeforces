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

	var h, w int
	if _, err := fmt.Fscan(reader, &h, &w); err != nil {
		return
	}
	grid := make([]string, h)
	for i := 0; i < h; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	hPref := make([][]int, h+1)
	vPref := make([][]int, h+1)
	for i := 0; i <= h; i++ {
		hPref[i] = make([]int, w+1)
		vPref[i] = make([]int, w+1)
	}

	for i := 1; i <= h; i++ {
		for j := 1; j <= w; j++ {
			if j < w && grid[i-1][j-1] == '.' && grid[i-1][j] == '.' {
				hPref[i][j] = 1
			}
			if i < h && grid[i-1][j-1] == '.' && grid[i][j-1] == '.' {
				vPref[i][j] = 1
			}
			hPref[i][j] += hPref[i-1][j] + hPref[i][j-1] - hPref[i-1][j-1]
			vPref[i][j] += vPref[i-1][j] + vPref[i][j-1] - vPref[i-1][j-1]
		}
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var r1, c1, r2, c2 int
		fmt.Fscan(reader, &r1, &c1, &r2, &c2)
		horizontal := hPref[r2][c2-1] - hPref[r1-1][c2-1] - hPref[r2][c1-1] + hPref[r1-1][c1-1]
		vertical := vPref[r2-1][c2] - vPref[r1-1][c2] - vPref[r2-1][c1-1] + vPref[r1-1][c1-1]
		fmt.Fprintln(writer, horizontal+vertical)
	}
}
