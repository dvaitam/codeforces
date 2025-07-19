package main

import (
	"fmt"
)

func isPrime(x int64) bool {
	if x < 2 {
		return false
	}
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
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
	var a, b int64
	// find largest prime a <= n
	for i := n; i >= 2; i-- {
		if isPrime(i) {
			a = i
			break
		}
	}
	// find b such that either a+b == n or n-(a+b) is prime
	for i := int64(2); i <= n; i++ {
		if !isPrime(i) {
			continue
		}
		sum := a + i
		if sum != n && !isPrime(n-sum) {
			continue
		}
		b = i
		break
	}
	// output
	if a == n {
		fmt.Println(1)
		fmt.Println(a)
	} else if a+b == n {
		fmt.Println(2)
		fmt.Printf("%d %d\n", a, b)
	} else {
		fmt.Println(3)
		fmt.Printf("%d %d %d\n", a, b, n-(a+b))
	}
}
