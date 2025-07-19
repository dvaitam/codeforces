package main

import (
	"fmt"
	"math"
)

func main() {
	var r, x1, y1, x2, y2 float64
	if _, err := fmt.Scan(&r, &x1, &y1, &x2, &y2); err != nil {
		return
	}
	const eps = 1e-15

	// If points coincide, choose any direction for new circle center
	if math.Abs(x1-x2) < eps && math.Abs(y1-y2) < eps {
		R := r / 2
		x0 := (2*x1 - r) / 2
		y0 := y1
		fmt.Printf("%.16f %.16f %.16f\n", x0, y0, R)
		return
	}
	// Distance between points
	dx := x2 - x1
	dy := y2 - y1
	dd := math.Hypot(dx, dy)
	// If second point is on or outside the circle, original circle suffices
	if dd >= r || math.Abs(dd-r) < eps {
		fmt.Printf("%.16f %.16f %.16f\n", x1, y1, r)
		return
	}
	// Compute new circle
	R := (r + dd) / 2
	// Move center from (x1,y1) towards (x2,y2)
	factor := (r - R) / dd
	x0 := (x1-x2)*factor + x1
	y0 := (y1-y2)*factor + y1
	fmt.Printf("%.16f %.16f %.16f\n", x0, y0, R)
}
