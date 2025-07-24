package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	const maxN = 300000

	// precompute factorials and inverse factorials
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[maxN] = modPow(fact[maxN], MOD-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	// precompute involution counts: number of permutations with only 1- and 2-cycles
	invol := make([]int64, maxN+1)
	invol[0] = 1
	if maxN >= 1 {
		invol[1] = 1
	}
	for i := 2; i <= maxN; i++ {
		invol[i] = (invol[i-1] + int64(i-1)*invol[i-2]) % MOD
	}

	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var ans int64
		for s := 0; 4*s <= n; s++ {
			term := comb(n-2*s, 2*s)
			term = term * fact[2*s] % MOD * invFact[s] % MOD
			term = term * invol[n-4*s] % MOD
			ans = (ans + term) % MOD
		}
		fmt.Fprintln(out, ans)
	}
}
