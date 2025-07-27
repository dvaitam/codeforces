package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}

	rank := make([]int, n+1)
	for i := 0; i < n; i++ {
		rank[p[i]] = i
	}
	rank[n] = -1

	strict := 0
	for idx := 0; idx < n-1; idx++ {
		i, j := p[idx], p[idx+1]
		ri1 := -1
		if i+1 < n {
			ri1 = rank[i+1]
		}
		rj1 := -1
		if j+1 < n {
			rj1 = rank[j+1]
		}
		if ri1 > rj1 {
			strict++
		}
	}

	if k-strict <= 0 {
		fmt.Fprintln(out, 0)
		return
	}

	N := int64(n + k - strict - 1)
	fact := make([]int64, N+1)
	invFact := make([]int64, N+1)
	fact[0] = 1
	for i := int64(1); i <= N; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[N] = modPow(fact[N], MOD-2)
	for i := N; i > 0; i-- {
		invFact[i-1] = invFact[i] * i % MOD
	}

	R := int64(n)
	ans := fact[N]
	ans = ans * invFact[R] % MOD
	ans = ans * invFact[N-R] % MOD
	fmt.Fprintln(out, ans)
}
