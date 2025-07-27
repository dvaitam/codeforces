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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	speeds := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &speeds[i])
	}

	sort.Slice(speeds, func(i, j int) bool { return speeds[i] < speeds[j] })

	dp := make([][]int64, n)
	for i := range dp {
		dp[i] = make([]int64, n)
	}

	for length := 2; length <= n; length++ {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			left := dp[l+1][r]
			right := dp[l][r-1]
			if left < right {
				dp[l][r] = left + speeds[r] - speeds[l]
			} else {
				dp[l][r] = right + speeds[r] - speeds[l]
			}
		}
	}

	fmt.Fprintln(writer, dp[0][n-1])
}
