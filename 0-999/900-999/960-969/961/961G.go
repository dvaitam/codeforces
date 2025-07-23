package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func modInverse(a int64) int64 {
	return modPow(a, MOD-2)
}

func stirling(n, k int64, fact, invFact []int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	var sum int64
	for j := int64(0); j <= k; j++ {
		comb := fact[k] * invFact[j] % MOD * invFact[k-j] % MOD
		pow := modPow(j, n)
		if (k-j)%2 == 1 {
			sum -= comb * pow % MOD
		} else {
			sum += comb * pow % MOD
		}
	}
	sum %= MOD
	if sum < 0 {
		sum += MOD
	}
	sum = sum * invFact[k] % MOD
	return sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	weights := make([]int64, n)
	var sumW int64
	for i := int64(0); i < n; i++ {
		fmt.Fscan(in, &weights[i])
		sumW = (sumW + weights[i]) % MOD
	}

	maxK := k
	fact := make([]int64, maxK+1)
	invFact := make([]int64, maxK+1)
	fact[0] = 1
	for i := int64(1); i <= maxK; i++ {
		fact[i] = fact[i-1] * i % MOD
	}
	invFact[maxK] = modInverse(fact[maxK])
	for i := maxK - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * (i + 1) % MOD
	}

	s1 := stirling(n, k, fact, invFact)
	s2 := stirling(n-1, k, fact, invFact)

	P := (s1 + (n-1)%MOD*s2%MOD) % MOD
	ans := sumW * P % MOD
	fmt.Println(ans)
}
