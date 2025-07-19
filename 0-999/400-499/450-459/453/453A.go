package main

import (
	"fmt"
	"math"
)

func main() {
	var m, n int
	if _, err := fmt.Scan(&m, &n); err != nil {
		return
	}
	mf := float64(m)
	ans := mf
	for i := 1; i <= m; i++ {
		p := float64(i-1) / mf
		ans -= math.Pow(p, float64(n))
	}
	// Print answer with 6 decimal places
	fmt.Printf("%.6f\n", ans)
}
