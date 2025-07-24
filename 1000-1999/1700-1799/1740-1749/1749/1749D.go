package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	// sieve for primes up to n
	isPrime := make([]bool, n+1)
	if n >= 2 {
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
	}

	mm := m % MOD

	// total number of arrays of length 1..n
	total := int64(0)
	pw := int64(1)
	for i := 1; i <= n; i++ {
		pw = pw * mm % MOD
		total = (total + pw) % MOD
	}

	// count arrays with a unique removal sequence
	prod := int64(1)
	ways := mm
	unamb := ways % MOD
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			prod *= int64(i)
		}
		if prod > m {
			break
		}
		q := m / prod
		ways = ways * (q % MOD) % MOD
		unamb = (unamb + ways) % MOD
	}

	ans := (total - unamb) % MOD
	if ans < 0 {
		ans += MOD
	}
	fmt.Println(ans)
}
