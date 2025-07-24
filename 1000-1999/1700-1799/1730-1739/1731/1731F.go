package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func sumPow(k int64, p int) int64 {
	m := p + 1
	y := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		y[i] = (y[i-1] + modPow(int64(i), int64(p))) % MOD
	}
	if k <= int64(m) {
		return y[k]
	}
	fact := make([]int64, m+1)
	ifact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	ifact[m] = modPow(fact[m], MOD-2)
	for i := m; i >= 1; i-- {
		ifact[i-1] = ifact[i] * int64(i) % MOD
	}
	pre := make([]int64, m+2)
	suf := make([]int64, m+2)
	pre[0] = 1
	for i := 0; i <= m; i++ {
		pre[i+1] = pre[i] * ((k - int64(i) + MOD) % MOD) % MOD
	}
	suf[m+1] = 1
	for i := m; i >= 0; i-- {
		suf[i] = suf[i+1] * ((k - int64(i) + MOD) % MOD) % MOD
	}
	ans := int64(0)
	for i := 0; i <= m; i++ {
		num := pre[i] * suf[i+1] % MOD
		den := ifact[i] * ifact[m-i] % MOD
		if (m-i)%2 == 1 {
			den = (MOD - den) % MOD
		}
		ans = (ans + y[i]*num%MOD*den) % MOD
	}
	return ans
}

func polyMul(a, b []int64) []int64 {
	c := make([]int64, len(a)+len(b)-1)
	for i := 0; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}
		for j := 0; j < len(b); j++ {
			if b[j] == 0 {
				continue
			}
			c[i+j] = (c[i+j] + a[i]*b[j]) % MOD
		}
	}
	return c
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
	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	// binom coefficients
	C := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int64, i+1)
		C[i][0] = 1
		C[i][i] = 1
		for j := 1; j < i; j++ {
			C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
		}
	}
	// power sums
	sumPowList := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		sumPowList[i] = sumPow(k, i)
	}
	powK := make([]int64, n+1)
	powK1 := make([]int64, n+1)
	powK[0] = 1
	powK1[0] = 1
	for i := 1; i <= n; i++ {
		powK[i] = powK[i-1] * (k % MOD) % MOD
		powK1[i] = powK1[i-1] * ((k + 1) % MOD) % MOD
	}
	ans := int64(0)
	for i := 1; i <= n; i++ {
		for l := 0; l <= i-1; l++ {
			for g := l + 1; g <= n-i; g++ {
				coef := C[i-1][l] * C[n-i][g] % MOD
				base := n - i - g + 1
				poly := make([]int64, base+1)
				poly[base] = 1
				// (v-1)^l
				a := make([]int64, l+1)
				for t := 0; t <= l; t++ {
					val := C[l][t]
					if (l-t)%2 == 1 {
						val = (MOD - val) % MOD
					}
					a[t] = val
				}
				poly = polyMul(poly, a)
				// (k-v+1)^{i-1-l}
				bdeg := i - 1 - l
				b := make([]int64, bdeg+1)
				for s := 0; s <= bdeg; s++ {
					val := C[bdeg][s] * powK1[bdeg-s] % MOD
					if s%2 == 1 {
						val = (MOD - val) % MOD
					}
					b[s] = val
				}
				poly = polyMul(poly, b)
				// (k-v)^g
				c := make([]int64, g+1)
				for r := 0; r <= g; r++ {
					val := C[g][r] * powK[g-r] % MOD
					if r%2 == 1 {
						val = (MOD - val) % MOD
					}
					c[r] = val
				}
				poly = polyMul(poly, c)
				sum := int64(0)
				for d := 0; d < len(poly); d++ {
					if poly[d] == 0 {
						continue
					}
					sum = (sum + poly[d]*sumPowList[d]) % MOD
				}
				ans = (ans + coef*sum) % MOD
			}
		}
	}
	fmt.Fprintln(out, ans%MOD)
}
