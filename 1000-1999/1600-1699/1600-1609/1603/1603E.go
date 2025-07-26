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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var mod int64
	if _, err := fmt.Fscan(in, &n, &mod); err != nil {
		return
	}

	if n == 1 {
		fmt.Fprintln(out, 2%mod)
		return
	}

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2, mod)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}

	perm := func(x, y, z int) int64 {
		res := fact[n]
		res = res * invFact[x] % mod
		res = res * invFact[y] % mod
		res = res * invFact[z] % mod
		return res
	}

	ans := int64(0)
	// case with no (n-1)
	for cntN := 0; cntN <= n; cntN++ {
		cntNp1 := n - cntN
		ans = (ans + perm(0, cntN, cntNp1)) % mod
	}
	// cases with at least one (n-1)
	for z := 1; z <= n-1; z++ { // z = count of (n+1)
		for x := z + 1; x <= n-z; x++ { // x = count of (n-1)
			y := n - x - z
			ans = (ans + perm(x, y, z)) % mod
		}
	}

	fmt.Fprintln(out, ans%mod)
}
