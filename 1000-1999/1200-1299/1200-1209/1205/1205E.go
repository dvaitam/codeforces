package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	a %= MOD
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func mobius(n int) []int {
	mu := make([]int, n+1)
	mu[1] = 1
	primes := []int{}
	isComp := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if p*i > n {
				break
			}
			isComp[p*i] = true
			if i%p == 0 {
				mu[p*i] = 0
				break
			} else {
				mu[p*i] = -mu[i]
			}
		}
	}
	return mu
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	mu := mobius(n)
	invK := modPow(k%MOD, MOD-2)
	powInv := make([]int64, n+1)
	powInv[0] = 1
	for i := 1; i <= n; i++ {
		powInv[i] = powInv[i-1] * invK % MOD
	}
	res := int64(0)
	for t := 1; t < n; t++ {
		m := (n - 1) / t
		divisors := []int{}
		for d := 1; d*d <= t; d++ {
			if t%d == 0 {
				divisors = append(divisors, d)
				if d*d != t {
					divisors = append(divisors, t/d)
				}
			}
		}
		for _, g := range divisors {
			muVal := mu[t/g]
			if muVal == 0 {
				continue
			}
			contrib := int64(0)
			ng := n - g
			for r := 2; r <= 2*m; r++ {
				var cnt int64
				if r <= m+1 {
					cnt = int64(r - 1)
				} else {
					cnt = int64(2*m - r + 1)
				}
				val := 2*n - r*t
				if val > ng {
					val = ng
				}
				contrib = (contrib + cnt*powInv[val]) % MOD
			}
			if muVal == 1 {
				res = (res + contrib) % MOD
			} else {
				res = (res - contrib) % MOD
				if res < 0 {
					res += MOD
				}
			}
		}
	}
	res %= MOD
	if res < 0 {
		res += MOD
	}
	fmt.Println(res)
}
