package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func sieve(limit int) []int {
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

func primeFactors(x int64, primes []int) []int64 {
	factors := make([]int64, 0)
	for _, p := range primes {
		if int64(p)*int64(p) > x {
			break
		}
		if x%int64(p) == 0 {
			factors = append(factors, int64(p))
			for x%int64(p) == 0 {
				x /= int64(p)
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return factors
}

func countCoprime(limit int64, factors []int64) int64 {
	if limit <= 0 {
		return 0
	}
	m := len(factors)
	var bad int64
	for mask := 1; mask < (1 << m); mask++ {
		prod := int64(1)
		bits := 0
		for i := 0; i < m; i++ {
			if (mask>>i)&1 == 1 {
				prod *= factors[i]
				bits++
			}
		}
		if bits%2 == 1 {
			bad += limit / prod
		} else {
			bad -= limit / prod
		}
	}
	return limit - bad
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	primes := sieve(31623)

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var m int64
		fmt.Fscan(reader, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := int64(1)
		ok := true
		for i := 1; i < n && ok; i++ {
			if a[i-1]%a[i] != 0 {
				ok = false
				break
			}
			x := a[i-1] / a[i]
			limit := m / a[i]
			pf := primeFactors(x, primes)
			cnt := countCoprime(limit, pf)
			ans = (ans * (cnt % mod)) % mod
		}
		if ok {
			fmt.Fprintln(writer, ans%mod)
		} else {
			fmt.Fprintln(writer, 0)
		}
	}
}
