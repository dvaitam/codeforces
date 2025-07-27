package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var H, W int
	if _, err := fmt.Fscan(in, &H, &W); err != nil {
		return
	}
	grid := make([][]byte, H)
	for i := 0; i < H; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}

	r, c := 0, 0
	berries := 0
	if grid[r][c] == '*' {
		berries++
		grid[r][c] = '.'
	}

	for {
		targetI, targetJ := -1, -1
		bestDist := 1 << 30
		for i := r; i < H; i++ {
			for j := c; j < W; j++ {
				if grid[i][j] == '*' {
					d := (i - r) + (j - c)
					if d < bestDist || (d == bestDist && (i < targetI || (i == targetI && j < targetJ))) {
						bestDist = d
						targetI, targetJ = i, j
					}
				}
			}
		}
		if targetI == -1 {
			break
		}
		for r < targetI {
			r++
			if grid[r][c] == '*' {
				berries++
				grid[r][c] = '.'
			}
		}
		for c < targetJ {
			c++
			if grid[r][c] == '*' {
				berries++
				grid[r][c] = '.'
			}
		}
		if grid[r][c] == '*' {
			berries++
			grid[r][c] = '.'
		}
	}

	for r < H-1 {
		r++
		if grid[r][c] == '*' {
			berries++
			grid[r][c] = '.'
		}
	}
	for c < W-1 {
		c++
		if grid[r][c] == '*' {
			berries++
			grid[r][c] = '.'
		}
	}

	fmt.Println(berries)
}
