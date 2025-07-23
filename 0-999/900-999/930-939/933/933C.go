package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Circle struct {
	x, y, r float64
}

type Point struct {
	x, y float64
}

const eps = 1e-9

func intersections(a, b Circle) []Point {
	dx := b.x - a.x
	dy := b.y - a.y
	d := math.Hypot(dx, dy)
	if d > a.r+b.r+eps || d < math.Abs(a.r-b.r)-eps || d == 0 {
		return nil
	}
	// distance from a's center to the line between intersection points
	alpha := (a.r*a.r - b.r*b.r + d*d) / (2 * d)
	h2 := a.r*a.r - alpha*alpha
	if h2 < eps {
		h2 = 0
	}
	xm := a.x + alpha*dx/d
	ym := a.y + alpha*dy/d
	if h2 == 0 {
		return []Point{{xm, ym}}
	}
	h := math.Sqrt(h2)
	rx := -dy * h / d
	ry := dx * h / d
	return []Point{{xm + rx, ym + ry}, {xm - rx, ym - ry}}
}

func addPoint(pts []Point, p Point) ([]Point, int) {
	for i, q := range pts {
		if (p.x-q.x)*(p.x-q.x)+(p.y-q.y)*(p.y-q.y) < eps*eps {
			return pts, i
		}
	}
	return append(pts, p), len(pts)
}

func find(par []int, x int) int {
	if par[x] != x {
		par[x] = find(par, par[x])
	}
	return par[x]
}

func union(par []int, x, y int) {
	rx := find(par, x)
	ry := find(par, y)
	if rx != ry {
		par[ry] = rx
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	circles := make([]Circle, n)
	for i := 0; i < n; i++ {
		var xi, yi, ri float64
		fmt.Fscan(reader, &xi, &yi, &ri)
		circles[i] = Circle{xi, yi, ri}
	}

	points := make([]Point, 0)
	sets := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		sets[i] = make(map[int]struct{})
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for _, p := range intersections(circles[i], circles[j]) {
				var idx int
				points, idx = addPoint(points, p)
				sets[i][idx] = struct{}{}
				sets[j][idx] = struct{}{}
			}
		}
	}

	V := len(points)
	E := 0
	loops := 0
	for i := 0; i < n; i++ {
		if len(sets[i]) == 0 {
			loops++
			E++
		} else {
			E += len(sets[i])
		}
	}
	V += loops

	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// if they share at least one intersection point
			for idx := range sets[i] {
				if _, ok := sets[j][idx]; ok {
					union(parent, i, j)
					break
				}
			}
		}
	}
	compSet := make(map[int]struct{})
	for i := 0; i < n; i++ {
		compSet[find(parent, i)] = struct{}{}
	}
	C := len(compSet)

	result := E - V + C + 1
	fmt.Println(result)
}
