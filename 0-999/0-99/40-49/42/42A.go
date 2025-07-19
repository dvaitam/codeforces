package main

import (
	"fmt"
)

func main() {
	var n int
	var V float64
	if _, err := fmt.Scan(&n, &V); err != nil {
		return
	}
	a := make([]float64, n)
	b := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Scan(&b[i])
	}
	sumA := 0.0
	for i := 0; i < n; i++ {
		sumA += a[i]
	}
	// maximum x by pan volume
	maxXByPan := V / sumA
	// maximum x by ingredient limits
	maxXByIng := 1e18
	for i := 0; i < n; i++ {
		// a[i] >= 1 by problem constraints
		xi := b[i] / a[i]
		if xi < maxXByIng {
			maxXByIng = xi
		}
	}
	// choose the smaller x
	x := maxXByIng
	if maxXByPan < x {
		x = maxXByPan
	}
	volume := sumA * x
	// output with sufficient precision
	fmt.Printf("%.10f\n", volume)
}
