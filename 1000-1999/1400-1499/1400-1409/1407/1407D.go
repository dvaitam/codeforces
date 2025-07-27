package main

import (
	"bufio"
	"fmt"
	"os"
)

// min returns the smaller of a and b.
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

	var n int
	fmt.Fscan(reader, &n)
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &h[i])
	}
	dp := make([]int, n)
	// stacks of indices for increasing and decreasing heights
	inc := make([]int, 0, n)
	dec := make([]int, 0, n)
	inc = append(inc, 0)
	dec = append(dec, 0)
	for i := 1; i < n; i++ {
		// default: jump from previous
		dp[i] = dp[i-1] + 1
		// decreasing stack: handle jumps where h[i] > previous
		for len(dec) > 0 && h[dec[len(dec)-1]] < h[i] {
			j := dec[len(dec)-1]
			dp[i] = min(dp[i], dp[j]+1)
			dec = dec[:len(dec)-1]
		}
		if len(dec) > 0 {
			j := dec[len(dec)-1]
			dp[i] = min(dp[i], dp[j]+1)
			if h[j] == h[i] {
				dec = dec[:len(dec)-1]
			}
		}
		dec = append(dec, i)
		// increasing stack: handle jumps where h[i] < previous
		for len(inc) > 0 && h[inc[len(inc)-1]] > h[i] {
			j := inc[len(inc)-1]
			dp[i] = min(dp[i], dp[j]+1)
			inc = inc[:len(inc)-1]
		}
		if len(inc) > 0 {
			j := inc[len(inc)-1]
			dp[i] = min(dp[i], dp[j]+1)
			if h[j] == h[i] {
				inc = inc[:len(inc)-1]
			}
		}
		inc = append(inc, i)
	}
	// output result
	fmt.Fprintln(writer, dp[n-1])
}
