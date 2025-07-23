package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, b int64) int64 {
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

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	if n > k {
		fmt.Println(0)
		return
	}

	// Precompute factorials and inverse factorials
	fac := make([]int64, k+1)
	ifac := make([]int64, k+1)
	fac[0] = 1
	for i := 1; i <= k; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[k] = modPow(fac[k], MOD-2)
	for i := k; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}

	comb := func(n, r int) int64 {
		if r < 0 || r > n {
			return 0
		}
		return fac[n] * ifac[r] % MOD * ifac[n-r] % MOD
	}

	pow2 := make([]int64, k+1)
	pow2[0] = 1
	for i := 1; i <= k; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}

	dp := make([]int64, k+1)
	dp[0] = 1
	for step := 1; step <= n; step++ {
		next := make([]int64, k+1)
		for used := 0; used <= k; used++ {
			val := dp[used]
			if val == 0 {
				continue
			}
			for add := 1; used+add <= k; add++ {
				ways := comb(k-used, add) * pow2[used] % MOD
				next[used+add] = (next[used+add] + val*ways) % MOD
			}
		}
		dp = next
	}

	var ans int64
	for _, v := range dp {
		ans = (ans + v) % MOD
	}
	fmt.Println(ans)
}
