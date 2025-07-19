package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// isPrime checks if n is a prime number.
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	limit := int(math.Sqrt(float64(n))) + 1
	for i := 2; i < limit; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// nextPrime returns the smallest prime >= n.
func nextPrime(n int) int {
	for !isPrime(n) {
		n++
	}
	return n
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var N int
	if _, err := fmt.Fscan(reader, &N); err != nil {
		return
	}
	M := nextPrime(N)
	// print total edges
	fmt.Fprintln(writer, M)
	// print cycle edges
	for i := 1; i < N; i++ {
		fmt.Fprintf(writer, "%d %d\n", i, i+1)
	}
	fmt.Fprintf(writer, "1 %d\n", N)
	// add extra edges
	half := N / 2
	extra := M - N
	for i := 1; i <= extra; i++ {
		fmt.Fprintf(writer, "%d %d\n", i, i+half)
	}
}
