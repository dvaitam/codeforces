package main

import (
	"fmt"
)

func main() {
	var x1, y1, x2, y2, x3, y3 int64
	if _, err := fmt.Scan(&x1, &y1, &x2, &y2, &x3, &y3); err != nil {
		return
	}
	// Compute parallelogram completion points
	candidates := [][2]int64{
		{x1 + x2 - x3, y1 + y2 - y3},
		{x2 + x3 - x1, y2 + y3 - y1},
		{x1 + x3 - x2, y1 + y3 - y2},
	}
	uniq := make(map[[2]int64]bool)
	result := make([][2]int64, 0, 3)
	for _, p := range candidates {
		if !uniq[p] {
			uniq[p] = true
			result = append(result, p)
		}
	}
	// Output
	fmt.Println(len(result))
	for _, p := range result {
		fmt.Println(p[0], p[1])
	}
}
