package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

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
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	// factorials for binomial coefficients
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(a, b int) int64 {
		if b < 0 || b > a {
			return 0
		}
		return fact[a] * invFact[b] % mod * invFact[a-b] % mod
	}
	h := make([]int64, n+1)
	bArr := make([]int64, n+1)
	h[1] = 1
	bArr[1] = 1
	for i := 2; i <= n; i++ {
		var val int64
		for s := 1; s < i; s++ {
			k := i - s
			bk := h[k]
			if k != 1 {
				bk = bk * 2 % mod
			}
			val = (val + comb(i-1, s-1)*h[s]%mod*bk) % mod
		}
		h[i] = val % mod
		bArr[i] = 2 * h[i] % mod
	}
	ans := bArr[n] - 2
	if ans < 0 {
		ans += mod
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
