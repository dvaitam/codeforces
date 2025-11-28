package main

import (
	"fmt"
	"math"
)

const eps = 1e-4

func main() {
	var ax, ay, bx, by, cx, cy float64
	if _, err := fmt.Scan(&ax, &ay, &bx, &by, &cx, &cy); err != nil {
		return
	}

	a := math.Hypot(ax-bx, ay-by)
	b := math.Hypot(bx-cx, by-cy)
	c := math.Hypot(cx-ax, cy-ay)

	p := (a + b + c) / 2
	s := math.Sqrt(p * (p - a) * (p - b) * (p - c))

	// R = abc / 4S
	R := (a * b * c) / (4 * s)

	// Calculate central angles subtended by the three sides
	angles := make([]float64, 3)
	sides := []float64{a, b, c}
	for i, side := range sides {
		val := 1.0 - (side*side)/(2.0*R*R)
		if val < -1.0 {
			val = -1.0
		}
		if val > 1.0 {
			val = 1.0
		}
		angles[i] = math.Acos(val)
	}

	for n := 3; n <= 100; n++ {
		delta := 2.0 * math.Pi / float64(n)
		allInt := true
		for _, theta := range angles {
			quotient := theta / delta
			diff := math.Abs(quotient - math.Round(quotient))
			if diff > eps {
				allInt = false
				break
			}
		}
		if allInt {
			area := (float64(n) / 2.0) * R * R * math.Sin(2.0*math.Pi/float64(n))
			fmt.Printf("%.8f\n", area)
			return
		}
	}
}

