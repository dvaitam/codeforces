package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 999999893

func modPow(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return res
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
		if n == 1 {
			fmt.Fprintln(out, 0)
			continue
		}
		exp := 2*n - 3
		denom := (modPow(2, exp) - 1 + MOD) % MOD
		ans := modPow(denom, MOD-2)
		fmt.Fprintln(out, ans)
	}
}
