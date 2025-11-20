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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var x, y int64
		fmt.Fscan(in, &n, &x, &y)
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		pos := make([]int, 0)
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				pos = append(pos, i)
			}
		}
		k := len(pos)
		if k%2 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		if k == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		if x >= y {
			if k == 2 && pos[1]-pos[0] == 1 {
				fmt.Fprintln(out, min(x, 2*y))
			} else {
				fmt.Fprintln(out, int64(k/2)*y)
			}
			continue
		}
		dp := make([][2]int64, k)
		const inf int64 = 1e18
		for i := range dp {
			dp[i][0], dp[i][1] = inf, inf
		}
		dp[0][1] = 0
		for i := 1; i < k; i++ {
			dist := int64(pos[i] - pos[i-1])
			costSwap := dist * x
			dp[i][0] = min(dp[i-1][1]+y, costSwap+func() int64 {
				if i >= 2 {
					return dp[i-2][0]
				}
				return 0
			}())
			if i >= 2 {
				dp[i][0] = min(dp[i][0], dp[i-2][0]+costSwap)
			}
			dp[i][1] = min(dp[i-1][0], func() int64 {
				if i >= 2 {
					return dp[i-2][1] + costSwap
				}
				return inf
			}())
		}
		fmt.Fprintln(out, dp[k-1][0])
	}
}
