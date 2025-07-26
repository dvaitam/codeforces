package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Point struct {
	x, y float64
}

type Segment struct {
	a, b Point
}

func sub(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y}
}

func dot(a, b Point) float64 {
	return a.x*b.x + a.y*b.y
}

func cross(a, b Point) float64 {
	return a.x*b.y - a.y*b.x
}

func norm(a Point) float64 {
	return math.Hypot(a.x, a.y)
}

func normalize(a Point) Point {
	r := norm(a)
	if r == 0 {
		return Point{0, 0}
	}
	return Point{a.x / r, a.y / r}
}

func angleDiff(a, b Point) float64 {
	a = normalize(a)
	b = normalize(b)
	c := math.Abs(dot(a, b))
	if c > 1 {
		c = 1
	}
	return math.Acos(c)
}

func isOnRay(p, d, t Point) (bool, float64) {
	v := sub(t, p)
	if math.Abs(cross(v, d)) > 1e-9 {
		return false, 0
	}
	k := dot(v, d) / dot(d, d)
	if k >= -1e-9 {
		return true, k
	}
	return false, 0
}

func intersectRaySegment(p, d Point, s Segment) (bool, float64) {
	a := s.a
	b := s.b
	v := sub(b, a)
	w := sub(a, p)
	denom := cross(d, v)
	if math.Abs(denom) < 1e-9 {
		if math.Abs(cross(w, d)) < 1e-9 {
			// collinear
			t1 := dot(sub(a, p), d) / dot(d, d)
			t2 := dot(sub(b, p), d) / dot(d, d)
			best := math.Inf(1)
			if t1 >= -1e-9 {
				best = math.Min(best, t1)
			}
			if t2 >= -1e-9 {
				best = math.Min(best, t2)
			}
			if best < math.Inf(1) {
				return true, best
			}
		}
		return false, 0
	}
	t := cross(w, v) / denom
	u := cross(w, d) / denom
	if t >= -1e-9 && u >= -1e-9 && u <= 1+1e-9 {
		return true, t
	}
	return false, 0
}

func attempt(start Point, d Point, segs []Segment) bool {
	p := start
	visited := make(map[[2]int]bool)
	for iter := 0; iter <= len(segs); iter++ {
		hasT, distT := isOnRay(p, d, Point{0, 0})
		bestDist := math.Inf(1)
		bestIdx := -1
		for i, s := range segs {
			ok, dist := intersectRaySegment(p, d, s)
			if ok && dist > 1e-9 && dist < bestDist {
				bestDist = dist
				bestIdx = i
			}
		}
		if hasT && (bestIdx == -1 || distT <= bestDist+1e-9) {
			return true
		}
		if bestIdx == -1 {
			return false
		}
		seg := segs[bestIdx]
		v := sub(seg.b, seg.a)
		endpoint := seg.b
		slideVec := v
		if dot(d, v) < 0 {
			endpoint = seg.a
			slideVec = Point{-v.x, -v.y}
		}
		if angleDiff(d, slideVec) >= math.Pi/4-1e-9 {
			return false
		}
		key := [2]int{int(endpoint.x*1000 + 0.5), int(endpoint.y*1000 + 0.5)}
		if visited[key] {
			return false
		}
		visited[key] = true
		p = endpoint
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	segs := make([]Segment, n)
	for i := 0; i < n; i++ {
		var ax, ay, bx, by float64
		fmt.Fscan(reader, &ax, &ay, &bx, &by)
		segs[i] = Segment{Point{ax, ay}, Point{bx, by}}
	}
	var q int
	fmt.Fscan(reader, &q)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; q > 0; q-- {
		var sx, sy float64
		fmt.Fscan(reader, &sx, &sy)
		start := Point{sx, sy}
		dirs := []Point{{-sx, -sy}}
		for _, s := range segs {
			dirs = append(dirs, Point{-s.a.x, -s.a.y}, Point{-s.b.x, -s.b.y})
		}
		found := false
		for _, d := range dirs {
			if d.x == 0 && d.y == 0 {
				continue
			}
			if attempt(start, d, segs) {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
