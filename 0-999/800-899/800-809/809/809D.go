package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	intervals := make([][2]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &intervals[i][0], &intervals[i][1])
	}
	const INF int64 = 1 << 60
	dp := make([]int64, n+2)
	for i := 1; i < len(dp); i++ {
		dp[i] = INF
	}
	dp[0] = -INF
	maxLen := 0
	for _, it := range intervals {
		l, r := it[0], it[1]
		for j := maxLen; j >= 0; j-- {
			if dp[j] == INF {
				continue
			}
			x := dp[j] + 1
			if x < l {
				x = l
			}
			if x <= r && x < dp[j+1] {
				dp[j+1] = x
				if j+1 > maxLen {
					maxLen = j + 1
				}
			}
		}
	}
	fmt.Println(maxLen)
}
