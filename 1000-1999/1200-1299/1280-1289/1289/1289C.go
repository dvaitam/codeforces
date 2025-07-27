package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if n <= 0 {
		fmt.Println(0)
		return
	}
	var x int
	// Read first element to initialize max
	fmt.Scan(&x)
	max := x
	// Iterate over remaining elements
	for i := 1; i < n; i++ {
		fmt.Scan(&x)
		if x > max {
			max = x
		}
	}
	fmt.Println(max)
}
