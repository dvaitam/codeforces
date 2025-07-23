package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	fmt.Fscan(in, &n)
	// compute C(2n+2, n+1) - 1 modulo mod
	m := 2*n + 2
	fact := make([]int64, m+1)
	fact[0] = 1
	for i := int64(1); i <= m; i++ {
		fact[i] = fact[i-1] * i % mod
	}
	invFact := make([]int64, m+1)
	invFact[m] = modPow(fact[m], mod-2)
	for i := m; i > 0; i-- {
		invFact[i-1] = invFact[i] * i % mod
	}
	choose := func(n, k int64) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}
	ans := (choose(m, n+1) - 1) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
