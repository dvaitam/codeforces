package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

// modPow computes x^p mod MOD
func modPow(x, p int) int {
	res := 1
	x %= MOD
	for p > 0 {
		if p&1 == 1 {
			res = res * x % MOD
		}
		x = x * x % MOD
		p >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	fact := make([]int, n+2)
	inv := make([]int, n+2)
	fact[0] = 1
	inv[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * i % MOD
		inv[i] = modPow(fact[i], MOD-2)
	}

	// comb returns n choose k modulo MOD
	comb := func(n, k int) int {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * inv[k] % MOD * inv[n-k] % MOD
	}

	d := make([]int, n+2)
	sol := 0
	for i := n; i >= 1; i-- {
		if a[i] <= 0 {
			continue
		}
		ram := n - i
		if a[i] > ram {
			continue
		}
		d[i] = comb(ram, a[i])
		for j, nr := i+a[i]+1, a[i]; j <= n; j, nr = j+1, nr+1 {
			aux := comb(nr, a[i]) * d[j] % MOD
			d[i] += aux
			if d[i] >= MOD {
				d[i] -= MOD
			}
		}
		sol += d[i]
		if sol >= MOD {
			sol -= MOD
		}
	}
	fmt.Print(sol)
}
