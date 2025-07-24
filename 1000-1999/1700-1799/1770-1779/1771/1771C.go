package main

import (
	"bufio"
	"fmt"
	"os"
)

// sieve returns all primes up to n using the simple sieve of Eratosthenes.
func sieve(n int) []int {
	mark := make([]bool, n+1)
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if !mark[i] {
			primes = append(primes, i)
			if i*i <= n {
				for j := i * i; j <= n; j += i {
					mark[j] = true
				}
			}
		}
	}
	return primes
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	primes := sieve(31623) // sqrt(1e9) ~= 31623

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		found := false
		mp := make(map[int]bool)
	outer:
		for _, val := range a {
			x := val
			for _, p := range primes {
				if p*p > x {
					break
				}
				if x%p == 0 {
					if mp[p] {
						found = true
						break outer
					}
					mp[p] = true
					for x%p == 0 {
						x /= p
					}
				}
			}
			if x > 1 {
				if mp[x] {
					found = true
					break outer
				}
				mp[x] = true
			}
		}

		if found {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
