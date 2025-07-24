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

func expectedMoves(n, d int) int64 {
	if d == 0 {
		return 0
	}
	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[MOD%int64(i)]%MOD
	}
	A := make([]int64, n+1)
	B := make([]int64, n+1)
	if n >= 1 {
		A[1] = 1
	}
	for i := 1; i < n; i++ {
		denom := n - i
		x1 := (int64(n)*A[i] - int64(i)*A[i-1]) % MOD
		if x1 < 0 {
			x1 += MOD
		}
		A[i+1] = x1 * inv[denom] % MOD

		x2 := (int64(n)*B[i] - int64(i)*B[i-1] - int64(n)) % MOD
		if x2 < 0 {
			x2 += MOD
		}
		B[i+1] = x2 * inv[denom] % MOD
	}
	diff := (A[n] - A[n-1]) % MOD
	if diff < 0 {
		diff += MOD
	}
	x := (B[n-1] + 1 - B[n]) % MOD
	if x < 0 {
		x += MOD
	}
	x = x * modPow(diff, MOD-2) % MOD
	res := (A[d]*x + B[d]) % MOD
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		d := 0
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				d++
			}
		}
		ans := expectedMoves(n, d)
		fmt.Fprintln(out, ans)
	}
}
