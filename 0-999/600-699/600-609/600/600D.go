package main

import (
	"fmt"
	"math"
)

// compute returns the area of intersection of two circles
func compute(x1, y1, r1, x2, y2, r2 float64) float64 {
	if r1 < r2 {
		return compute(x2, y2, r2, x1, y1, r1)
	}
	dx := x1 - x2
	dy := y1 - y2
	dsq := dx*dx + dy*dy
	// one circle is completely inside the other
	if dsq <= (r1-r2)*(r1-r2) {
		return math.Pi * r2 * r2
	}
	// circles are separate
	if dsq >= (r1+r2)*(r1+r2) {
		return 0.0
	}
	d := math.Sqrt(dsq)
	// angles for circular segment
	cos1 := (r1*r1 + dsq - r2*r2) / (2 * r1 * d)
	cos2 := (r2*r2 + dsq - r1*r1) / (2 * r2 * d)
	theta1 := math.Acos(cos1)
	theta2 := math.Acos(cos2)
	// area of segments
	p1 := r1 * r1 * (theta1 - cos1*math.Sin(theta1))
	p2 := r2 * r2 * (theta2 - cos2*math.Sin(theta2))
	return p1 + p2
}

func main() {
	var x1, y1, r1, x2, y2, r2 float64
	if _, err := fmt.Scan(&x1, &y1, &r1, &x2, &y2, &r2); err != nil {
		return
	}
	res := compute(x1, y1, r1, x2, y2, r2)
	fmt.Printf("%.20f\n", res)
}
