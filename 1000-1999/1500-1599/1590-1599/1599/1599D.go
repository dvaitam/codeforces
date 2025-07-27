package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct {
	x, y int64
	idx  int
}

func cross(ax, ay, bx, by int64) int64 {
	return ax*by - ay*bx
}

func convexHull(pts []Point) []Point {
	if len(pts) <= 1 {
		res := make([]Point, len(pts))
		copy(res, pts)
		return res
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	lower := make([]Point, 0)
	for _, p := range pts {
		for len(lower) >= 2 {
			a := lower[len(lower)-2]
			b := lower[len(lower)-1]
			if cross(b.x-a.x, b.y-a.y, p.x-b.x, p.y-b.y) <= 0 {
				lower = lower[:len(lower)-1]
			} else {
				break
			}
		}
		lower = append(lower, p)
	}
	upper := make([]Point, 0)
	for i := len(pts) - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 {
			a := upper[len(upper)-2]
			b := upper[len(upper)-1]
			if cross(b.x-a.x, b.y-a.y, p.x-b.x, p.y-b.y) <= 0 {
				upper = upper[:len(upper)-1]
			} else {
				break
			}
		}
		upper = append(upper, p)
	}
	upper = upper[1 : len(upper)-1]
	res := append(lower, upper...)
	return res
}

func onionLayers(pts []Point) [][]Point {
	res := [][]Point{}
	remaining := make([]Point, len(pts))
	copy(remaining, pts)
	for len(remaining) > 0 {
		hull := convexHull(remaining)
		res = append(res, hull)
		marked := make(map[int]bool, len(hull))
		for _, p := range hull {
			marked[p.idx] = true
		}
		tmp := make([]Point, 0, len(remaining)-len(hull))
		for _, p := range remaining {
			if !marked[p.idx] {
				tmp = append(tmp, p)
			}
		}
		remaining = tmp
	}
	return res
}

func extremeIdx(hull []Point, dx, dy int64) int {
	nx, ny := -dy, dx
	best := 0
	bestVal := hull[0].x*nx + hull[0].y*ny
	for i := 1; i < len(hull); i++ {
		val := hull[i].x*nx + hull[i].y*ny
		if val < bestVal {
			bestVal = val
			best = i
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &pts[i].x, &pts[i].y)
		pts[i].idx = i + 1
	}
	layers := onionLayers(pts)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var dx, dy int64
		var k int
		fmt.Fscan(reader, &dx, &dy, &k)
		curDx, curDy := dx, dy
		layerIdx := 0
		for {
			hull := layers[layerIdx]
			start := extremeIdx(hull, curDx, curDy)
			m := len(hull)
			if k <= m {
				ansIdx := hull[(start+k-1)%m].idx
				fmt.Fprintln(writer, ansIdx)
				break
			}
			k -= m
			if m > 1 {
				prev := hull[(start+m-2)%m]
				last := hull[(start+m-1)%m]
				curDx = last.x - prev.x
				curDy = last.y - prev.y
			}
			layerIdx++
		}
	}
}
