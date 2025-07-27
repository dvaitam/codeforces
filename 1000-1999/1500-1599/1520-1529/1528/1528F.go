package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	// limit j up to min(k,n)
	m := int(k)
	if n < k {
		m = int(n)
	}
	// compute Stirling numbers of the second kind S[k][j] for 0<=j<=m
	S := make([][]int64, m+1)
	for i := range S {
		S[i] = make([]int64, m+1)
	}
	S[0][0] = 1
	for i := 1; i <= int(k); i++ {
		upper := i
		if upper > m {
			upper = m
		}
		S[i][0] = 0
		for j := 1; j <= upper; j++ {
			val := (S[i-1][j-1] + int64(j)*S[i-1][j]) % MOD
			S[i][j] = val
		}
	}

	// precompute falling factorial n^{(j)} and powers of n+1
	fall := make([]int64, m+1)
	fall[0] = 1
	for j := 1; j <= m; j++ {
		fall[j] = fall[j-1] * ((n - int64(j) + 1) % MOD) % MOD
	}

	powAll := modPow(n+1, n)
	invN1 := modPow(n+1, MOD-2)
	pow := powAll

	ans := int64(0)
	for j := 1; j <= m; j++ {
		pow = pow * invN1 % MOD // now pow = (n+1)^{n-j}
		term := S[k][j] * fall[j] % MOD
		term = term * pow % MOD
		ans = (ans + term) % MOD
	}
	fmt.Println(ans % MOD)
}
