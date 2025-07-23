package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func modPow(a, e int) int {
	res := 1
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow(a, MOD-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}

	g := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &g[i])
	}

	L := make([]int, n)
	R := make([]int, n)
	maxN := n
	for i := 0; i < n; i++ {
		var s int
		var str string
		fmt.Fscan(in, &s, &str)
		R[i] = s
		cnt := 0
		for j := 0; j < len(str); j++ {
			if str[j] == '0' {
				cnt++
			}
		}
		L[i] = cnt
		if s > maxN {
			maxN = s
		}
	}

	// precompute factorials up to n (at most 5000)
	fact := make([]int, maxN+1)
	invFact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[maxN] = modInv(fact[maxN])
	for i := maxN - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * (i + 1) % MOD
	}
	comb := func(n, k int) int {
		if n < 0 || k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	cntTop := 0
	for i := 0; i < n; i++ {
		rank := 1
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if L[j] > R[i] {
				rank++
			}
		}
		if rank <= a {
			cntTop++
		}
	}

	ans := comb(cntTop, b)
	fmt.Fprintln(out, ans)
}
