package main

import (
	"bufio"
	"fmt"
	"os"
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

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	h := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}

	var ans int64
	for i := 0; i < n; i++ {
		if h[i] > 1 {
			ans = (ans + h[i] - 1) % MOD
		}
	}

	pair := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		m := min(h[i], h[i+1])
		if m > 1 {
			pair[i] = m - 1
			ans = (ans + pair[i]*pair[i]) % MOD
		}
	}

	if n <= 2 {
		fmt.Fprintln(out, ans%MOD)
		return
	}

	B := make([]int64, n)
	for i := 1; i < n-1; i++ {
		m := min(h[i-1], h[i])
		m = min(m, h[i+1])
		if m > 1 {
			B[i] = m - 1
		}
	}

	i := 1
	for i < n-1 {
		if B[i] == 0 {
			i++
			continue
		}
		L := i
		for i <= n-2 && B[i] > 0 {
			i++
		}
		R := i - 1
		size := R - L + 1
		prefProd := make([]int64, size+1)
		prefProd[0] = 1
		for j := 0; j < size; j++ {
			prefProd[j+1] = prefProd[j] * B[L+j] % MOD
		}
		arrVal := make([]int64, size)
		prefixSum := make([]int64, size)
		for j := 0; j < size; j++ {
			arrVal[j] = pair[L-1+j] * modInv(prefProd[j]) % MOD
			if j == 0 {
				prefixSum[j] = arrVal[j] % MOD
			} else {
				prefixSum[j] = (prefixSum[j-1] + arrVal[j]) % MOD
			}
		}
		for r := L + 1; r <= R+1; r++ {
			idx := r - 1 - L
			prod := prefProd[idx+1]
			sumVal := prefixSum[idx]
			contrib := pair[r-1] % MOD
			contrib = contrib * prod % MOD
			contrib = contrib * sumVal % MOD
			ans = (ans + contrib) % MOD
		}
	}

	fmt.Fprintln(out, ans%MOD)
}
