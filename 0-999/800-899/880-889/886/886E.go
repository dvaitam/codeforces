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
	if k >= n {
		fmt.Println(0)
		return
	}
	// factorials and modular inverses up to n
	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv := make([]int64, n)
	inv[1] = 1
	for i := 2; i < n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[int(MOD%int64(i))]%MOD
	}

	w := make([]int64, k)
	w[0] = 1
	dpSum := int64(1) // dp[1]
	for i := 1; i < n; i++ {
		total := int64(0)
		for j := 0; j < k; j++ {
			total += w[j]
		}
		total %= MOD
		new0 := total * inv[i] % MOD
		for j := k - 1; j >= 1; j-- {
			w[j] = w[j-1]
		}
		w[0] = new0
		dpSum += w[0]
		if dpSum >= MOD {
			dpSum -= MOD
		}
	}
	good := fact[n-1] * dpSum % MOD
	bad := (fact[n] - good) % MOD
	if bad < 0 {
		bad += MOD
	}
	fmt.Println(bad)
}
