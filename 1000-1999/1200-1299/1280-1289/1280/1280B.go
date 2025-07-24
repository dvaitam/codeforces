package main

import (
	"bufio"
	"fmt"
	"os"
)

func rowAllA(row string) bool {
	for i := 0; i < len(row); i++ {
		if row[i] != 'A' {
			return false
		}
	}
	return true
}

func colAllA(grid []string, col int) bool {
	for i := 0; i < len(grid); i++ {
		if grid[i][col] != 'A' {
			return false
		}
	}
	return true
}

func solve(grid []string) int {
	r := len(grid)
	c := len(grid[0])
	hasA := false
	allA := true
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if grid[i][j] == 'A' {
				hasA = true
			} else {
				allA = false
			}
		}
	}
	if !hasA {
		return -1 // MORTAL
	}
	if allA {
		return 0
	}

	if rowAllA(grid[0]) || rowAllA(grid[r-1]) || colAllA(grid, 0) || colAllA(grid, c-1) {
		return 1
	}

	if grid[0][0] == 'A' || grid[0][c-1] == 'A' || grid[r-1][0] == 'A' || grid[r-1][c-1] == 'A' {
		return 2
	}

	for i := 0; i < r; i++ {
		if rowAllA(grid[i]) {
			return 2
		}
	}
	for j := 0; j < c; j++ {
		if colAllA(grid, j) {
			return 2
		}
	}

	for j := 0; j < c; j++ {
		if grid[0][j] == 'A' || grid[r-1][j] == 'A' {
			return 3
		}
	}
	for i := 0; i < r; i++ {
		if grid[i][0] == 'A' || grid[i][c-1] == 'A' {
			return 3
		}
	}

	return 4
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var r, c int
		fmt.Fscan(reader, &r, &c)
		grid := make([]string, r)
		for i := 0; i < r; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		res := solve(grid)
		if res == -1 {
			fmt.Fprintln(writer, "MORTAL")
		} else {
			fmt.Fprintln(writer, res)
		}
	}
}
