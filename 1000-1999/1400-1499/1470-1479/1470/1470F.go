package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	x int64
	y int64
}

func area(minX, maxX, minY, maxY int64) int64 {
	if maxX < minX || maxY < minY {
		return 0
	}
	return (maxX - minX) * (maxY - minY)
}

func minArea(points []Point) int64 {
	n := len(points)
	if n == 0 {
		return 0
	}
	sort.Slice(points, func(i, j int) bool {
		if points[i].x == points[j].x {
			return points[i].y < points[j].y
		}
		return points[i].x < points[j].x
	})

	prefMinX := make([]int64, n)
	prefMaxX := make([]int64, n)
	prefMinY := make([]int64, n)
	prefMaxY := make([]int64, n)

	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y
	for i := 0; i < n; i++ {
		if points[i].x < minX {
			minX = points[i].x
		}
		if points[i].x > maxX {
			maxX = points[i].x
		}
		if points[i].y < minY {
			minY = points[i].y
		}
		if points[i].y > maxY {
			maxY = points[i].y
		}
		prefMinX[i] = minX
		prefMaxX[i] = maxX
		prefMinY[i] = minY
		prefMaxY[i] = maxY
	}

	sufMinX := make([]int64, n)
	sufMaxX := make([]int64, n)
	sufMinY := make([]int64, n)
	sufMaxY := make([]int64, n)

	minX, maxX = points[n-1].x, points[n-1].x
	minY, maxY = points[n-1].y, points[n-1].y
	for i := n - 1; i >= 0; i-- {
		if points[i].x < minX {
			minX = points[i].x
		}
		if points[i].x > maxX {
			maxX = points[i].x
		}
		if points[i].y < minY {
			minY = points[i].y
		}
		if points[i].y > maxY {
			maxY = points[i].y
		}
		sufMinX[i] = minX
		sufMaxX[i] = maxX
		sufMinY[i] = minY
		sufMaxY[i] = maxY
	}

	ans := area(prefMinX[n-1], prefMaxX[n-1], prefMinY[n-1], prefMaxY[n-1])
	for i := 0; i < n-1; i++ {
		a1 := area(prefMinX[i], prefMaxX[i], prefMinY[i], prefMaxY[i])
		a2 := area(sufMinX[i+1], sufMaxX[i+1], sufMinY[i+1], sufMaxY[i+1])
		if a1+a2 < ans {
			ans = a1 + a2
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		points := make([]Point, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &points[i].x, &points[i].y)
		}

		if n <= 1 {
			fmt.Fprintln(writer, 0)
			continue
		}

		// Vertical splits
		ptsX := make([]Point, n)
		copy(ptsX, points)
		ans := minArea(ptsX)

		// Horizontal splits (swap coordinates)
		ptsY := make([]Point, n)
		for i := 0; i < n; i++ {
			ptsY[i] = Point{x: points[i].y, y: points[i].x}
		}
		val := minArea(ptsY)
		if val < ans {
			ans = val
		}

		fmt.Fprintln(writer, ans)
	}
}
