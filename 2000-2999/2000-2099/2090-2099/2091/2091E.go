package main

import (
	"bufio"
	"fmt"
	"os"
)

func sieve(limit int) []int {
	if limit < 2 {
		return nil
	}
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	ns := make([]int, t)
	maxN := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(in, &ns[i])
		if ns[i] > maxN {
			maxN = ns[i]
		}
	}

	primes := sieve(maxN)

	for _, n := range ns {
		var ans int64
		for _, p := range primes {
			if p > n {
				break
			}
			ans += int64(n / p)
		}
		fmt.Fprintln(out, ans)
	}
}
