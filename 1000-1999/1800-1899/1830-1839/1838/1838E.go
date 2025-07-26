package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007
const MAXN = 200000

var inv [MAXN + 2]int64

func init() {
	inv[1] = 1
	for i := 2; i <= MAXN+1; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[MOD%int64(i)]%MOD
	}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m, k int64
		fmt.Fscan(in, &n, &m, &k)
		for i := int64(0); i < n; i++ {
			var tmp int64
			fmt.Fscan(in, &tmp)
		}
		if k == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		total := modPow(k, m)
		comb := int64(1) // C(m,0)
		powTerm := modPow(k-1, m)
		sum := comb * powTerm % MOD
		invK1 := modPow(k-1, MOD-2)
		limit := n - 1
		if limit > m {
			limit = m
		}
		for i := int64(1); i <= limit; i++ {
			comb = comb * ((m - i + 1) % MOD) % MOD
			comb = comb * inv[i] % MOD
			powTerm = powTerm * invK1 % MOD
			sum = (sum + comb*powTerm%MOD) % MOD
		}
		ans := (total - sum) % MOD
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(out, ans)
	}
}
