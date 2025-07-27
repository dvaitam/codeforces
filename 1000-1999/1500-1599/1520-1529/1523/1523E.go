package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

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

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modInv(fact[n])
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	return fact, invFact
}

func comb(n, r int, fact, invFact []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * invFact[r] % MOD * invFact[n-r] % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	const maxN = 100000
	fact, invFact := prepareFact(maxN)

	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)

		maxM := (n + k - 1) / k
		ans := int64(1)
		for m := 1; m <= maxM; m++ {
			N := n - (k-1)*(m-1)
			if N < m {
				break
			}
			num := comb(N, m, fact, invFact)
			den := comb(n, m, fact, invFact)
			val := num * modInv(den) % MOD
			ans += val
			if ans >= MOD {
				ans -= MOD
			}
		}
		fmt.Fprintln(out, ans%MOD)
	}
}
