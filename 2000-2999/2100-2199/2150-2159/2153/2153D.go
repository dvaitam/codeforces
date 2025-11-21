package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func range3(a, b, c int64) int64 {
	return max64(a, max64(b, c)) - min64(a, min64(b, c))
}

func solvePath(nums []int64) int64 {
	m := len(nums)
	if m == 0 {
		return 0
	}
	if m == 1 {
		return inf
	}
	dp := make([]int64, m+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for i := 2; i <= m; i++ {
		pairCost := abs64(nums[i-2] - nums[i-1])
		if dp[i-2]+pairCost < dp[i] {
			dp[i] = dp[i-2] + pairCost
		}
		if i >= 3 {
			tripleCost := range3(nums[i-3], nums[i-2], nums[i-1])
			if dp[i-3]+tripleCost < dp[i] {
				dp[i] = dp[i-3] + tripleCost
			}
		}
	}
	return dp[m]
}

func solveCase(a []int64) int64 {
	n := len(a)
	best := solvePath(a)
	if n >= 2 {
		rest := solvePath(a[1 : n-1])
		if rest < inf {
			cost := rest + abs64(a[n-1]-a[0])
			if cost < best {
				best = cost
			}
		}
	}
	if n >= 3 {
		rest := solvePath(a[1 : n-2])
		if rest < inf {
			cost := rest + range3(a[n-2], a[n-1], a[0])
			if cost < best {
				best = cost
			}
		}
		rest = solvePath(a[2 : n-1])
		if rest < inf {
			cost := rest + range3(a[n-1], a[0], a[1])
			if cost < best {
				best = cost
			}
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solveCase(a))
	}
}
