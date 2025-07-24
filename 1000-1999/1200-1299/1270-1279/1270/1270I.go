package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cell struct {
	x int
	y int
}

func apply(grid [][]uint64, fig []Cell, x, y, n int, p uint64) {
	for _, c := range fig {
		i := (x + c.x) & (n - 1)
		j := (y + c.y) & (n - 1)
		grid[i][j] ^= p
	}
}

func allZero(grid [][]uint64) bool {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			if grid[i][j] != 0 {
				return false
			}
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	n := 1 << uint(k)
	grid := make([][]uint64, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]uint64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}
	var t int
	fmt.Fscan(in, &t)
	fig := make([]Cell, t)
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &fig[i].x, &fig[i].y)
		fig[i].x--
		fig[i].y--
	}
	if t == 0 {
		fmt.Println(-1)
		return
	}
	base := fig[0]
	ops := 0
	for iter := 0; iter < n*n && !allZero(grid); iter++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if grid[i][j] != 0 {
					p := grid[i][j]
					tx := (i - base.x + n) & (n - 1)
					ty := (j - base.y + n) & (n - 1)
					apply(grid, fig, tx, ty, n, p)
					ops++
				}
			}
		}
	}
	if !allZero(grid) {
		fmt.Println(-1)
	} else {
		fmt.Println(ops)
	}
}
