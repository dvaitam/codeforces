package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, b int64) int64 {
	a %= MOD
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		if n%2 == 1 {
			base := (modPow(2, n-1) + 1) % MOD
			ans := modPow(base, k)
			fmt.Fprintln(writer, ans)
		} else {
			a := (modPow(2, n-1) - 1 + MOD) % MOD
			pow2n := modPow(2, n)
			ans := int64(1)
			cur := int64(1)
			for i := int64(1); i <= k; i++ {
				ans = (a*ans + cur) % MOD
				cur = cur * pow2n % MOD
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
