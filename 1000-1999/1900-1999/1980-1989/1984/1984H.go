package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x, y float64
}

func cross(a, b Point) float64 {
	return a.x*b.y - a.y*b.x
}

func sub(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y}
}

func pointInTriangle(a, b, c, p Point) bool {
	ab := sub(b, a)
	ap := sub(p, a)
	bc := sub(c, b)
	bp := sub(p, b)
	ca := sub(a, c)
	cp := sub(p, c)
	c1 := cross(ab, ap)
	c2 := cross(bc, bp)
	c3 := cross(ca, cp)
	return (c1 >= 0 && c2 >= 0 && c3 >= 0) || (c1 <= 0 && c2 <= 0 && c3 <= 0)
}

func circumcircle(a, b, c Point) (float64, float64, float64) {
	x1, y1 := a.x, a.y
	x2, y2 := b.x, b.y
	x3, y3 := c.x, c.y
	d := 2 * (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2))
	cx := ((x1*x1+y1*y1)*(y2-y3) + (x2*x2+y2*y2)*(y3-y1) + (x3*x3+y3*y3)*(y1-y2)) / d
	cy := ((x1*x1+y1*y1)*(x3-x2) + (x2*x2+y2*y2)*(x1-x3) + (x3*x3+y3*y3)*(x2-x1)) / d
	dx := cx - x1
	dy := cy - y1
	r2 := dx*dx + dy*dy
	return cx, cy, r2
}

func circleContainsAll(a, b, c Point, pts []Point) bool {
	cx, cy, r2 := circumcircle(a, b, c)
	eps := 1e-7
	for _, p := range pts {
		dx := p.x - cx
		dy := p.y - cy
		if dx*dx+dy*dy > r2+eps {
			return false
		}
	}
	return true
}

func insideSet(a, b, c Point, pts []Point) []bool {
	res := make([]bool, len(pts))
	for i, p := range pts {
		if pointInTriangle(a, b, c, p) {
			res[i] = true
		}
	}
	return res
}

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		pts := make([]Point, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pts[i].x, &pts[i].y)
		}
		if n <= 2 {
			fmt.Fprintln(out, 1)
			continue
		}
		count1 := 0
		validR := make([]int, 0)
		for r := 2; r < n; r++ {
			if circleContainsAll(pts[0], pts[1], pts[r], pts) {
				inside := insideSet(pts[0], pts[1], pts[r], pts)
				all := true
				for i := 0; i < n; i++ {
					if !inside[i] {
						all = false
						break
					}
				}
				if all {
					count1++
				}
				validR = append(validR, r)
			}
		}
		if count1 > 0 {
			fmt.Fprintln(out, count1%mod)
			continue
		}
		count2 := 0
		for _, r := range validR {
			inside1 := insideSet(pts[0], pts[1], pts[r], pts)
			owned := inside1
			pairs := [][2]int{{0, 1}, {0, r}, {1, r}}
			for _, pq := range pairs {
				p, q := pq[0], pq[1]
				for s := 0; s < n; s++ {
					if owned[s] {
						continue
					}
					if circleContainsAll(pts[p], pts[q], pts[s], pts) {
						inside2 := insideSet(pts[p], pts[q], pts[s], pts)
						all := true
						for i := 0; i < n; i++ {
							if !owned[i] && !inside2[i] {
								all = false
								break
							}
						}
						if all {
							count2++
						}
					}
				}
			}
		}
		if count2 > 0 {
			fmt.Fprintln(out, count2%mod)
		} else {
			fmt.Fprintln(out, 0)
		}
	}
}
