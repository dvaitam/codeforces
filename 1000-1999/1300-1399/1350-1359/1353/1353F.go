package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([][]int64, n)
		for i := 0; i < n; i++ {
			a[i] = make([]int64, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
		fmt.Fprintln(out, solve(n, m, a))
	}
}

func solve(n, m int, a [][]int64) int64 {
	best := int64(1 << 62)
	// Enumerate candidate base values from each cell
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			base := a[i][j] - int64(i+j)
			cost := computeCost(base, n, m, a)
			if cost < best {
				best = cost
			}
		}
	}
	return best
}

func computeCost(base int64, n, m int, a [][]int64) int64 {
	// dp[i][j] minimal operations to reach (i,j)
	dp := make([][]int64, n)
	for i := range dp {
		dp[i] = make([]int64, m)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			target := base + int64(i+j)
			if target > a[i][j] {
				dp[i][j] = INF
				continue
			}
			cost := a[i][j] - target
			if i == 0 && j == 0 {
				dp[i][j] = cost
			} else {
				prev := INF
				if i > 0 && dp[i-1][j] < prev {
					prev = dp[i-1][j]
				}
				if j > 0 && dp[i][j-1] < prev {
					prev = dp[i][j-1]
				}
				if prev != INF {
					dp[i][j] = prev + cost
				}
			}
		}
	}
	return dp[n-1][m-1]
}
