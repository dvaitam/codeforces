package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
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
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &h[i])
	}

	m := 0
	for i := 0; i < n; i++ {
		if h[i] != h[(i+1)%n] {
			m++
		}
	}

	if m == 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	fac := make([]int64, m+1)
	invFac := make([]int64, m+1)
	fac[0] = 1
	for i := 1; i <= m; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	invFac[m] = modPow(fac[m], mod-2)
	for i := m; i > 0; i-- {
		invFac[i-1] = invFac[i] * int64(i) % mod
	}

	pow := make([]int64, m+1)
	pow[0] = 1
	base := (k - 2) % mod
	if base < 0 {
		base += mod
	}
	for i := 1; i <= m; i++ {
		pow[i] = pow[i-1] * base % mod
	}

	comb := func(n, r int) int64 {
		if r < 0 || r > n {
			return 0
		}
		return fac[n] * invFac[r] % mod * invFac[n-r] % mod
	}

	var coef0 int64
	for j := 0; 2*j <= m; j++ {
		term := comb(m, 2*j) * comb(2*j, j) % mod * pow[m-2*j] % mod
		coef0 += term
		if coef0 >= mod {
			coef0 -= mod
		}
	}

	total := modPow(k%mod, int64(m))
	diff := (total - coef0) % mod
	if diff < 0 {
		diff += mod
	}
	inv2 := (mod + 1) / 2
	diff = diff * inv2 % mod

	ans := modPow(k%mod, int64(n-m)) * diff % mod
	fmt.Fprintln(writer, ans)
}
