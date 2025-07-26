package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(base, exp int64) int64 {
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	maxN := n - 1
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[maxN] = modPow(fact[maxN], MOD-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	comb := func(n, r int) int64 {
		if r < 0 || r > n {
			return 0
		}
		return fact[n] * invFact[r] % MOD * invFact[n-r] % MOD
	}

	dp := make([]int64, maxN+1)
	dp[0] = 1
	for w := 1; w <= k; w++ {
		base := int64(k - w + 1)
		ndp := make([]int64, maxN+1)
		for t := 0; t <= maxN; t++ {
			if dp[t] == 0 {
				continue
			}
			remain := maxN - t
			valT := dp[t]
			for x := 0; x <= remain; x++ {
				exp := int64(x*t + x*(x-1)/2)
				val := valT * comb(remain, x) % MOD
				val = val * modPow(base, exp) % MOD
				ndp[t+x] = (ndp[t+x] + val) % MOD
			}
		}
		dp = ndp
	}
	fmt.Println(dp[maxN])
}
