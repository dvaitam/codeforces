package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k, q int
	if _, err := fmt.Fscan(in, &n, &k, &q); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// dp[t][i]: number of ways to be at cell i after t moves starting from any cell
	dp := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int64, n+2) // add padding for boundaries
	}
	for i := 1; i <= n; i++ {
		dp[0][i] = 1
	}
	for t := 1; t <= k; t++ {
		for i := 1; i <= n; i++ {
			var val int64
			if i > 1 {
				val += dp[t-1][i-1]
			}
			if i < n {
				val += dp[t-1][i+1]
			}
			if val >= MOD {
				val %= MOD
			}
			dp[t][i] = val
		}
	}

	coef := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var c int64
		for t := 0; t <= k; t++ {
			c = (c + dp[t][i]*dp[k-t][i]) % MOD
		}
		coef[i] = c
	}

	var total int64
	for i := 1; i <= n; i++ {
		total = (total + a[i]*coef[i]) % MOD
	}

	for ; q > 0; q-- {
		var idx int
		var x int64
		fmt.Fscan(in, &idx, &x)
		total = (total - a[idx]*coef[idx]) % MOD
		if total < 0 {
			total += MOD
		}
		a[idx] = x
		total = (total + a[idx]*coef[idx]) % MOD
		fmt.Fprintln(out, total)
	}
}
