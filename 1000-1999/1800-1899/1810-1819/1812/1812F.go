package main

import (
	"bufio"
	"fmt"
	"os"
)

// largestPrimeFactor returns the largest prime factor of n.
func largestPrimeFactor(n int64) int64 {
	maxPrime := int64(-1)

	for n%2 == 0 {
		maxPrime = 2
		n /= 2
	}

	for i := int64(3); i*i <= n; i += 2 {
		for n%i == 0 {
			maxPrime = i
			n /= i
		}
	}

	if n > 2 {
		maxPrime = n
	}

	return maxPrime
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Println(largestPrimeFactor(n))
}
