package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Point struct{ x, y int64 }

func cross(o, a, b Point) int64 {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

func convexHull(points []Point) []Point {
	sort.Slice(points, func(i, j int) bool {
		if points[i].x == points[j].x {
			return points[i].y < points[j].y
		}
		return points[i].x < points[j].x
	})
	n := len(points)
	if n <= 1 {
		return points
	}
	hull := make([]Point, 0, 2*n)
	for _, p := range points {
		for len(hull) >= 2 && cross(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	k := len(hull)
	for i := n - 2; i >= 0; i-- {
		p := points[i]
		for len(hull) > k && cross(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	if len(hull) > 1 {
		hull = hull[:len(hull)-1]
	}
	return hull
}

func inside(a, b, p Point) bool { return cross(a, b, p) >= 0 }

func lineInter(p1, p2, a, b Point) Point {
	A1 := p2.y - p1.y
	B1 := p1.x - p2.x
	C1 := A1*p1.x + B1*p1.y
	A2 := b.y - a.y
	B2 := a.x - b.x
	C2 := A2*a.x + B2*a.y
	det := A1*B2 - A2*B1
	x := (B2*C1 - B1*C2) / det
	y := (A1*C2 - A2*C1) / det
	return Point{x, y}
}

func clip(poly, clipPoly []Point) []Point {
	out := poly
	for i := 0; i < len(clipPoly); i++ {
		a := clipPoly[i]
		b := clipPoly[(i+1)%len(clipPoly)]
		if len(out) == 0 {
			break
		}
		var res []Point
		prev := out[len(out)-1]
		prevIn := inside(a, b, prev)
		for _, cur := range out {
			curIn := inside(a, b, cur)
			if curIn {
				if !prevIn {
					res = append(res, lineInter(prev, cur, a, b))
				}
				res = append(res, cur)
			} else if prevIn {
				res = append(res, lineInter(prev, cur, a, b))
			}
			prev = cur
			prevIn = curIn
		}
		out = res
	}
	return out
}

func shift(poly []Point, dx, dy int64) []Point {
	r := make([]Point, len(poly))
	for i, p := range poly {
		r[i] = Point{p.x + dx, p.y + dy}
	}
	return r
}

func dedup(poly []Point) []Point {
	if len(poly) == 0 {
		return poly
	}
	res := []Point{poly[0]}
	for i := 1; i < len(poly); i++ {
		if poly[i].x != poly[i-1].x || poly[i].y != poly[i-1].y {
			res = append(res, poly[i])
		}
	}
	if len(res) > 1 && res[0].x == res[len(res)-1].x && res[0].y == res[len(res)-1].y {
		res = res[:len(res)-1]
	}
	return res
}

func polygonArea(poly []Point) int64 {
	s := int64(0)
	n := len(poly)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		s += poly[i].x*poly[j].y - poly[j].x*poly[i].y
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		var N, M int
		if _, err := fmt.Fscan(in, &N, &M); err != nil {
			return
		}
		if N == 0 && M == 0 {
			break
		}
		pts := make([]Point, 0, 4*M)
		for i := 0; i < M; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			pts = append(pts, Point{int64(x - 1), int64(y - 1)}, Point{int64(x), int64(y - 1)}, Point{int64(x - 1), int64(y)}, Point{int64(x), int64(y)})
		}
		hull := convexHull(pts)
		// ensure CCW orientation
		if polygonArea(hull) < 0 {
			for i, j := 0, len(hull)-1; i < j; i, j = i+1, j-1 {
				hull[i], hull[j] = hull[j], hull[i]
			}
		}
		poly := shift(hull, -1, -1)
		shifts := []Point{{-1, 1}, {1, -1}, {1, 1}}
		for _, s := range shifts {
			poly = clip(poly, shift(hull, s.x, s.y))
		}
		poly = dedup(poly)
		if polygonArea(poly) > 0 {
			for i, j := 0, len(poly)-1; i < j; i, j = i+1, j-1 {
				poly[i], poly[j] = poly[j], poly[i]
			}
		}
		idx := 0
		for i := 1; i < len(poly); i++ {
			if poly[i].x < poly[idx].x || (poly[i].x == poly[idx].x && poly[i].y < poly[idx].y) {
				idx = i
			}
		}
		fmt.Println(len(poly))
		for i := 0; i < len(poly); i++ {
			v := poly[(idx+i)%len(poly)]
			fmt.Printf("%d %d\n", v.x, v.y)
		}
	}
}
