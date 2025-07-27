package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const LIMIT int64 = 10000000000 // 1e10

type Segment struct {
	x1, y1, x2, y2 int64
	horizontal     bool
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// coordinate sets
	xSet := make(map[int64]struct{})
	ySet := make(map[int64]struct{})

	addX := func(v int64) { xSet[v] = struct{}{} }
	addY := func(v int64) { ySet[v] = struct{}{} }

	// borders
	addX(0)
	addX(LIMIT)
	addX(LIMIT + 1)
	addY(0)
	addY(LIMIT)
	addY(LIMIT + 1)

	// start position (center)
	start := LIMIT / 2
	posX, posY := start, start
	addX(posX)
	addX(posX + 1)
	addY(posY)
	addY(posY + 1)

	segments := make([]Segment, 0, n)

	for i := 0; i < n; i++ {
		var d string
		var l int64
		fmt.Fscan(in, &d, &l)
		nx, ny := posX, posY
		switch d[0] {
		case 'L':
			nx = posX - l
			// horizontal
			addY(posY)
			addY(posY + 1)
			addX(min(posX, nx))
			addX(max(posX, nx)+1)
			segments = append(segments, Segment{posX, posY, nx, ny, true})
		case 'R':
			nx = posX + l
			addY(posY)
			addY(posY + 1)
			addX(min(posX, nx))
			addX(max(posX, nx)+1)
			segments = append(segments, Segment{posX, posY, nx, ny, true})
		case 'D':
			ny = posY - l
			// vertical
			addX(posX)
			addX(posX + 1)
			addY(min(posY, ny))
			addY(max(posY, ny)+1)
			segments = append(segments, Segment{posX, posY, nx, ny, false})
		case 'U':
			ny = posY + l
			addX(posX)
			addX(posX + 1)
			addY(min(posY, ny))
			addY(max(posY, ny)+1)
			segments = append(segments, Segment{posX, posY, nx, ny, false})
		}
		posX, posY = nx, ny
	}

	// coordinate compression
	xs := make([]int64, 0, len(xSet))
	for v := range xSet {
		xs = append(xs, v)
	}
	ys := make([]int64, 0, len(ySet))
	for v := range ySet {
		ys = append(ys, v)
	}
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })

	xIndex := make(map[int64]int, len(xs))
	for i, v := range xs {
		xIndex[v] = i
	}
	yIndex := make(map[int64]int, len(ys))
	for i, v := range ys {
		yIndex[v] = i
	}

	w := len(xs) - 1
	h := len(ys) - 1

	// sprayed grid
	sprayed := make([][]bool, h)
	visited := make([][]bool, h)
	for i := range sprayed {
		sprayed[i] = make([]bool, w)
		visited[i] = make([]bool, w)
	}

	// mark segments
	for _, s := range segments {
		if s.horizontal {
			yIdx := yIndex[s.y1]
			lx := min(s.x1, s.x2)
			rx := max(s.x1, s.x2)
			lIdx := xIndex[lx]
			rIdx := xIndex[rx+1]
			for i := lIdx; i < rIdx; i++ {
				sprayed[yIdx][i] = true
			}
		} else { // vertical
			xIdx := xIndex[s.x1]
			ly := min(s.y1, s.y2)
			ry := max(s.y1, s.y2)
			bIdx := yIndex[ly]
			tIdx := yIndex[ry+1]
			for j := bIdx; j < tIdx; j++ {
				sprayed[j][xIdx] = true
			}
		}
	}

	// starting cell
	sprayed[yIndex[start]][xIndex[start]] = true

	// BFS from border
	type pair struct{ x, y int }
	queue := make([]pair, 0)

	push := func(i, j int) {
		visited[j][i] = true
		queue = append(queue, pair{i, j})
	}

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if sprayed[j][i] || visited[j][i] {
				continue
			}
			if xs[i] == 0 || ys[j] == 0 || xs[i+1] == LIMIT+1 || ys[j+1] == LIMIT+1 {
				push(i, j)
			}
		}
	}

	dir4 := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for idx := 0; idx < len(queue); idx++ {
		c := queue[idx]
		for _, d := range dir4 {
			ni, nj := c.x+d[0], c.y+d[1]
			if ni < 0 || nj < 0 || ni >= w || nj >= h {
				continue
			}
			if sprayed[nj][ni] || visited[nj][ni] {
				continue
			}
			push(ni, nj)
		}
	}

	// compute answer
	var total int64 = 0
	for j := 0; j < h; j++ {
		dy := ys[j+1] - ys[j]
		for i := 0; i < w; i++ {
			dx := xs[i+1] - xs[i]
			area := dx * dy
			if sprayed[j][i] || !visited[j][i] {
				// safe: either sprayed or enclosed
				total += area
			}
		}
	}

	fmt.Println(total)
}