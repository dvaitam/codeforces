package main

import (
	"fmt"
)

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := int64(3); i <= n/i; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int64
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if isPrime(n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
