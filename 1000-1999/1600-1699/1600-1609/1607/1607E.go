package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s string
		fmt.Fscan(reader, &s)

		minX, maxX := 0, 0
		minY, maxY := 0, 0
		x, y := 0, 0
		startRow, startCol := 1, 1

		for _, c := range s {
			nx, ny := x, y
			switch c {
			case 'L':
				ny--
			case 'R':
				ny++
			case 'U':
				nx--
			case 'D':
				nx++
			}
			newMinX := min(minX, nx)
			newMaxX := max(maxX, nx)
			newMinY := min(minY, ny)
			newMaxY := max(maxY, ny)
			if newMaxX-newMinX+1 > n || newMaxY-newMinY+1 > m {
				break
			}
			x, y = nx, ny
			minX, maxX = newMinX, newMaxX
			minY, maxY = newMinY, newMaxY
			startRow = 1 - minX
			startCol = 1 - minY
		}
		fmt.Fprintln(writer, startRow, startCol)
	}
}
