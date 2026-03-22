package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b%2 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		b /= 2
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int64
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}

	if m*(k+1) > n {
		fmt.Println(0)
		return
	}

	D := n - m*(k+1)

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	invFact[0] = 1
	for i := int64(1); i <= n; i++ {
		fact[i] = (fact[i-1] * i) % MOD
	}
	invFact[n] = powMod(fact[n], MOD-2)
	for i := n - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * (i + 1)) % MOD
	}

	nCr := func(n, r int64) int64 {
		if r < 0 || r > n {
			return 0
		}
		return fact[n] * invFact[r] % MOD * invFact[n-r] % MOD
	}

	ans := int64(0)
	limit := D / k
	if m < limit {
		limit = m
	}

	for j := int64(0); j <= limit; j++ {
		term := nCr(m, j)
		term = (term * powMod(2, m-j)) % MOD
		term = (term * nCr(D-j*k+m, m)) % MOD

		if j%2 == 1 {
			ans = (ans - term + MOD) % MOD
		} else {
			ans = (ans + term) % MOD
		}
	}

	fmt.Println(ans)
}
