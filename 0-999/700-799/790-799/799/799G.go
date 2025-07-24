package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Point struct{ x, y float64 }

func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func sub(a, b Point) Point     { return Point{a.x - b.x, a.y - b.y} }

func polygonArea(poly []Point) float64 {
	area := 0.0
	n := len(poly)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += cross(poly[i], poly[j])
	}
	if area < 0 {
		area = -area
	}
	return area / 2
}

// cutPolygon cuts poly by line through p along dir and keeps points on the left side.
func cutPolygon(poly []Point, p, dir Point) []Point {
	res := make([]Point, 0, len(poly)+2)
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		da := cross(dir, sub(a, p))
		db := cross(dir, sub(b, p))
		if da >= 0 { // point a is kept
			res = append(res, a)
		}
		if da*db < 0 { // segment intersects line
			t := da / (da - db)
			inter := Point{a.x + t*(b.x-a.x), a.y + t*(b.y-a.y)}
			res = append(res, inter)
		}
	}
	return res
}

func findAngle(poly []Point, p Point, total float64) float64 {
	lo, hi := 0.0, math.Pi
	for iter := 0; iter < 60; iter++ {
		mid := (lo + hi) / 2
		dir := Point{math.Cos(mid), math.Sin(mid)}
		half := cutPolygon(poly, p, dir)
		area := polygonArea(half)
		if area > total/2 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return (lo + hi) / 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	poly := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &poly[i].x, &poly[i].y)
	}
	total := polygonArea(poly)
	for ; q > 0; q-- {
		var x, y float64
		fmt.Fscan(in, &x, &y)
		angle := findAngle(poly, Point{x, y}, total)
		fmt.Fprintf(out, "%.10f\n", angle)
	}
}
