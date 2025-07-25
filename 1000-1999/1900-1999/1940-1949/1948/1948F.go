package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(a, b int) int {
	res := 1
	base := a % mod
	exp := b
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	prefA := make([]int, n+1)
	prefB := make([]int, n+1)
	totalA := 0
	totalB := 0
	for i := 1; i <= n; i++ {
		totalA += a[i]
		totalB += b[i]
		prefA[i] = prefA[i-1] + a[i]
		prefB[i] = prefB[i-1] + b[i]
	}

	// precompute combinations prefix for totalB
	fac := make([]int, totalB+1)
	invfac := make([]int, totalB+1)
	fac[0] = 1
	for i := 1; i <= totalB; i++ {
		fac[i] = fac[i-1] * i % mod
	}
	invfac[totalB] = modPow(fac[totalB], mod-2)
	for i := totalB; i >= 1; i-- {
		invfac[i-1] = invfac[i] * i % mod
	}

	combPrefix := make([]int, totalB+1)
	prefix := 0
	for i := 0; i <= totalB; i++ {
		c := fac[totalB] * invfac[i] % mod * invfac[totalB-i] % mod
		prefix += c
		if prefix >= mod {
			prefix -= mod
		}
		combPrefix[i] = prefix
	}

	pow2B := modPow(2, totalB)
	invPow2B := modPow(pow2B, mod-2)

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		ARange := prefA[r] - prefA[l-1]
		BRange := prefB[r] - prefB[l-1]
		nOut := totalB - BRange
		diff := 2*ARange - totalA
		threshold := nOut - diff
		if threshold < 0 {
			fmt.Fprintln(writer, 1)
		} else if threshold >= totalB {
			fmt.Fprintln(writer, 0)
		} else {
			val := pow2B - combPrefix[threshold]
			if val < 0 {
				val += mod
			}
			ans := val * invPow2B % mod
			fmt.Fprintln(writer, ans)
		}
	}
}
