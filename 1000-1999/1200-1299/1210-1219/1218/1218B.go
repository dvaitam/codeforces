package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
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

func cutPolygon(poly []Point, dir Point) []Point {
	res := make([]Point, 0, len(poly)+2)
	n := len(poly)
	for i := 0; i < n; i++ {
		a := poly[i]
		b := poly[(i+1)%n]
		da := cross(dir, a)
		db := cross(dir, b)
		if da >= -1e-9 {
			res = append(res, a)
		}
		if da*db < -1e-9 {
			t := da / (da - db)
			res = append(res, Point{a.x + (b.x-a.x)*t, a.y + (b.y-a.y)*t})
		}
	}
	return res
}

func raySegIntersect(dir, a, b Point) (float64, bool) {
	d2 := sub(b, a)
	den := cross(dir, d2)
	if math.Abs(den) < 1e-12 {
		return 0, false
	}
	t := cross(a, d2) / den
	u := cross(a, dir) / den
	if t >= 0 && u >= 0 && u <= 1 {
		return t, true
	}
	return 0, false
}

type Polygon struct {
	verts  []Point
	angles []float64 // increasing, len = m+1
	start  float64
	span   float64
}

func newPolygon(pts []Point) Polygon {
	m := len(pts)
	angles := make([]float64, m)
	for i, p := range pts {
		ang := math.Atan2(p.y, p.x)
		if ang < 0 {
			ang += 2 * math.Pi
		}
		angles[i] = ang
	}
	start := 0
	for i := 1; i < m; i++ {
		if angles[i] < angles[start] {
			start = i
		}
	}
	ordVerts := make([]Point, m)
	ordAng := make([]float64, m+1)
	prev := angles[start]
	ordVerts[0] = pts[start]
	ordAng[0] = prev
	for k := 1; k < m; k++ {
		idx := (start + k) % m
		ang := angles[idx]
		if ang < prev {
			ang += 2 * math.Pi
		}
		ordAng[k] = ang
		ordVerts[k] = pts[idx]
		prev = ang
	}
	ordAng[m] = ordAng[0] + 2*math.Pi
	span := ordAng[m-1] - ordAng[0]
	normAng := make([]float64, m+1)
	base := ordAng[0]
	for i := 0; i <= m; i++ {
		normAng[i] = ordAng[i] - base
	}
	return Polygon{verts: ordVerts, angles: normAng, start: math.Mod(base, 2*math.Pi), span: span}
}

func (p Polygon) rayDist(angle float64) (float64, bool) {
	a := angle
	if a < p.start {
		a += 2 * math.Pi
	}
	if a < p.start || a > p.start+p.span {
		return 0, false
	}
	a -= p.start
	idx := sort.Search(len(p.angles)-1, func(i int) bool { return p.angles[i+1] >= a-1e-12 })
	if idx >= len(p.angles)-1 {
		idx = len(p.angles) - 2
	}
	v1 := p.verts[idx]
	v2 := p.verts[(idx+1)%len(p.verts)]
	return raySegIntersect(Point{math.Cos(angle), math.Sin(angle)}, v1, v2)
}

func (p Polygon) wedgeArea(a1, a2 float64) float64 {
	poly := make([]Point, len(p.verts))
	copy(poly, p.verts)
	dir1 := Point{math.Cos(a1), math.Sin(a1)}
	poly = cutPolygon(poly, dir1)
	dir2 := Point{math.Cos(a2), math.Sin(a2)}
	poly = cutPolygon(poly, Point{-dir2.x, -dir2.y})
	if len(poly) < 3 {
		return 0
	}
	return polygonArea(poly)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	polys := make([]Polygon, n)
	angles := []float64{0, 2 * math.Pi}
	type event struct {
		ang float64
		id  int
		typ int
	}
	var events []event
	for i := 0; i < n; i++ {
		var c int
		fmt.Fscan(in, &c)
		pts := make([]Point, c)
		for j := 0; j < c; j++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			pts[j] = Point{float64(x), float64(y)}
		}
		poly := newPolygon(pts)
		polys[i] = poly
		s := poly.start
		e := poly.start + poly.span
		if e <= 2*math.Pi {
			events = append(events, event{s, i, 1}, event{e, i, -1})
			angles = append(angles, s, e)
		} else {
			events = append(events, event{s, i, 1}, event{2 * math.Pi, i, -1}, event{0, i, 1}, event{e - 2*math.Pi, i, -1})
			angles = append(angles, s, 2*math.Pi, 0, e-2*math.Pi)
		}
		for _, ang := range poly.angles[:len(poly.angles)-1] {
			a := math.Mod(poly.start+ang, 2*math.Pi)
			angles = append(angles, a)
		}
	}
	sort.Float64s(angles)
	anglesUni := []float64{angles[0]}
	for _, a := range angles[1:] {
		if a-anglesUni[len(anglesUni)-1] > 1e-12 {
			anglesUni = append(anglesUni, a)
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].ang < events[j].ang })
	active := make(map[int]struct{})
	idxEv := 0
	area := 0.0
	for i := 0; i < len(anglesUni)-1; i++ {
		a1 := anglesUni[i]
		a2 := anglesUni[i+1]
		mid := (a1 + a2) / 2
		for idxEv < len(events) && events[idxEv].ang <= a1+1e-12 {
			ev := events[idxEv]
			if ev.typ == 1 {
				active[ev.id] = struct{}{}
			} else {
				delete(active, ev.id)
			}
			idxEv++
		}
		bestDist := math.Inf(1)
		bestID := -1
		for id := range active {
			if d, ok := polys[id].rayDist(mid); ok {
				if d < bestDist {
					bestDist = d
					bestID = id
				}
			}
		}
		if bestID >= 0 {
			area += polys[bestID].wedgeArea(a1, a2)
		}
	}
	fmt.Printf("%.10f\n", area)
}
