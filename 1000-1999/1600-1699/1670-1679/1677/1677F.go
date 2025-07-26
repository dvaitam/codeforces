package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func powMod(a, b int64) int64 {
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

func calcB(a, p, k int64) int64 {
	var res int64
	for t := int64(1); t <= a; t++ {
		term := powMod(p, t) * powMod(t, k) % MOD
		res = (res + term) % MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k, p int64
	fmt.Fscan(in, &n, &k, &p)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	A := make([]int64, n+1)
	B := make([]int64, n+1)
	invA := make([]int64, n+1)
	prefProd := make([]int64, n+1)
	prefX := make([]int64, n+1)
	prefY := make([]int64, n+1)
	prefZ := make([]int64, n+1)

	prefProd[0] = 1
	for i := 1; i <= n; i++ {
		A[i] = (a[i] + 1) % MOD
		B[i] = calcB(a[i], p, k) % MOD
		invA[i] = powMod(A[i], MOD-2)

		prefProd[i] = prefProd[i-1] * A[i] % MOD
		prefX[i] = (prefX[i-1] + B[i]*invA[i]) % MOD
		prefY[i] = (prefY[i-1] + a[i]%MOD*invA[i]) % MOD
		prefZ[i] = (prefZ[i-1] + B[i]*invA[i]%MOD*a[i]%MOD*invA[i]) % MOD
	}

	var ans int64
	for l := 1; l <= n; l++ {
		for r := l; r <= n; r++ {
			prod := prefProd[r] * powMod(prefProd[l-1], MOD-2) % MOD
			X := (prefX[r] - prefX[l-1]) % MOD
			if X < 0 {
				X += MOD
			}
			Y := (prefY[r] - prefY[l-1]) % MOD
			if Y < 0 {
				Y += MOD
			}
			Z := (prefZ[r] - prefZ[l-1]) % MOD
			if Z < 0 {
				Z += MOD
			}
			val := (X + X*Y%MOD - Z) % MOD
			if val < 0 {
				val += MOD
			}
			ans = (ans + prod*val) % MOD
		}
	}

	fmt.Fprintln(out, ans)
}
