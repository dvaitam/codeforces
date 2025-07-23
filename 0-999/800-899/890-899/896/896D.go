package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(a, e, mod int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, p, l, r int
	if _, err := fmt.Fscan(in, &n, &p, &l, &r); err != nil {
		return
	}
	mod := int64(p)
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2, mod)
	for i := n - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % mod
	}
	comb := func(nn, kk int) int64 {
		if kk < 0 || kk > nn {
			return 0
		}
		return fact[nn] * invFact[kk] % mod * invFact[nn-kk] % mod
	}

	ans := int64(0)
	for k := l; k <= r; k++ {
		maxJ := (n - k) / 2
		resK := int64(0)
		for j := 0; j <= maxJ; j++ {
			c1 := comb(n, k+2*j)
			c2 := comb(k+2*j, j)
			c3 := comb(k+2*j, j-1)
			v := (c2 - c3) % mod
			if v < 0 {
				v += mod
			}
			v = v * c1 % mod
			resK += v
			if resK >= mod {
				resK -= mod
			}
		}
		ans += resK
		if ans >= mod {
			ans -= mod
		}
	}

	fmt.Println(ans % mod)
}
