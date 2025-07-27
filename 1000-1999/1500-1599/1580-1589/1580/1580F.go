package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)

	// naive DP over all possible values 0..m-1
	ans := int64(0)
	for start := 0; start < m; start++ {
		dp := make([]int64, m)
		dp[start] = 1
		for step := 1; step < n; step++ {
			prefix := make([]int64, m+1)
			for i := 0; i < m; i++ {
				prefix[i+1] = (prefix[i] + dp[i]) % MOD
			}
			ndp := make([]int64, m)
			for y := 0; y < m; y++ {
				ndp[y] = prefix[m-y] % MOD
			}
			dp = ndp
		}
		for last := 0; last < m; last++ {
			if last+start < m {
				ans = (ans + dp[last]) % MOD
			}
		}
	}

	fmt.Println(ans % MOD)
}
