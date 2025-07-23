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
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	const inf = int(1e9)
	dp := make([][3]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 0; j < 3; j++ {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0

	for i := 1; i <= n; i++ {
		// rest
		dp[i][0] = min(min(dp[i-1][0], dp[i-1][1]), dp[i-1][2]) + 1
		if a[i-1] == 1 || a[i-1] == 3 {
			dp[i][1] = min(dp[i-1][0], dp[i-1][2])
		}
		if a[i-1] == 2 || a[i-1] == 3 {
			dp[i][2] = min(dp[i-1][0], dp[i-1][1])
		}
	}

	ans := min(min(dp[n][0], dp[n][1]), dp[n][2])
	fmt.Fprintln(writer, ans)
}
