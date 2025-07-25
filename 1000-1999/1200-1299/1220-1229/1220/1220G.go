package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type Point struct {
	x, y int64
}

func int64sqrt(x int64) float64 {
	return math.Sqrt(float64(x))
}

func circleIntersections(p1 Point, r1 float64, p2 Point, r2 float64) []Point {
	dx := float64(p2.x - p1.x)
	dy := float64(p2.y - p1.y)
	d := math.Hypot(dx, dy)
	if d > r1+r2+1e-6 || d < math.Abs(r1-r2)-1e-6 || d == 0 {
		return nil
	}
	a := (r1*r1 - r2*r2 + d*d) / (2 * d)
	h2 := r1*r1 - a*a
	if h2 < -1e-6 {
		return nil
	}
	if h2 < 0 {
		h2 = 0
	}
	h := math.Sqrt(h2)
	xm := float64(p1.x) + a*dx/d
	ym := float64(p1.y) + a*dy/d
	rx := -dy * (h / d)
	ry := dx * (h / d)
	points := []Point{}
	x1 := math.Round(xm + rx)
	y1 := math.Round(ym + ry)
	if math.Abs((xm+rx)-x1) < 1e-6 && math.Abs((ym+ry)-y1) < 1e-6 {
		points = append(points, Point{int64(x1), int64(y1)})
	}
	x2 := math.Round(xm - rx)
	y2 := math.Round(ym - ry)
	if math.Abs((xm-rx)-x2) < 1e-6 && math.Abs((ym-ry)-y2) < 1e-6 {
		p := Point{int64(x2), int64(y2)}
		if len(points) == 0 || points[0] != p {
			points = append(points, p)
		}
	}
	return points
}

func checkPoint(p Point, ants []Point, d []int64) bool {
	n := len(ants)
	arr := make([]int64, n)
	for i, a := range ants {
		dx := p.x - a.x
		dy := p.y - a.y
		arr[i] = dx*dx + dy*dy
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	for i := 0; i < n; i++ {
		if arr[i] != d[i] {
			return false
		}
	}
	return true
}

func solveQuery(ants []Point, d []int64) []Point {
	n := len(ants)
	dmin := d[0]
	dmax := d[n-1]
	var ans []Point
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			pts := circleIntersections(ants[i], int64sqrt(dmin), ants[j], int64sqrt(dmax))
			for _, p := range pts {
				if checkPoint(p, ants, d) {
					ans = append(ans, p)
				}
			}
		}
	}
	// deduplicate
	if len(ans) > 1 {
		sort.Slice(ans, func(i, j int) bool {
			if ans[i].x == ans[j].x {
				return ans[i].y < ans[j].y
			}
			return ans[i].x < ans[j].x
		})
		uniq := ans[:1]
		for k := 1; k < len(ans); k++ {
			if ans[k] != ans[k-1] {
				uniq = append(uniq, ans[k])
			}
		}
		ans = uniq
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ants := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &ants[i].x, &ants[i].y)
	}
	var m int
	fmt.Fscan(reader, &m)
	for q := 0; q < m; q++ {
		d := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &d[i])
		}
		sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
		pts := solveQuery(ants, d)
		fmt.Fprintln(writer, len(pts))
		sort.Slice(pts, func(i, j int) bool {
			if pts[i].x == pts[j].x {
				return pts[i].y < pts[j].y
			}
			return pts[i].x < pts[j].x
		})
		for _, p := range pts {
			fmt.Fprintf(writer, "%d %d\n", p.x, p.y)
		}
	}
}
