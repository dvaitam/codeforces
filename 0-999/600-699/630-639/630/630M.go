package main

import (
	"fmt"
	"math"
)

func main() {
	var x int64
	if _, err := fmt.Scan(&x); err != nil {
		return
	}
	bestK := int64(0)
	bestDiff := int64(1<<62 - 1)
	for k := int64(0); k < 4; k++ {
		angle := -x + 90*k
		// normalize angle to range [-180,180]
		mod := ((angle % 360) + 360) % 360
		if mod > 180 {
			mod -= 360
		}
		diff := int64(math.Abs(float64(mod)))
		if diff < bestDiff {
			bestDiff = diff
			bestK = k
		}
	}
	fmt.Print(bestK)
}
