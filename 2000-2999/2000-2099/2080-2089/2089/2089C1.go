package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

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

func main() {
	const maxL = 5000
	inv := make([]int64, maxL+1)
	for i := 1; i <= maxL; i++ {
		inv[i] = modPow(int64(i), MOD-2)
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, l, k int
		fmt.Fscan(in, &n, &l, &k)

		dpNext := make([]int64, n)
		dpNext[0] = 1
		ans := make([]int64, n)

		for r := l; r >= 1; r-- {
			curr := make([]int64, n)
			q := r / n
			rem := r % n
			qMod := int64(q) % MOD

			pref := make([]int64, 2*n+1)
			for i := 0; i < 2*n; i++ {
				pref[i+1] = (pref[i] + dpNext[i%n]) % MOD
			}

			for j := 0; j < n; j++ {
				base := n + j
				partial := pref[base] - pref[base-rem]
				partial %= MOD
				if partial < 0 {
					partial += MOD
				}
				sumVal := (qMod + partial) % MOD
				curr[j] = sumVal * inv[r] % MOD
			}

			for i := 1; i <= n; i++ {
				idx := i % n
				ans[i-1] = (ans[i-1] + curr[idx]) % MOD
			}

			dpNext = curr
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
