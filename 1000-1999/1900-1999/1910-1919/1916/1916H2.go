package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= MOD
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var p int64
	var k int
	fmt.Fscan(in, &n, &p, &k)

	powQ := make([]int64, k+2)
	powQ[0] = 1
	for i := 1; i <= k+1; i++ {
		powQ[i] = powQ[i-1] * (p % MOD) % MOD
	}
	qn := powMod(p%MOD, n)

	ans := make([]int64, k+1)
	ans[0] = 1

	limit := k
	if n < int64(limit) {
		limit = int(n)
	}

	for r := 1; r <= limit; r++ {
		diff := qn - powQ[r-1]
		diff %= MOD
		if diff < 0 {
			diff += MOD
		}
		val := diff * diff % MOD
		val = val * ans[r-1] % MOD
		den1 := powQ[r] - powQ[r-1]
		den1 %= MOD
		if den1 < 0 {
			den1 += MOD
		}
		den2 := powQ[r] - 1
		den2 %= MOD
		if den2 < 0 {
			den2 += MOD
		}
		val = val * (p - 1) % MOD
		val = val * powMod(den1, MOD-2) % MOD
		val = val * powMod(den2, MOD-2) % MOD
		ans[r] = val
	}

	for r := limit + 1; r <= k; r++ {
		ans[r] = 0
	}

	for i := 0; i <= k; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
