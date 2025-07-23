package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const eps = 1e-9

type Point struct {
	x, y float64
}

func sub(a, b Point) Point         { return Point{a.x - b.x, a.y - b.y} }
func add(a, b Point) Point         { return Point{a.x + b.x, a.y + b.y} }
func mul(a Point, t float64) Point { return Point{a.x * t, a.y * t} }
func dot(a, b Point) float64       { return a.x*b.x + a.y*b.y }
func cross(a, b Point) float64     { return a.x*b.y - a.y*b.x }

func onSegment(p, a, b Point) bool {
	if math.Abs(cross(sub(b, a), sub(p, a))) > eps {
		return false
	}
	return dot(sub(p, a), sub(p, b)) <= eps
}

func pointInPoly(p Point, poly []Point) bool {
	// boundary check
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		if onSegment(p, a, b) {
			return true
		}
	}
	wn := 0
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		if a.y <= p.y {
			if b.y > p.y && cross(sub(b, a), sub(p, a)) > eps {
				wn++
			}
		} else {
			if b.y <= p.y && cross(sub(b, a), sub(p, a)) < -eps {
				wn--
			}
		}
	}
	return wn != 0
}

func intersectionLength(poly []Point, p0, p1 Point) float64 {
	d := sub(p1, p0)
	l := math.Hypot(d.x, d.y)
	dir := Point{d.x / l, d.y / l}
	var ts []float64
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		da := sub(b, a)
		denom := cross(d, da)
		if math.Abs(denom) < eps {
			if math.Abs(cross(d, sub(a, p0))) < eps {
				t1 := dot(sub(a, p0), dir)
				t2 := dot(sub(b, p0), dir)
				if t2 < t1 {
					t1, t2 = t2, t1
				}
				ts = append(ts, t1, t2)
			}
			continue
		}
		t := cross(sub(a, p0), da) / denom
		s := cross(sub(a, p0), d) / denom
		if s >= -eps && s <= 1+eps {
			ts = append(ts, t)
		}
	}
	if len(ts) == 0 {
		return 0
	}
	sort.Float64s(ts)
	uniq := []float64{ts[0]}
	for i := 1; i < len(ts); i++ {
		if math.Abs(ts[i]-ts[i-1]) > eps {
			uniq = append(uniq, ts[i])
		}
	}
	res := 0.0
	for i := 0; i+1 < len(uniq); i++ {
		tmid := (uniq[i] + uniq[i+1]) / 2
		pt := add(p0, mul(dir, tmid))
		if pointInPoly(pt, poly) {
			res += uniq[i+1] - uniq[i]
		}
	}
	return res * l
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	poly := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &poly[i].x, &poly[i].y)
	}
	for i := 0; i < m; i++ {
		var x1, y1, x2, y2 float64
		fmt.Fscan(reader, &x1, &y1, &x2, &y2)
		p0 := Point{x1, y1}
		p1 := Point{x2, y2}
		length := intersectionLength(poly, p0, p1)
		fmt.Fprintf(writer, "%.10f\n", length)
	}
}
