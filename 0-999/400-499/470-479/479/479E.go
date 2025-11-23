package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, a, b, k int
	fmt.Fscan(reader, &n, &a, &b, &k)

	const MOD = 1000000007

	dp := make([]int, n+1)
	dp[a] = 1

	P := make([]int, n+1)
	newDp := make([]int, n+1)

	for step := 0; step < k; step++ {
		currentSum := 0
		for i := 1; i <= n; i++ {
			currentSum += dp[i]
			if currentSum >= MOD {
				currentSum -= MOD
			}
			P[i] = currentSum
		}

		for i := range newDp {
			newDp[i] = 0
		}

		if a < b {
			for y := 1; y < b; y++ {
				limit := (b + y - 1) / 2
				if limit < 1 {
					continue
				}
				val := P[limit]
				if y <= limit {
					val -= dp[y]
					if val < 0 {
						val += MOD
					}
				}
				newDp[y] = val
			}
		} else {
			for y := b + 1; y <= n; y++ {
				start := (b + y) / 2 + 1
				if start > n {
					continue
				}
				val := P[n] - P[start-1]
				if val < 0 {
					val += MOD
				}
				if y >= start {
					val -= dp[y]
					if val < 0 {
						val += MOD
					}
				}
				newDp[y] = val
			}
		}
		dp, newDp = newDp, dp
	}

	ans := 0
	for _, v := range dp {
		ans += v
		if ans >= MOD {
			ans -= MOD
		}
	}
	fmt.Println(ans)
}