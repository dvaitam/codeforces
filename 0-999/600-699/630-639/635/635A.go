package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var r, c, n, k int
	if _, err := fmt.Fscan(in, &r, &c, &n, &k); err != nil {
		return
	}
	grid := make([][]int, r+1)
	for i := range grid {
		grid[i] = make([]int, c+1)
	}
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x >= 1 && x <= r && y >= 1 && y <= c {
			grid[x][y] = 1
		}
	}

	prefix := make([][]int, r+1)
	for i := range prefix {
		prefix[i] = make([]int, c+1)
	}
	for i := 1; i <= r; i++ {
		row := 0
		for j := 1; j <= c; j++ {
			row += grid[i][j]
			prefix[i][j] = prefix[i-1][j] + row
		}
	}

	ans := 0
	for x1 := 1; x1 <= r; x1++ {
		for x2 := x1; x2 <= r; x2++ {
			for y1 := 1; y1 <= c; y1++ {
				for y2 := y1; y2 <= c; y2++ {
					count := prefix[x2][y2] - prefix[x1-1][y2] - prefix[x2][y1-1] + prefix[x1-1][y1-1]
					if count >= k {
						ans++
					}
				}
			}
		}
	}

	fmt.Println(ans)
}
