package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type point struct {
	x int64
	y int64
}

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt64() int64 {
	sign := int64(1)
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func cross(o, a, b point) int64 {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

func convexHull(points []point) []point {
	n := len(points)
	if n <= 1 {
		cp := make([]point, n)
		copy(cp, points)
		return cp
	}
	pts := make([]point, n)
	copy(pts, points)
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	lower := make([]point, 0, n)
	for _, p := range pts {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}
	upper := make([]point, 0, n)
	for i := n - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}
	if len(lower) == 0 {
		return nil
	}
	lower = lower[:len(lower)-1]
	upper = upper[:len(upper)-1]
	hull := append(lower, upper...)
	return hull
}

func main() {
	fs := newFastScanner()
	n := int(fs.nextInt64())
	rInt := fs.nextInt64()
	r := float64(rInt)

	points := make([]point, n)
	angles := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		x := fs.nextInt64()
		y := fs.nextInt64()
		points[i] = point{x: x, y: y}
		if x != 0 || y != 0 {
			angles = append(angles, math.Atan2(float64(y), float64(x)))
		}
	}

	halfArea := 0.5 * math.Pi * r * r
	if len(angles) <= 1 {
		fmt.Printf("%.12f\n", halfArea)
		return
	}

	sort.Float64s(angles)
	m := len(angles)
	angles2 := make([]float64, 2*m)
	for i := 0; i < m; i++ {
		angles2[i] = angles[i]
		angles2[i+m] = angles[i] + 2*math.Pi
	}
	for i := 0; i < m; i++ {
		j := i + m - 1
		if angles2[j]-angles2[i] <= math.Pi+1e-12 {
			fmt.Printf("%.12f\n", halfArea)
			return
		}
	}

	hull := convexHull(points)
	if len(hull) < 2 {
		fmt.Printf("%.12f\n", halfArea)
		return
	}

	minDist := math.MaxFloat64
	for i := 0; i < len(hull); i++ {
		a := hull[i]
		b := hull[(i+1)%len(hull)]
		num := math.Abs(float64(a.x*b.y - a.y*b.x))
		dx := float64(b.x - a.x)
		dy := float64(b.y - a.y)
		den := math.Hypot(dx, dy)
		if den == 0 {
			continue
		}
		dist := num / den
		if dist < minDist {
			minDist = dist
		}
	}

	if minDist < 0 {
		minDist = 0
	}
	if minDist > r {
		minDist = r
	}
	c := minDist / r
	if c > 1 {
		c = 1
	}
	if c < -1 {
		c = -1
	}
	tmp := r*r - minDist*minDist
	if tmp < 0 {
		tmp = 0
	}
	area := r*r*math.Acos(c) - minDist*math.Sqrt(tmp)
	fmt.Printf("%.12f\n", area)
}
