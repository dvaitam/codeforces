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

	var N, K int
	fmt.Fscan(in, &N, &K)

	limit := N
	if K > limit {
		limit = K
	}
	fact := make([]int64, limit+2)
	invFact := make([]int64, limit+2)
	fact[0] = 1
	for i := 1; i <= limit+1; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[limit+1] = modPow(fact[limit+1], MOD-2)
	for i := limit + 1; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	minNK := N
	if K < minNK {
		minNK = K
	}
	ans := int64(0)
	for m := 0; m <= minNK; m++ {
		var F int64
		if m == 0 {
			F = 1
		} else {
			if m > K {
				F = 0
			} else {
				pow1 := modPow(int64(m+1), int64(K-m+1))
				pow2 := modPow(int64(m), int64(K-m+1))
				diff := (pow1 - pow2) % MOD
				if diff < 0 {
					diff += MOD
				}
				F = fact[m] * diff % MOD
			}
		}
		term := fact[N] * invFact[N-m] % MOD
		ans = (ans + term*F) % MOD
	}
	fmt.Fprintln(out, ans)
}
