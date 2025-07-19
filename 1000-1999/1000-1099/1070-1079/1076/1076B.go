package main

import (
	"bufio"
	"fmt"
	"os"
)

// isPrime returns true if n is a prime number (n >= 2).
func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// If n is prime, only one move possible
	if isPrime(n) {
		fmt.Println(1)
		return
	}
	// If n is even, can remove 1 each time: n/2 moves
	if n%2 == 0 {
		fmt.Println(n / 2)
		return
	}
	// n is odd and composite: find smallest divisor
	for i := int64(3); i <= n; i++ {
		if n%i == 0 {
			// Remove divisor i, remainder is n-i
			// Then we can remove 1 each move: (n-i)/2 + 1 moves total
			result := (n-i)/2 + 1
			fmt.Println(result)
			return
		}
	}
}
