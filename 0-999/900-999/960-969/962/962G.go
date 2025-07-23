package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ x, y int }
type Rect struct{ x1, y1, x2, y2 int }

func insideRect(p Point, r Rect) bool {
	return p.x > r.x1 && p.x < r.x2 && p.y > r.y2 && p.y < r.y1
}

func pointInPolygon(pt Point, poly []Point) bool {
	inside := false
	n := len(poly)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := poly[i].x, poly[i].y
		xj, yj := poly[j].x, poly[j].y
		if (yi > pt.y) != (yj > pt.y) {
			if int64(pt.x-xi) < int64(xj-xi)*int64(pt.y-yi)/int64(yj-yi) {
				inside = !inside
			}
		}
	}
	return inside
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var r Rect
	if _, err := fmt.Fscan(in, &r.x1, &r.y1, &r.x2, &r.y2); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	poly := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &poly[i].x, &poly[i].y)
	}

	cross := 0
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		ia := insideRect(a, r)
		ib := insideRect(b, r)
		if ia != ib {
			cross++
			continue
		}
		if ia || ib { // both inside already handled
			continue
		}
		if a.x == b.x {
			x := a.x
			if x > r.x1 && x < r.x2 {
				low, high := a.y, b.y
				if low > high {
					low, high = high, low
				}
				if low < r.y2 && high > r.y2 {
					cross++
				}
				if low < r.y1 && high > r.y1 {
					cross++
				}
			}
		} else if a.y == b.y {
			y := a.y
			if y > r.y2 && y < r.y1 {
				low, high := a.x, b.x
				if low > high {
					low, high = high, low
				}
				if low < r.x1 && high > r.x1 {
					cross++
				}
				if low < r.x2 && high > r.x2 {
					cross++
				}
			}
		}
	}

	if cross > 0 {
		fmt.Fprintln(out, cross/2)
		return
	}

	center := Point{(r.x1 + r.x2) / 2, (r.y1 + r.y2) / 2}
	if pointInPolygon(center, poly) || insideRect(poly[0], r) {
		fmt.Fprintln(out, 1)
	} else {
		fmt.Fprintln(out, 0)
	}
}
