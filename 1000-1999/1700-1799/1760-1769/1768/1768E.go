package main

import (
	"bufio"
	"fmt"
	"os"
)

func modPow(a, e, mod int64) int64 {
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
	var n int
	var mod int64
	fmt.Fscan(in, &n, &mod)
	max := 3 * n
	fact := make([]int64, max+1)
	invFact := make([]int64, max+1)
	fact[0] = 1
	for i := 1; i <= max; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[max] = modPow(fact[max], mod-2, mod)
	for i := max; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}

	factN := fact[n]
	fact2N := fact[2*n]
	fact3N := fact[3*n]
	pow3FactN := factN * factN % mod * factN % mod
	c2n := comb(2*n, n)
	countNoA := c2n * c2n % mod * pow3FactN % mod
	var s int64
	for k := 0; k <= n; k++ {
		term := comb(n, k) * comb(n, k) % mod * comb(2*n-k, n) % mod
		s = (s + term) % mod
	}
	intersection := s * pow3FactN % mod
	ans := (3*fact3N%mod - 2*fact2N%mod - 2*countNoA%mod + factN%mod + intersection%mod - 1) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Println(ans)
}
