package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	// sort skills
	sort.Ints(arr[1:])
	// best[i]: smallest index j such that arr[i] - arr[j] <= 5
	best := make([]int, n+1)
	for i := 1; i <= n; i++ {
		for j := i; j >= 1; j-- {
			if arr[i]-arr[j] > 5 {
				break
			}
			best[i] = j
		}
	}
	// dp[i][j]: max students using j teams among first i students
	dp := make([][]int16, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int16, k+1)
	}
	for i := 1; i <= n; i++ {
		// baseline: not take student i
		copy(dp[i], dp[i-1])
		l := best[i]
		sz := int16(i - l + 1)
		// one team
		if dp[i][1] < sz {
			dp[i][1] = sz
		}
		// multiple teams
		for j := 2; j <= k; j++ {
			if dp[l-1][j-1] != 0 {
				cand := dp[l-1][j-1] + sz
				if dp[i][j] < cand {
					dp[i][j] = cand
				}
			}
		}
	}
	// compute result
	var sol int16
	for j := 1; j <= k; j++ {
		if dp[n][j] >= int16(k) && dp[n][j] > sol {
			sol = dp[n][j]
		}
	}
	fmt.Fprint(writer, sol)
}
