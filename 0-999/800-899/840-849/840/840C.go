package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1_000_000_007

// squareFree returns the product of primes in x that appear
// an odd number of times in its factorization.
func squareFree(x int64) int64 {
	var res int64 = 1
	for p := int64(2); p*p <= x; p++ {
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		if cnt%2 == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sf := make([]int64, n)
	for i := 0; i < n; i++ {
		sf[i] = squareFree(a[i])
	}

	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			if i != j && sf[i] != sf[j] {
				adj[i][j] = true
			}
		}
	}

	full := 1 << n
	dp := make([][]int, full)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = 1
	}
	for mask := 0; mask < full; mask++ {
		for last := 0; last < n; last++ {
			val := dp[mask][last]
			if val == 0 {
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) == 0 && adj[last][next] {
					nm := mask | (1 << next)
					dp[nm][next] = (dp[nm][next] + val) % MOD
				}
			}
		}
	}
	ans := 0
	finalMask := full - 1
	for i := 0; i < n; i++ {
		ans = (ans + dp[finalMask][i]) % MOD
	}
	fmt.Println(ans)
}
