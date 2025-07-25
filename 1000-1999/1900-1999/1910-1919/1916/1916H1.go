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

func modInv(x int64) int64 {
	return modPow(x, MOD-2)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	var p int64
	var k int
	if _, err := fmt.Fscan(in, &n, &p, &k); err != nil {
		return
	}

	if k > int(n) {
		// no ranks higher than n
		// but we still need powers up to k for loops
	}

	powP := make([]int64, k+1)
	powP[0] = 1
	for i := 1; i <= k; i++ {
		powP[i] = powP[i-1] * p % MOD
	}
	pn := modPow(p, n)

	results := make([]int64, k+1)
	for r := 0; r <= k; r++ {
		if int64(r) > n {
			results[r] = 0
			continue
		}
		if r == 0 {
			results[r] = 1
			continue
		}
		prod1 := int64(1)
		for i := 0; i < r; i++ {
			prod1 = prod1 * ((pn - powP[i] + MOD) % MOD) % MOD
		}
		pr := powP[r]
		denom := int64(1)
		for i := 0; i < r; i++ {
			denom = denom * ((pr - powP[i] + MOD) % MOD) % MOD
		}
		results[r] = prod1 * prod1 % MOD * modInv(denom) % MOD
	}

	for i := 0; i <= k; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, results[i])
	}
}
