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

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, c int
	if _, err := fmt.Fscan(in, &n, &k, &c); err != nil {
		return
	}
	r := n % k
	var t string
	fmt.Fscan(in, &t)
	q := n / k

	letters := make([]int, r)
	for i := 0; i < r; i++ {
		letters[i] = int(t[i] - 'a')
	}

	powAlpha := make([]int64, q+1)
	powAlpha[0] = 1
	for i := 1; i <= q; i++ {
		powAlpha[i] = powAlpha[i-1] * int64(c) % MOD
	}

	dpPrev := make([]int64, q+1)
	dpPrev[0] = 1

	for _, v := range letters {
		B := c - v - 1
		C := c - v

		powB := make([]int64, q+1)
		powC := make([]int64, q+1)
		powB[0] = 1
		powC[0] = 1
		for i := 1; i <= q; i++ {
			powB[i] = powB[i-1] * int64(B) % MOD
			powC[i] = powC[i-1] * int64(C) % MOD
		}

		dpCur := make([]int64, q+1)
		if B == 0 {
			for tPrev := 0; tPrev <= q; tPrev++ {
				dpCur[tPrev] = (dpCur[tPrev] + dpPrev[tPrev]*powAlpha[tPrev]) % MOD
			}
			dpPrev = dpCur
			continue
		}
		invB := modInv(int64(B))
		invPowB := make([]int64, q+1)
		invPowB[0] = 1
		for i := 1; i <= q; i++ {
			invPowB[i] = invPowB[i-1] * invB % MOD
		}
		prefix := make([]int64, q+1)
		prefix[0] = dpPrev[0] * powAlpha[0] % MOD * invPowB[0] % MOD
		for i := 1; i <= q; i++ {
			prefix[i] = (prefix[i-1] + dpPrev[i]*powAlpha[i]%MOD*invPowB[i]%MOD) % MOD
		}
		for tIdx := 0; tIdx <= q; tIdx++ {
			dpCur[tIdx] = powB[tIdx] * powC[q-tIdx] % MOD * prefix[tIdx] % MOD
		}
		dpPrev = dpCur
	}

	total := int64(0)
	for _, v := range dpPrev {
		total = (total + v) % MOD
	}
	extra := int64(q * (k - r))
	total = total * modPow(int64(c), extra) % MOD
	fmt.Fprintln(out, total)
}
