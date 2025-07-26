package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// Compute Fibonacci number iteratively
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	fmt.Println(a)
}
