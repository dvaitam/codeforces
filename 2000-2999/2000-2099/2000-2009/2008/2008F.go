package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(base, exp int64) int64 {
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		sum := int64(0)
		sumSq := int64(0)
		for i := int64(0); i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			x %= MOD
			sum = (sum + x) % MOD
			sumSq = (sumSq + x*x%MOD) % MOD
		}
		numerator := (sum*sum%MOD - sumSq + MOD) % MOD
		denominator := (n % MOD) * ((n - 1) % MOD) % MOD
		ans := numerator * modInv(denominator) % MOD
		fmt.Fprintln(out, ans)
	}
}
