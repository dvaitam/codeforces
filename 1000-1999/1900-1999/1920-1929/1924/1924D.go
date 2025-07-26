package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007
const MAX int = 4000

var fac, inv [MAX + 1]int64

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
func C(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fac[n] * inv[r] % MOD * inv[n-r] % MOD
}
func F(n, m, U int) int64 {
	if m > n+U {
		return 0
	}
	res := (C(n+m, n) - C(n+m, n+U+1)) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}
func solve(n, m, k int) int64 {
	if k > n || k > m {
		return 0
	}
	U := m - k
	if U == 0 {
		return F(n, m, 0)
	}
	res := (F(n, m, U) - F(n, m, U-1)) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}
func main() {
	fac[0] = 1
	for i := 1; i <= MAX; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	inv[MAX] = modPow(fac[MAX], MOD-2)
	for i := MAX; i >= 1; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		fmt.Fprintln(out, solve(n, m, k))
	}
}
