package main

import (
	"fmt"
	"math"
)

func main() {
	var a, b, m float64
	var vx, vy, vz float64

	if _, err := fmt.Scan(&a, &b, &m); err != nil {
		return
	}
	if _, err := fmt.Scan(&vx, &vy, &vz); err != nil {
		return
	}

	t := -m / vy

	X := a/2.0 + vx*t
	Z := vz * t

	Xrem := math.Mod(X, 2*a)
	if Xrem < 0 {
		Xrem += 2 * a
	}
	x0 := Xrem
	if Xrem > a {
		x0 = 2*a - Xrem
	}

	Zrem := math.Mod(Z, 2*b)
	if Zrem < 0 {
		Zrem += 2 * b
	}
	z0 := Zrem
	if Zrem > b {
		z0 = 2*b - Zrem
	}

	fmt.Printf("%.10f %.10f\n", x0, z0)
}
