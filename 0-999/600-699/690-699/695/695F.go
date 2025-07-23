package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	// Smallest prime factor sieve
	spf := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			for j := i; j <= n; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	// Compute radical (product of distinct primes) for every number
	rad := make([]int, n+1)
	rad[1] = 1
	for i := 2; i <= n; i++ {
		p0 := spf[i]
		x := i / p0
		rad[i] = rad[x]
		if x%p0 != 0 {
			rad[i] *= p0
		}
	}

	freq := make([]int, n+1)
	for i := 1; i <= n; i++ {
		freq[rad[i]]++
	}
	fixed := make([]int, n+1)

	for i := 1; i <= n; i++ {
		if p[i] > 0 {
			if rad[i] != rad[p[i]] {
				fmt.Println(0)
				return
			}
			fixed[rad[i]]++
		}
	}

	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}

	ans := int64(1)
	for r := 1; r <= n; r++ {
		if freq[r] > 0 {
			m := freq[r] - fixed[r]
			if m < 0 {
				fmt.Println(0)
				return
			}
			ans = ans * fact[m] % MOD
		}
	}

	fmt.Println(ans)
}
