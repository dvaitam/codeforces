package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func modPow(a, b int64) int64 {
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

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a0, x, y, k, M int64
	if _, err := fmt.Fscan(in, &n, &a0, &x, &y, &k, &M); err != nil {
		return
	}

	// compute L = lcm(1..k)
	L := int64(1)
	for i := int64(1); i <= k; i++ {
		L = L * i / gcd(L, i)
	}

	invN := modInv(n % MOD)
	mulStay := (MOD + 1 - invN) % MOD

	// dp arrays for residues
	dp := make([]int64, L)
	next := make([]int64, L)

	for step := k; step >= 1; step-- {
		s := int64(step)
		for r := int64(0); r < L; r++ {
			m := r % s
			// value when this residue is chosen at this step
			val := (invN*(r+dp[r-m]) + mulStay*dp[r]) % MOD
			next[r] = val
		}
		dp, next = next, dp
	}

	// contribution from the quotient part (multiples of L)
	constPart := (L % MOD) * (k % MOD) % MOD * invN % MOD

	ans := int64(0)
	val := a0
	for i := int64(0); i < n; i++ {
		q := val / L
		r := val % L
		ans += dp[r]
		ans += (q % MOD) * constPart % MOD
		ans %= MOD
		val = (val*x + y) % int64(M)
	}

	ans = ans * modPow(n%MOD, k) % MOD
	fmt.Println(ans)
}
