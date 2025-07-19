package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type pt struct{ x, y float64 }

func sub(a, b pt) pt        { return pt{a.x - b.x, a.y - b.y} }
func cross(a, b pt) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b pt) float64   { return a.x*b.x + a.y*b.y }

func work(ext []pt, n int, ans *float64) {
	sum := 0.0
	j := 1
	for i := 0; i < n; i++ {
		// advance j while angle is obtuse
		for j+1 < len(ext) && dot(sub(ext[i+1], ext[i]), sub(ext[j], ext[j+1])) < 0 {
			sum += cross(ext[j+1], ext[j])
			j++
		}
		u := sub(ext[i+1], ext[i])
		l1 := math.Hypot(u.x, u.y)
		w := sub(ext[j], ext[i])
		l2 := math.Hypot(w.x, w.y)
		h := cross(w, u) / l1
		l := math.Sqrt(l2*l2 - h*h)
		// point on edge at distance l from ext[i]
		t := l / l1
		res := pt{ext[i].x + u.x*t, ext[i].y + u.y*t}
		// area components
		area := cross(ext[j], res) + cross(res, ext[i+1]) - sum
		if area < *ans {
			*ans = area
		}
		// move window
		if i+2 < len(ext) {
			sum -= cross(ext[i+2], ext[i+1])
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pts := make([]pt, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}
	// extend points
	ext := make([]pt, 0, 2*n)
	ext = append(ext, pts...)
	ext = append(ext, pts...)
	// ensure orientation
	if n >= 3 {
		u := sub(ext[1], ext[0])
		v := sub(ext[2], ext[1])
		if cross(u, v) > 0 {
			// reverse all
			for i, j := 0, len(ext)-1; i < j; i, j = i+1, j-1 {
				ext[i], ext[j] = ext[j], ext[i]
			}
		}
	}
	ans := math.Inf(1)
	work(ext, n, &ans)
	// mirror and reverse
	for i := range ext {
		ext[i].x = -ext[i].x
	}
	for i, j := 0, len(ext)-1; i < j; i, j = i+1, j-1 {
		ext[i], ext[j] = ext[j], ext[i]
	}
	work(ext, n, &ans)
	// result
	fmt.Printf("%.10f\n", ans/2)
}
