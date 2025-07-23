package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(grid [][]int) bool {
	N := len(grid)
	A := make([][]int, N+1)
	for i := 0; i <= N; i++ {
		A[i] = make([]int, N+1)
	}
	for j := 1; j <= N; j++ {
		for i := 1; i <= N; i++ {
			val := grid[j-1][i-1] - A[i-1][j-1] - A[i][j-1] - A[i-1][j]
			if val != 0 && val != 1 {
				return false
			}
			A[i][j] = val
		}
	}
	minX, minY, maxX, maxY := N+1, N+1, -1, -1
	for i := 0; i <= N; i++ {
		for j := 0; j <= N; j++ {
			if A[i][j] == 1 {
				if i < minX {
					minX = i
				}
				if i > maxX {
					maxX = i
				}
				if j < minY {
					minY = j
				}
				if j > maxY {
					maxY = j
				}
			}
		}
	}
	if maxX < minX || maxY < minY {
		return false
	}
	if maxX == minX || maxY == minY {
		return false
	}
	for i := 0; i <= N; i++ {
		for j := 0; j <= N; j++ {
			exp := 0
			if i >= minX && i <= maxX && j >= minY && j <= maxY {
				exp = 1
			}
			if A[i][j] != exp {
				return false
			}
		}
	}
	return true
}
func main() {
	in := bufio.NewReader(os.Stdin)
	var N int
	if _, err := fmt.Fscan(in, &N); err != nil {
		return
	}
	grid := make([][]int, N)
	for row := 0; row < N; row++ {
		var line string
		fmt.Fscan(in, &line)
		y := N - 1 - row
		grid[y] = make([]int, N)
		for x := 0; x < N; x++ {
			grid[y][x] = int(line[x] - '0')
		}
	}
	if check(grid) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
