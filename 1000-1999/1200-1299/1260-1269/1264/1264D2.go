package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	base := a % mod
	for b > 0 {
		if b&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		b >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)

	// precompute factorials and powers of two
	fac := make([]int64, n+2)
	invf := make([]int64, n+2)
	pow2 := make([]int64, n+2)
	fac[0] = 1
	pow2[0] = 1
	for i := 1; i <= n+1; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
		pow2[i] = pow2[i-1] * 2 % mod
	}
	invf[n+1] = modPow(fac[n+1], mod-2)
	for i := n; i >= 0; i-- {
		invf[i] = invf[i+1] * int64(i+1) % mod
	}

	comb := func(n, k int64) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * invf[k] % mod * invf[n-k] % mod
	}

	preOpen := make([]int, n+1)
	preQ := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preOpen[i] = preOpen[i-1]
		preQ[i] = preQ[i-1]
		switch s[i-1] {
		case '(':
			preOpen[i]++
		case '?':
			preQ[i]++
		}
	}

	sufClose := make([]int, n+2)
	sufQ := make([]int, n+2)
	for i := n; i >= 1; i-- {
		sufClose[i] = sufClose[i+1]
		sufQ[i] = sufQ[i+1]
		switch s[i-1] {
		case ')':
			sufClose[i]++
		case '?':
			sufQ[i]++
		}
	}

	ans := int64(0)
	for i := 1; i <= n; i++ {
		if s[i-1] != '(' && s[i-1] != '?' {
			continue
		}
		Ai := int64(preQ[i-1])
		Oi := int64(preOpen[i-1])
		for j := i + 1; j <= n; j++ {
			if s[j-1] != ')' && s[j-1] != '?' {
				continue
			}
			midQ := preQ[j-1] - preQ[i]
			Bj := int64(sufQ[j+1])
			Cj := int64(sufClose[j+1])
			N := Ai + Bj
			R := Ai + Oi - Cj
			if R < 0 || R > N {
				continue
			}
			ways := comb(N, R)
			ways = ways * pow2[midQ] % mod
			ans += ways
			if ans >= mod {
				ans -= mod
			}
		}
	}
	fmt.Println(ans % mod)
}
