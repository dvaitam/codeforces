package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

var inv2 int64 = modPow(2, mod-2)

func solve(n, k, x int64) int64 {
	if x == 0 {
		if n > k {
			return 0
		}
		pow2k := modPow(2, k)
		pow := int64(1)
		ans := int64(1)
		for i := int64(0); i < n; i++ {
			ans = ans * ((pow2k - pow + mod) % mod) % mod
			pow = pow * 2 % mod
		}
		return ans
	}
	if k == 0 {
		return 0
	}
	m := n
	if k < m {
		m = k
	}
	pow2n := modPow(2, n)
	pow2r := int64(1)
	pow2rMinus1 := int64(1)
	pow2kRPlus1 := modPow(2, k+1)
	prod := int64(1)
	g1 := int64(1)
	g0 := int64(0)
	ans := int64(0)
	for r := int64(0); r <= m; r++ {
		if r > 0 {
			term := (pow2n - pow2rMinus1) % mod
			if term < 0 {
				term += mod
			}
			prod = prod * term % mod
			pow2kRPlus1 = pow2kRPlus1 * inv2 % mod
			g1 = g1 * ((pow2kRPlus1 - 1 + mod) % mod) % mod
			g1 = g1 * modPow((pow2r-1+mod)%mod, mod-2) % mod
			if r == 1 {
				g0 = 1
			} else {
				g0 = g0 * ((pow2kRPlus1 - 1 + mod) % mod) % mod
				g0 = g0 * modPow((pow2rMinus1-1+mod)%mod, mod-2) % mod
			}
		}
		sub := g1 - g0
		if sub < 0 {
			sub += mod
		}
		ans = (ans + prod*sub) % mod
		pow2rMinus1 = pow2r
		pow2r = pow2r * 2 % mod
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t++ {
		var n, k, x int64
		fmt.Fscan(in, &n, &k, &x)
		fmt.Fprintln(out, solve(n, k, x))
	}
}
