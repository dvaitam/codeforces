package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	a %= MOD
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	m := k + 2
	f := make([]int64, m+1)
	for i := int64(1); i <= m; i++ {
		f[i] = (f[i-1] + modPow(i, k)) % MOD
	}

	if n <= m {
		fmt.Fprintln(writer, f[n])
		return
	}

	fact := make([]int64, m+1)
	invFact := make([]int64, m+1)
	fact[0] = 1
	for i := int64(1); i <= m; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[m] = modPow(fact[m], MOD-2)
	for i := m; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % MOD
	}

	pre := make([]int64, m+2)
	suf := make([]int64, m+2)
	pre[0] = 1
	x := n % MOD
	for i := int64(1); i <= m; i++ {
		pre[i] = pre[i-1] * ((x - i + MOD) % MOD) % MOD
	}
	suf[m+1] = 1
	for i := m; i >= 1; i-- {
		suf[i] = suf[i+1] * ((x - i + MOD) % MOD) % MOD
	}

	ans := int64(0)
	for i := int64(1); i <= m; i++ {
		numerator := pre[i-1] * suf[i+1] % MOD
		term := f[i] * numerator % MOD
		term = term * invFact[i-1] % MOD
		term = term * invFact[m-i] % MOD
		if (m-i)%2 == 1 {
			term = MOD - term
		}
		ans = (ans + term) % MOD
	}

	fmt.Fprintln(writer, ans)
}
