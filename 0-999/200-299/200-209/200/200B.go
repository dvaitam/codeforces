package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	sum := 0
	for i := 0; i < n; i++ {
		var a int
		fmt.Scan(&a)
		sum += a
	}
	avg := float64(sum) / float64(n)
	fmt.Printf("%f\n", avg)
}
