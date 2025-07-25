package main

import (
	"bufio"
	"fmt"
	"os"
)

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		const Inf int64 = 1 << 60
		dp := make([][]int64, n+1)
		for i := range dp {
			dp[i] = make([]int64, k+1)
			for j := range dp[i] {
				dp[i][j] = Inf
			}
		}
		dp[0][0] = 0

		for i := 1; i <= n; i++ {
			val := arr[i-1]
			for j := 0; j <= k; j++ {
				if dp[i-1][j]+val < dp[i][j] {
					dp[i][j] = dp[i-1][j] + val
				}
			}
			minVal := val
			for length := 2; length <= k+1 && length <= i; length++ {
				if arr[i-length] < minVal {
					minVal = arr[i-length]
				}
				cost := length - 1
				for j := cost; j <= k; j++ {
					cand := dp[i-length][j-cost] + int64(length)*minVal
					if cand < dp[i][j] {
						dp[i][j] = cand
					}
				}
			}
		}

		ans := dp[n][0]
		for j := 1; j <= k; j++ {
			if dp[n][j] < ans {
				ans = dp[n][j]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
