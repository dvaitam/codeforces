package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		sum += arr[i]
	}
	if sum%int64(n) != 0 {
		fmt.Fprintln(out, 0)
		return
	}
	avg := sum / int64(n)

	// frequency map
	freq := make(map[int64]int)
	Lvals := make(map[int64]int)
	Gvals := make(map[int64]int)
	L, G, M := 0, 0, 0
	for _, v := range arr {
		freq[v]++
		if v < avg {
			Lvals[v]++
			L++
		} else if v > avg {
			Gvals[v]++
			G++
		} else {
			M++
		}
	}

	// precompute factorials
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	if L == 0 || G == 0 {
		// all permutations
		ans := fact[n]
		for _, c := range freq {
			ans = ans * invFact[c] % MOD
		}
		fmt.Fprintln(out, ans)
		return
	}

	// number of sign sequences
	signSeq := (2 * comb(n, M)) % MOD

	// permutations within groups
	denom := int64(1)
	for _, c := range Lvals {
		denom = denom * fact[c] % MOD
	}
	permL := fact[L] * modPow(denom, MOD-2) % MOD

	denom = 1
	for _, c := range Gvals {
		denom = denom * fact[c] % MOD
	}
	permG := fact[G] * modPow(denom, MOD-2) % MOD

	ans := signSeq * permL % MOD * permG % MOD
	fmt.Fprintln(out, ans)
}
