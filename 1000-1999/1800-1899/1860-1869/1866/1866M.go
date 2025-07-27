package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pArr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pArr[i])
	}

	inv100 := modPow(100, MOD-2)
	pval := make([]int64, 100)
	for i := 0; i < 100; i++ {
		pval[i] = int64(i) * inv100 % MOD
	}

	H := make([]int64, 100)
	s := int64(0)

	for i := 0; i < n; i++ {
		p := pval[pArr[i]]
		inv1 := modPow((1-p+MOD)%MOD, MOD-2)
		f := p * H[pArr[i]] % MOD
		t := (inv1 + p*inv1%MOD*s%MOD - f) % MOD
		if t < 0 {
			t += MOD
		}
		s = (s + t) % MOD
		for r := 0; r < 100; r++ {
			H[r] = (H[r]*pval[r] + s) % MOD
		}
	}

	fmt.Fprintln(out, s)
}
