package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

func modPow(a, b int64) int64 {
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	fac := make([]int64, n+1)
	invFac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	invFac[n] = modPow(fac[n], mod-2)
	for i := n; i >= 1; i-- {
		invFac[i-1] = invFac[i] * int64(i) % mod
	}

	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = mod - int64(mod/int64(i))*inv[int(mod%int64(i))]%mod
	}

	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fac[a] * invFac[b] % mod * invFac[a-b] % mod
	}

	ans := int64(0)
	for k := 2; k <= n; k += 2 {
		m := n - k
		var ways int64
		if m == 0 {
			if n%2 == 0 {
				ways = 2
			} else {
				continue
			}
		} else {
			ways = 2 * int64(n) % mod
			ways = ways * comb(k-1, m-1) % mod
			ways = ways * inv[m] % mod
		}
		ans = (ans + ways*fac[k]) % mod
	}

	fmt.Fprintln(out, ans)
}
