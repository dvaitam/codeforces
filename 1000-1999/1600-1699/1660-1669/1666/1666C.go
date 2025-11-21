package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct {
	x, y int64
}

type segment struct {
	x1, y1, x2, y2 int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var pts [3]point
	for i := 0; i < 3; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}

	xs := []int64{pts[0].x, pts[1].x, pts[2].x}
	ys := []int64{pts[0].y, pts[1].y, pts[2].y}
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
	my := ys[1]

	minX, maxX := pts[0].x, pts[0].x
	for i := 1; i < 3; i++ {
		if pts[i].x < minX {
			minX = pts[i].x
		}
		if pts[i].x > maxX {
			maxX = pts[i].x
		}
	}

	groups := make(map[int64][]int64)
	for _, pt := range pts {
		groups[pt.x] = append(groups[pt.x], pt.y)
	}

	var segments []segment

	for x, ys := range groups {
		minY := my
		maxY := my
		for _, y := range ys {
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
		if minY != maxY {
			segments = append(segments, segment{x, minY, x, maxY})
		}
	}

	if minX != maxX {
		lx := minX
		rx := maxX
		if lx > rx {
			lx, rx = rx, lx
		}
		segments = append(segments, segment{lx, my, rx, my})
	}

	fmt.Fprintln(out, len(segments))
	for _, seg := range segments {
		fmt.Fprintf(out, "%d %d %d %d\n", seg.x1, seg.y1, seg.x2, seg.y2)
	}
}
