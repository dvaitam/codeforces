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

	var n, k, s int
	if _, err := fmt.Fscan(reader, &n, &k, &s); err != nil {
		return
	}
	q := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &q[i])
	}

	// maximum useful swaps is k*(n-k)
	maxS := k * (n - k)
	if s > maxS {
		s = maxS
	}

	const inf int64 = 1 << 60
	dp := make([][]int64, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int64, s+1)
		for j := 0; j <= s; j++ {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0

	for i := 1; i <= n; i++ {
		val := int64(q[i-1])
		upper := min(i, k)
		for j := upper; j >= 1; j-- {
			cost := i - j
			if cost > s {
				continue
			}
			for c := cost; c <= s; c++ {
				prev := dp[j-1][c-cost]
				if prev+val < dp[j][c] {
					dp[j][c] = prev + val
				}
			}
		}
	}

	ans := inf
	for c := 0; c <= s; c++ {
		if dp[k][c] < ans {
			ans = dp[k][c]
		}
	}
	fmt.Fprintln(writer, ans)
}
