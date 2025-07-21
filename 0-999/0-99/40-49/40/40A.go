package main

import (
	"fmt"
	"math"
)

func main() {
	var x, y int
	if _, err := fmt.Scanf("%d %d", &x, &y); err != nil {
		return
	}
	d2 := x*x + y*y
	// Check if on integer circle boundary
	fs := math.Sqrt(float64(d2))
	if math.Abs(fs-math.Round(fs)) < 1e-9 {
		fmt.Println("black")
		return
	}
	// Determine ring by floor distance
	k := int(math.Floor(fs))
	if k%2 == 0 {
		fmt.Println("white")
	} else {
		fmt.Println("black")
	}
}
