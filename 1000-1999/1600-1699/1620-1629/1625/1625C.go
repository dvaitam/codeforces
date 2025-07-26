package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, l, k int
	if _, err := fmt.Fscan(reader, &n, &l, &k); err != nil {
		return
	}
	d := make([]int, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &d[i])
	}
	d[n] = l
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	const INF int64 = 1 << 60
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, k+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0

	for i := 0; i < n; i++ {
		for j := 0; j <= k; j++ {
			if dp[i][j] == INF {
				continue
			}
			for t := i + 1; t <= n; t++ {
				removed := t - i - 1
				nj := j + removed
				if nj > k {
					break
				}
				cost := int64(d[t]-d[i]) * int64(a[i])
				if dp[t][nj] > dp[i][j]+cost {
					dp[t][nj] = dp[i][j] + cost
				}
			}
		}
	}

	ans := INF
	for j := 0; j <= k; j++ {
		ans = min(ans, dp[n][j])
	}
	fmt.Fprintln(writer, ans)
}
