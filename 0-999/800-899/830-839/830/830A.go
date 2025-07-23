package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	var p int64
	if _, err := fmt.Fscan(in, &n, &k, &p); err != nil {
		return
	}

	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &b[i])
	}

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

	const inf int64 = 1 << 60
	dp := make([]int64, k+1)
	for j := 0; j <= k; j++ {
		dp[j] = 0
	}

	for i := 1; i <= n; i++ {
		newdp := make([]int64, k+1)
		for j := 0; j <= k; j++ {
			newdp[j] = inf
		}
		for j := 1; j <= k; j++ {
			cost := abs(a[i-1]-b[j-1]) + abs(b[j-1]-p)
			val := max(dp[j-1], cost)
			if val < newdp[j] {
				newdp[j] = val
			}
			if newdp[j-1] < newdp[j] {
				newdp[j] = newdp[j-1]
			}
		}
		dp = newdp
	}

	fmt.Println(dp[k])
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
