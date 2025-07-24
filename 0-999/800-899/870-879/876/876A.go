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
	reader := bufio.NewReader(os.Stdin)
	var n, a, b, c int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fmt.Fscan(reader, &a)
	fmt.Fscan(reader, &b)
	fmt.Fscan(reader, &c)

	const INF int = int(1e9)
	dp := make([][3]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i][0], dp[i][1], dp[i][2] = INF, INF, INF
	}
	dp[1][0] = 0
	for i := 2; i <= n; i++ {
		dp[i][0] = min(dp[i-1][1]+a, dp[i-1][2]+b)
		dp[i][1] = min(dp[i-1][0]+a, dp[i-1][2]+c)
		dp[i][2] = min(dp[i-1][0]+b, dp[i-1][1]+c)
	}
	ans := min(dp[n][0], min(dp[n][1], dp[n][2]))
	fmt.Println(ans)
}
