package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	N := 3 * (n + 1)
	fac := make([]int64, N+1)
	ifac := make([]int64, N+1)
	fac[0] = 1
	for i := 1; i <= N; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[N] = modPow(fac[N], MOD-2)
	for i := N; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
	}

	inv3 := modPow(3, MOD-2)
	ans := make([]int64, 3*n+1)
	ans[0] = int64(n+1) % MOD
	for k := 1; k <= 3*n; k++ {
		val := comb(N, k+1)
		val = (val - 3*ans[k-1]) % MOD
		if val < 0 {
			val += MOD
		}
		if k >= 2 {
			val = (val - ans[k-2]) % MOD
			if val < 0 {
				val += MOD
			}
		}
		ans[k] = val * inv3 % MOD
	}

	for i := 0; i < q; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x >= 0 && x < len(ans) {
			fmt.Fprintln(writer, ans[x])
		} else {
			fmt.Fprintln(writer, 0)
		}
	}
}
