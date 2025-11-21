package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
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
		var a, b, c string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		fmt.Fscan(in, &c)

		n, m := len(a), len(b)
		const inf = int(1e9)
		dp := make([][]int, n+1)
		for i := range dp {
			dp[i] = make([]int, m+1)
			for j := range dp[i] {
				dp[i][j] = inf
			}
		}
		dp[0][0] = 0

		total := len(c)
		for i := 0; i <= n; i++ {
			for j := 0; j <= m; j++ {
				idx := i + j
				if idx >= total {
					continue
				}
				if i < n {
					cost := dp[i][j]
					if a[i] != c[idx] {
						cost++
					}
					if cost < dp[i+1][j] {
						dp[i+1][j] = cost
					}
				}
				if j < m {
					cost := dp[i][j]
					if b[j] != c[idx] {
						cost++
					}
					if cost < dp[i][j+1] {
						dp[i][j+1] = cost
					}
				}
			}
		}
		fmt.Fprintln(out, dp[n][m])
	}
}
