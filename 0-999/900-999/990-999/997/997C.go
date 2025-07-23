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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	pow3 := make([]int64, n+1)
	fact[0] = 1
	pow3[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
		pow3[i] = pow3[i-1] * 3 % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % mod
	}
	pow3n := pow3[n]
	pow3nn := modPow(pow3n, int64(n))
	invPow3n := modPow(pow3n, mod-2)
	pow3nm1 := pow3[n-1]
	invPow3nm1 := modPow(pow3nm1, mod-2)
	base0 := (pow3nm1 - 1 + mod) % mod
	p0 := pow3n * modPow(base0, int64(n)) % mod
	var S int64
	invPow3nPow := int64(1)
	invPow3nm1Pow := int64(1)
	for a := 0; a <= n; a++ {
		comb := fact[n] * invFact[a] % mod * invFact[n-a] % mod
		var p int64
		if a == 0 {
			p = p0
		} else {
			invPow3nPow = invPow3nPow * invPow3n % mod
			invPow3nm1Pow = invPow3nm1Pow * invPow3nm1 % mod
			term1 := pow3nn * invPow3nm1Pow % mod
			base := (pow3[n-a] - 1 + mod) % mod
			term2 := 3 * modPow(base, int64(n)) % mod
			term3 := 3 * pow3nn % mod * invPow3nPow % mod
			p = (term1 + term2 - term3) % mod
		}
		add := comb * p % mod
		if a%2 == 0 {
			S = (S + add) % mod
		} else {
			S = (S - add) % mod
		}
	}
	if S < 0 {
		S += mod
	}
	ans := pow3nn - S
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
