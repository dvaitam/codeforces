package main

import (
	"bufio"
	"fmt"
	"os"
)

// sieve primes up to limit
func sieve(limit int) []int {
	isComp := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func countDistinctPrimeFactors(n int, primes []int) int {
	cnt := 0
	for _, p := range primes {
		if p*p > n {
			break
		}
		if n%p == 0 {
			cnt++
			for n%p == 0 {
				n /= p
			}
		}
	}
	if n > 1 {
		cnt++
	}
	return cnt
}

func countPairs(c, d, x, g int, primes []int) int {
	y := x / g
	if (y+d)%c != 0 {
		return 0
	}
	k := (y + d) / c
	if k <= 0 {
		return 0
	}
	cnt := countDistinctPrimeFactors(k, primes)
	return 1 << cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const LIMIT = 50000
	primes := sieve(LIMIT)

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var c, d, x int
		fmt.Fscan(in, &c, &d, &x)
		ans := 0
		for g := 1; g*g <= x; g++ {
			if x%g != 0 {
				continue
			}
			ans += countPairs(c, d, x, g, primes)
			if g*g != x {
				ans += countPairs(c, d, x, x/g, primes)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
