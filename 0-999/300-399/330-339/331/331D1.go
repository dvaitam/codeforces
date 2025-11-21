package main

import (
	"bufio"
	"fmt"
	"os"
)

type Arrow struct {
	x0, y0, x1, y1 int
}

type Elem struct {
	pos   int
	arrow int
}

type Query struct {
	x, y int
	dir  byte
	t    int64
}

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

func clamp(val, lo, hi int) int {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, b int
	fmt.Fscan(in, &n, &b)
	arrows := make([]Arrow, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arrows[i].x0, &arrows[i].y0, &arrows[i].x1, &arrows[i].y1)
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([]Query, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i].x, &queries[i].y)
		fmt.Fscan(in, &queries[i].dir, &queries[i].t)
	}

	type Key struct {
		pos int
		dir byte
	}

	results := make([][2]int, q)
	visited := make(map[Key]bool)
	for idx := range queries {
		curX := queries[idx].x
		curY := queries[idx].y
		curDir := queries[idx].dir
		remaining := queries[idx].t

		for remaining > 0 {
			key := Key{pos: curX*(b+1) + curY, dir: curDir}
			if visited[key] {
				break
			}
			visited[key] = true

			var nextX, nextY int
			var dist int64
			switch curDir {
			case 'U':
				nextX, nextY, dist = nextUp(curX, curY, b, arrows)
			case 'D':
				nextX, nextY, dist = nextDown(curX, curY, b, arrows)
			case 'L':
				nextX, nextY, dist = nextLeft(curX, curY, b, arrows)
			case 'R':
				nextX, nextY, dist = nextRight(curX, curY, b, arrows)
			}

			if dist >= remaining {
				results[idx] = advancePosition(curX, curY, curDir, remaining, b)
				remaining = 0
			} else {
				curX, curY = nextX, nextY
				remaining -= dist
				curDir = switchDirection(nextX, nextY, curDir, arrows)
			}
		}
		if remaining > 0 {
			results[idx] = clampPosition(curX, curY, b)
		}
	}

	for _, res := range results {
		fmt.Fprintf(out, "%d %d\n", res[0], res[1])
	}
}

func advancePosition(x, y int, dir byte, step int64, b int) [2]int {
	switch dir {
	case 'U':
		return [2]int{x, clamp(y+int(step), 0, b)}
	case 'D':
		return [2]int{x, clamp(y-int(step), 0, b)}
	case 'R':
		return [2]int{clamp(x+int(step), 0, b), y}
	case 'L':
		return [2]int{clamp(x-int(step), 0, b), y}
	}
	return [2]int{x, y}
}

func clampPosition(x, y, b int) [2]int {
	return [2]int{clamp(x, 0, b), clamp(y, 0, b)}
}

func nextUp(x, y, b int, arrows []Arrow) (int, int, int64) {
	minY := b
	nextX, nextY := x, b
	for _, arr := range arrows {
		if arr.x0 == arr.x1 && arr.x0 == x && arr.y0 >= y {
			if arr.y0 < minY {
				minY = arr.y0
				nextX, nextY = arr.x0, arr.y0
			}
		}
	}
	return nextX, nextY, int64(nextY - y)
}

func nextDown(x, y, b int, arrows []Arrow) (int, int, int64) {
	maxY := 0
	nextX, nextY := x, 0
	for _, arr := range arrows {
		if arr.x0 == arr.x1 && arr.x0 == x && arr.y0 <= y {
			if arr.y0 > maxY {
				maxY = arr.y0
				nextX, nextY = arr.x0, arr.y0
			}
		}
	}
	return nextX, nextY, int64(y - nextY)
}

func nextLeft(x, y, b int, arrows []Arrow) (int, int, int64) {
	maxX := 0
	nextX, nextY := 0, y
	for _, arr := range arrows {
		if arr.y0 == arr.y1 && arr.y0 == y && arr.x0 <= x {
			if arr.x0 > maxX {
				maxX = arr.x0
				nextX, nextY = arr.x0, arr.y0
			}
		}
	}
	return nextX, nextY, int64(x - nextX)
}

func nextRight(x, y, b int, arrows []Arrow) (int, int, int64) {
	minX := b
	nextX, nextY := b, y
	for _, arr := range arrows {
		if arr.y0 == arr.y1 && arr.y0 == y && arr.x0 >= x {
			if arr.x0 < minX {
				minX = arr.x0
				nextX, nextY = arr.x0, arr.y0
			}
		}
	}
	return nextX, nextY, int64(nextX - x)
}

func switchDirection(x, y int, dir byte, arrows []Arrow) byte {
	for _, arr := range arrows {
		if arr.x0 == arr.x1 && arr.x0 == x && arr.y0 == y {
			if arr.y1 > arr.y0 {
				return 'U'
			}
			return 'D'
		}
		if arr.y0 == arr.y1 && arr.x0 == x && arr.y0 == y {
			if arr.x1 > arr.x0 {
				return 'R'
			}
			return 'L'
		}
	}
	return dir
}
