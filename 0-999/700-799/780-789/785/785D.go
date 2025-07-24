package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const MOD int64 = 1000000007

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

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	n := len(s)

	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1]
		if s[i-1] == '(' {
			prefix[i]++
		}
	}
	suffix := make([]int, n+2)
	for i := n; i >= 1; i-- {
		suffix[i] = suffix[i+1]
		if s[i-1] == ')' {
			suffix[i]++
		}
	}

	// precompute factorials
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[n] = modInv(fact[n])
	for i := n - 1; i >= 0; i-- {
		invFact[i] = invFact[i+1] * int64(i+1) % MOD
	}

	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	var ans int64
	for i := 1; i <= n; i++ {
		if s[i-1] == '(' {
			A := prefix[i-1]
			B := suffix[i+1]
			ans = (ans + comb(A+B, A+1)) % MOD
		}
	}

	fmt.Println(ans)
}
