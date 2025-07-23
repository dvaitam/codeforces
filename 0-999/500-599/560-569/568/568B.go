package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

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
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	// precompute factorials and inverse factorials
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	comb := func(nn, kk int) int64 {
		if kk < 0 || kk > nn {
			return 0
		}
		return fact[nn] * invFact[kk] % MOD * invFact[nn-kk] % MOD
	}

	// compute Bell numbers up to n using recurrence
	bell := make([]int64, n+1)
	bell[0] = 1
	for i := 0; i < n; i++ {
		sum := int64(0)
		for k := 0; k <= i; k++ {
			sum = (sum + comb(i, k)*bell[k]) % MOD
		}
		bell[i+1] = sum
	}

	ans := int64(0)
	for m := 0; m < n; m++ {
		ans = (ans + comb(n, m)*bell[m]) % MOD
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans%MOD)
	out.Flush()
}
