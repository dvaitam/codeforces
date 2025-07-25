package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	if k == 1 {
		fmt.Fprintln(out, 1)
		return
	}

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	comb := func(n, r int) int64 {
		if r < 0 || r > n {
			return 0
		}
		return fact[n] * invFact[r] % mod * invFact[n-r] % mod
	}

	ans := int64(0)
	kMod := k % mod
	km1Mod := (k - 1) % mod
	for r := 0; r <= n; r++ {
		cr := comb(n, r)
		for c := 0; c <= n; c++ {
			sign := 1
			if (r+c)%2 == 1 {
				sign = -1
			}
			term := cr * comb(n, c) % mod
			exp1 := int64(n*r + n*c - r*c)
			exp2 := int64((n - r) * (n - c))
			term = term * modPow(km1Mod, exp1) % mod
			term = term * modPow(kMod, exp2) % mod
			if sign == 1 {
				ans += term
			} else {
				ans -= term
			}
		}
	}
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(out, ans)
}
