package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	inv2 := (MOD + 1) / 2
	for ; t > 0; t-- {
		var n, m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		sumF := int64(0)
		for i := int64(0); i < m; i++ {
			var a, b, f int64
			fmt.Fscan(reader, &a, &b, &f)
			sumF += f
		}
		totalPairs := n * (n - 1) / 2
		modP := totalPairs % MOD
		invP := modPow(modP, MOD-2)

		term1 := k % MOD
		term1 = term1 * invP % MOD
		term1 = term1 * (sumF % MOD) % MOD

		term2 := m % MOD
		term2 = term2 * (k % MOD) % MOD
		term2 = term2 * ((k - 1) % MOD) % MOD
		term2 = term2 * inv2 % MOD
		term2 = term2 * invP % MOD
		term2 = term2 * invP % MOD

		ans := (term1 + term2) % MOD
		if ans < 0 {
			ans += MOD
		}
		fmt.Fprintln(writer, ans)
	}
}
