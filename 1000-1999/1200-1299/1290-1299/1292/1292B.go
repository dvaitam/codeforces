package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y int64
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var x0, y0, ax, ay, bx, by int64
	fmt.Fscan(in, &x0, &y0, &ax, &ay, &bx, &by)
	var xs, ys, t int64
	fmt.Fscan(in, &xs, &ys, &t)

	// Generate data nodes within a reasonable bound
	pts := []Point{}
	x, y := x0, y0
	limit := xs + t
	if ys+t > limit {
		limit = ys + t
	}
	maxCoord := limit * 2
	for x <= maxCoord && y <= maxCoord {
		pts = append(pts, Point{x: x, y: y})
		x = ax*x + bx
		y = ay*y + by
	}

	ans := 0
	n := len(pts)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			dist := (pts[j].x - pts[i].x) + (pts[j].y - pts[i].y)
			startToI := abs(xs-pts[i].x) + abs(ys-pts[i].y)
			startToJ := abs(xs-pts[j].x) + abs(ys-pts[j].y)
			need := dist + startToI
			if startToJ < startToI {
				need = dist + startToJ
			}
			if need <= t {
				if j-i+1 > ans {
					ans = j - i + 1
				}
			}
		}
	}

	fmt.Println(ans)
}
