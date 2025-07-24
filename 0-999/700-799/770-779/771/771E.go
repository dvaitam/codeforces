package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}

	// dynamic programming over prefix sums
	prefix1, prefix2, prefixAll := int64(0), int64(0), int64(0)
	last1 := map[int64]int{0: 0}
	last2 := map[int64]int{0: 0}
	lastAll := map[int64]int{0: 0}
	lastPair := map[[2]int64]int{{0, 0}: 0}

	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix1 += a[i-1]
		prefix2 += b[i-1]
		prefixAll += a[i-1] + b[i-1]
		best := dp[i-1]
		if j, ok := last1[prefix1]; ok {
			if best < dp[j]+1 {
				best = dp[j] + 1
			}
		}
		if j, ok := last2[prefix2]; ok {
			if best < dp[j]+1 {
				best = dp[j] + 1
			}
		}
		if j, ok := lastAll[prefixAll]; ok {
			if best < dp[j]+1 {
				best = dp[j] + 1
			}
		}
		if j, ok := lastPair[[2]int64{prefix1, prefix2}]; ok {
			if best < dp[j]+2 {
				best = dp[j] + 2
			}
		}
		dp[i] = best
		last1[prefix1] = i
		last2[prefix2] = i
		lastAll[prefixAll] = i
		lastPair[[2]int64{prefix1, prefix2}] = i
	}

	fmt.Fprintln(out, dp[n])
}
