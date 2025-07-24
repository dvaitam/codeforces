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
	var n, m int64
	var k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	limit := k
	if n < int64(limit) {
		limit = int(n)
	}
	invM := modPow(m%mod, mod-2)
	powM := make([]int64, limit+1)
	powM[0] = 1
	for i := 1; i <= limit; i++ {
		powM[i] = powM[i-1] * invM % mod
	}
	fall := make([]int64, limit+1)
	fall[0] = 1
	for i := 1; i <= limit; i++ {
		fall[i] = fall[i-1] * ((n - int64(i) + 1) % mod) % mod
	}
	S := make([]int64, limit+1)
	S[0] = 1
	for i := 1; i <= k; i++ {
		upper := i
		if upper > limit {
			upper = limit
		}
		for j := upper; j >= 1; j-- {
			S[j] = (S[j-1] + S[j]*int64(j)) % mod
		}
		S[0] = 0
	}
	var ans int64
	for j := 0; j <= limit; j++ {
		ans = (ans + S[j]*fall[j]%mod*powM[j]) % mod
	}
	fmt.Println(ans)
}
