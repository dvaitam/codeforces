package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	// Build representation with maximum number of primes (2 and 3)
	primes := []int{}
	if n%2 == 1 {
		primes = append(primes, 3)
		n -= 3
	}
	for n > 0 {
		primes = append(primes, 2)
		n -= 2
	}
	// Output result
	fmt.Println(len(primes))
	for i, p := range primes {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(p)
	}
	fmt.Println()
}
