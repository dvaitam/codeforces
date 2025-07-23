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

	var d int
	if _, err := fmt.Fscan(reader, &d); err != nil {
		return
	}
	var n, m int
	fmt.Fscan(reader, &n, &m)

	type Sofa struct {
		minX, maxX int
		minY, maxY int
	}
	sofas := make([]Sofa, d)

	minXCount := make([]int, n+2)
	maxXCount := make([]int, n+2)
	minYCount := make([]int, m+2)
	maxYCount := make([]int, m+2)

	for i := 0; i < d; i++ {
		var x1, y1, x2, y2 int
		fmt.Fscan(reader, &x1, &y1, &x2, &y2)
		minX := x1
		if x2 < minX {
			minX = x2
		}
		maxX := x1
		if x2 > maxX {
			maxX = x2
		}
		minY := y1
		if y2 < minY {
			minY = y2
		}
		maxY := y1
		if y2 > maxY {
			maxY = y2
		}
		sofas[i] = Sofa{minX, maxX, minY, maxY}
		minXCount[minX]++
		maxXCount[maxX]++
		minYCount[minY]++
		maxYCount[maxY]++
	}

	var cntl, cntr, cntt, cntb int
	fmt.Fscan(reader, &cntl, &cntr, &cntt, &cntb)

	prefixMinX := make([]int, n+2)
	prefixMaxX := make([]int, n+2)
	for i := 1; i <= n; i++ {
		prefixMinX[i] = prefixMinX[i-1] + minXCount[i]
		prefixMaxX[i] = prefixMaxX[i-1] + maxXCount[i]
	}
	prefixMinY := make([]int, m+2)
	prefixMaxY := make([]int, m+2)
	for i := 1; i <= m; i++ {
		prefixMinY[i] = prefixMinY[i-1] + minYCount[i]
		prefixMaxY[i] = prefixMaxY[i-1] + maxYCount[i]
	}

	total := d
	for idx, s := range sofas {
		left := prefixMinX[s.maxX-1]
		if s.minX < s.maxX {
			left--
		}
		right := total - prefixMaxX[s.minX]
		if s.maxX > s.minX {
			right--
		}
		top := prefixMinY[s.maxY-1]
		if s.minY < s.maxY {
			top--
		}
		bottom := total - prefixMaxY[s.minY]
		if s.maxY > s.minY {
			bottom--
		}
		if left == cntl && right == cntr && top == cntt && bottom == cntb {
			fmt.Fprintln(writer, idx+1)
			return
		}
	}
	fmt.Fprintln(writer, -1)
}
