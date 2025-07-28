package main

import (
	"fmt"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	var a, b int
	fmt.Scan(&a, &b)
	fmt.Println(gcd(a, b))
}
