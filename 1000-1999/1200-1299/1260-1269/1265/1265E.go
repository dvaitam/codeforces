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

	inv100 := modPow(100, MOD-2)
	prefix := int64(1)
	sumPrefix := int64(0)
	for i := 0; i < n; i++ {
		var p int64
		fmt.Fscan(in, &p)
		sumPrefix = (sumPrefix + prefix) % MOD
		q := p * inv100 % MOD
		prefix = prefix * q % MOD
	}

	ans := sumPrefix * modPow(prefix, MOD-2) % MOD
	fmt.Fprintln(out, ans)
}
