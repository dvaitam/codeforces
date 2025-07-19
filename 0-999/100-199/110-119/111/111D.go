package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

var c [1005][1005]int64
var f [1005]int64
var fz [2005]int64
var fm [1005]int64

// modExp computes base^exp % mod
func modExp(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

// compute binomial coefficients up to n
func getC(n int) {
	for i := 0; i <= n; i++ {
		c[i][0] = 1
		for j := 1; j <= i; j++ {
			c[i][j] = c[i-1][j] + c[i-1][j-1]
			if c[i][j] >= mod {
				c[i][j] -= mod
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	// special case
	if m == 1 {
		fmt.Println(modExp(int64(k), int64(n)))
		return
	}
	getC(n)
	// compute f[i]
	for i := 1; i <= n; i++ {
		f[i] = modExp(int64(i), int64(n))
		for j := 1; j < i; j++ {
			f[i] = (f[i] - f[j]*c[i][j]%mod + mod) % mod
		}
	}
	fz[0], fm[0] = 1, 1
	limit := 2 * n
	if limit > k {
		limit = k
	}
	for i := 1; i <= limit; i++ {
		fz[i] = fz[i-1] * int64(k-i+1) % mod
	}
	for i := 1; i <= n; i++ {
		fm[i] = fm[i-1] * modExp(int64(i), mod-2) % mod
	}
	var ans int64
	for i := 0; i <= n && i <= k; i++ {
		temp := modExp(int64(i), int64(m-2)*int64(n))
		for j := 0; j <= n-i && i+2*j <= k; j++ {
			idx := i + 2*j
			term := temp * fz[idx] % mod
			term = term * fm[i] % mod
			term = term * fm[j] % mod
			term = term * fm[j] % mod
			val := f[i+j]
			term = term * val % mod
			term = term * val % mod
			ans = (ans + term) % mod
		}
	}
	fmt.Println(ans)
}
