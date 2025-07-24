package main

import (
	"bufio"
	"fmt"
	"os"
)

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	slices := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &slices[i])
		total += slices[i]
	}

	dp := make([][2]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		val := slices[i]
		dp[i][0] = max64(val+dp[i+1][1], -val+dp[i+1][0])
		dp[i][1] = min64(-val+dp[i+1][0], val+dp[i+1][1])
	}

	diff := dp[0][1]
	alice := (total + diff) / 2
	bob := total - alice
	fmt.Fprintf(out, "%d %d", alice, bob)
}
