package main

import (
	"fmt"
)

// isPrime checks if n is a prime number.
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n < 4 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func main() {
	var t int
	if _, err := fmt.Scan(&t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b int
		fmt.Scan(&a, &b)
		if isPrime(a+b) && a-b == 1 {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}
