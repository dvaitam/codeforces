package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const eps = 1e-9

// cut returns the polygon (x,y) clipped by the half-plane aa*x+bb*y+cc <= 0
func cut(x, y []float64, aa, bb, cc float64) ([]float64, []float64) {
	n := len(x)
	if n == 0 {
		return x, y
	}
	var nx, ny []float64
	for i := 0; i < n; i++ {
		xi, yi := x[i], y[i]
		zi := aa*xi + bb*yi + cc
		// include point if inside
		if zi < eps {
			nx = append(nx, xi)
			ny = append(ny, yi)
		}
		// edge to next
		j := i + 1
		if j == n {
			j = 0
		}
		xj, yj := x[j], y[j]
		zj := aa*xj + bb*yj + cc
		if (zi < -eps && zj > eps) || (zi > eps && zj < -eps) {
			// compute intersection
			a := yj - yi
			b := xi - xj
			c := -a*xi - b*yi
			d := a*bb - b*aa
			ix := (b*cc - c*bb) / d
			iy := (c*aa - a*cc) / d
			nx = append(nx, ix)
			ny = append(ny, iy)
		}
	}
	return nx, ny
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var w, h, n int
	if _, err := fmt.Fscan(in, &w, &h, &n); err != nil {
		return
	}
	pts := make([]struct{ x, y int }, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	// combine duplicates
	xs := make([]struct{ x, y int }, 0, n)
	ws := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if i > 0 && pts[i] == pts[i-1] {
			ws[len(ws)-1]++
		} else {
			xs = append(xs, pts[i])
			ws = append(ws, 1)
		}
	}
	pts = xs
	n = len(pts)
	ansSq := 0.0
	// for each site
	for st := 0; st < n; st++ {
		// start with rectangle
		x := []float64{0, float64(w), float64(w), 0}
		y := []float64{0, 0, float64(h), float64(h)}
		// cut by all other bisectors
		for i := 0; i < n; i++ {
			dx := float64(pts[i].x - pts[st].x)
			dy := float64(pts[i].y - pts[st].y)
			midx := float64(pts[i].x+pts[st].x) * 0.5
			midy := float64(pts[i].y+pts[st].y) * 0.5
			cc := -dx*midx - dy*midy
			x, y = cut(x, y, dx, dy, cc)
			if len(x) == 0 {
				break
			}
		}
		// if multiplicity >=2, take farthest vertex
		if ws[st] >= 2 {
			for i := 0; i < len(x); i++ {
				dx := x[i] - float64(pts[st].x)
				dy := y[i] - float64(pts[st].y)
				d := dx*dx + dy*dy
				if d > ansSq {
					ansSq = d
				}
			}
			continue
		}
		// find nearest neighbor for each edge midpoint
		m := len(x)
		ptsIdx := make([]int, m)
		for i := 0; i < m; i++ {
			j := (i + 1) % m
			mx := (x[i] + x[j]) * 0.5
			my := (y[i] + y[j]) * 0.5
			best := -1
			bestD := 1e100
			for k := 0; k < n; k++ {
				if k == st {
					continue
				}
				dx := mx - float64(pts[k].x)
				dy := my - float64(pts[k].y)
				d := dx*dx + dy*dy
				if d < bestD {
					bestD = d
					best = k
				}
			}
			ptsIdx[i] = best
		}
		// save original polygon
		ox := make([]float64, len(x))
		oy := make([]float64, len(y))
		copy(ox, x)
		copy(oy, y)
		sz := len(ptsIdx)
		// for each candidate neighbor
		for jj := 0; jj < sz; jj++ {
			pt := ptsIdx[jj]
			// reset polygon
			x = make([]float64, len(ox))
			y = make([]float64, len(oy))
			copy(x, ox)
			copy(y, oy)
			// cut by all other neighbors except this
			for u := 0; u < sz; u++ {
				if u == jj {
					continue
				}
				i := ptsIdx[u]
				dx := float64(pts[i].x - pts[pt].x)
				dy := float64(pts[i].y - pts[pt].y)
				midx := float64(pts[i].x+pts[pt].x) * 0.5
				midy := float64(pts[i].y+pts[pt].y) * 0.5
				cc := -dx*midx - dy*midy
				x, y = cut(x, y, dx, dy, cc)
				if len(x) == 0 {
					break
				}
			}
			// update answer for this region
			for i := 0; i < len(x); i++ {
				dx := x[i] - float64(pts[pt].x)
				dy := y[i] - float64(pts[pt].y)
				d := dx*dx + dy*dy
				if d > ansSq {
					ansSq = d
				}
			}
		}
	}
	// output
	fmt.Printf("%.17f\n", math.Sqrt(ansSq))
}
