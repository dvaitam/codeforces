package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
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

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		a[i] %= MOD
	}
	if n == 1 {
		if a[0] < 0 {
			a[0] += MOD
		}
		fmt.Fprintln(writer, a[0]%MOD)
		return
	}
	if n%2 == 1 {
		b := make([]int64, n-1)
		add := true
		for i := 0; i < n-1; i++ {
			if add {
				b[i] = (a[i] + a[i+1]) % MOD
			} else {
				b[i] = (a[i] - a[i+1]) % MOD
				if b[i] < 0 {
					b[i] += MOD
				}
			}
			add = !add
		}
		a = b
		n--
	}

	m := n / 2
	fact := make([]int64, m)
	invFact := make([]int64, m)
	fact[0] = 1
	for i := 1; i < m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[m-1] = modInv(fact[m-1])
	for i := m - 2; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	ans := int64(0)
	for j := 0; j < m; j++ {
		coef := comb(m-1, j)
		var val int64
		if m%2 == 0 {
			val = (a[2*j] - a[2*j+1]) % MOD
		} else {
			val = (a[2*j] + a[2*j+1]) % MOD
		}
		if val < 0 {
			val += MOD
		}
		ans = (ans + coef*val) % MOD
	}
	fmt.Fprintln(writer, ans%MOD)
}
