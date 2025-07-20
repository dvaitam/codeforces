package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

func modPow(a, b, m int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % m
		}
		a = a * a % m
		b >>= 1
	}
	return res
}

func prepareFact(n int) ([]int64, []int64) {
	fact := make([]int64, n+1)
	inv := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[n] = modPow(fact[n], mod-2, mod)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
	return fact, inv
}

func C(n, r int64, fact, inv []int64) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fact[n] * inv[r] % mod * inv[n-r] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	maxN := 5000
	fact, inv := prepareFact(maxN)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		res := int64(0)
		half := (n - 1) / 2
		for s := 0; s <= half; s++ {
			for x := s + 1; x <= 2*s+1 && x <= n; x++ {
				ways := C(int64(x-1), int64(s), fact, inv) * C(int64(n-x), int64(n-2*s-1), fact, inv) % mod
				res = (res + int64(x)*ways) % mod
			}
		}
		for s := (n + 1) / 2; s <= n; s++ {
			ways := C(int64(n), int64(s), fact, inv)
			mex := int64(2*s + 1)
			res = (res + mex*ways) % mod
		}
		fmt.Fprintln(out, res%mod)
	}
}
