package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type rect struct {
	x1, x2 int
	y1, y2 int
}

func uniqueSorted(vals []int) []int {
	sort.Ints(vals)
	res := vals[:0]
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			res = append(res, v)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		rects := make([]rect, 0, n)
		xVals := make([]int, 0, 2*n)
		yVals := make([]int, 0, 2*n)

		curX, curY := 0, 0
		for i := 0; i < n; i++ {
			var dx, dy int
			fmt.Fscan(in, &dx, &dy)
			curX += dx
			curY += dy
			r := rect{
				x1: curX,
				x2: curX + m,
				y1: curY,
				y2: curY + m,
			}
			rects = append(rects, r)
			xVals = append(xVals, r.x1, r.x2)
			yVals = append(yVals, r.y1, r.y2)
		}

		xVals = uniqueSorted(xVals)
		yVals = uniqueSorted(yVals)

		xIndex := make(map[int]int, len(xVals))
		for idx, v := range xVals {
			xIndex[v] = idx
		}
		yIndex := make(map[int]int, len(yVals))
		for idx, v := range yVals {
			yIndex[v] = idx
		}

		w := len(xVals) - 1
		h := len(yVals) - 1
		grid := make([][]bool, w)
		for i := range grid {
			grid[i] = make([]bool, h)
		}

		for _, r := range rects {
			x0 := xIndex[r.x1]
			x1 := xIndex[r.x2]
			y0 := yIndex[r.y1]
			y1 := yIndex[r.y2]
			for xi := x0; xi < x1; xi++ {
				for yi := y0; yi < y1; yi++ {
					grid[xi][yi] = true
				}
			}
		}

		perimeter := 0
		for xi := 0; xi < w; xi++ {
			width := xVals[xi+1] - xVals[xi]
			for yi := 0; yi < h; yi++ {
				if !grid[xi][yi] {
					continue
				}
				height := yVals[yi+1] - yVals[yi]
				if xi == 0 || !grid[xi-1][yi] {
					perimeter += height
				}
				if xi == w-1 || !grid[xi+1][yi] {
					perimeter += height
				}
				if yi == 0 || !grid[xi][yi-1] {
					perimeter += width
				}
				if yi == h-1 || !grid[xi][yi+1] {
					perimeter += width
				}
			}
		}

		fmt.Fprintln(out, perimeter)
	}
}
