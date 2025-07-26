package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	primes := sieve(31623) // sqrt(1e9) but ai<=1e7 -> sqrt=3162; use a bit bigger

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		counts := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			factorize(x, primes, counts)
		}
		pairs := 0
		leftover := 0
		for _, c := range counts {
			pairs += c / 2
			leftover += c % 2
		}
		result := pairs + leftover/3
		fmt.Fprintln(writer, result)
	}
}

func sieve(n int) []int {
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func factorize(x int, primes []int, counts map[int]int) {
	for _, p := range primes {
		if p*p > x {
			break
		}
		for x%p == 0 {
			counts[p]++
			x /= p
		}
	}
	if x > 1 {
		counts[x]++
	}
}
