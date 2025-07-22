package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, l, r int
	if _, err := fmt.Fscan(reader, &n, &l, &r); err != nil {
		return
	}

	heights := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &heights[i])
		sum += heights[i]
	}

	important := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &important[i])
	}

	const negInf = -1 << 30
	dp := make([]int, sum+1)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0

	for i := 0; i < n; i++ {
		h := heights[i]
		imp := important[i] == 1
		for s := sum - h; s >= 0; s-- {
			if dp[s] == negInf {
				continue
			}
			gain := 0
			if imp && s >= l && s <= r {
				gain = 1
			}
			if v := dp[s] + gain; v > dp[s+h] {
				dp[s+h] = v
			}
		}
	}

	fmt.Fprintln(writer, dp[sum])
}
