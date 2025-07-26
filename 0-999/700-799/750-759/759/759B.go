package main

import "fmt"

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a % b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	var a, b int
	// read two integers
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	// output gcd
	fmt.Println(gcd(a, b))
}
