package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s, t string
	if _, err := fmt.Fscanln(reader, &s); err != nil {
		return
	}
	if _, err := fmt.Fscanln(reader, &t); err != nil {
		return
	}
	n, m := len(s), len(t)
	// dp[i][j]: min cost to convert s[:i] to t[:j]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	// initialize
	for i := 1; i <= n; i++ {
		dp[i][0] = dp[i-1][0] + int(s[i-1]-'a'+1)
	}
	for j := 1; j <= m; j++ {
		dp[0][j] = dp[0][j-1] + int(t[j-1]-'a'+1)
	}
	// compute
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			// deletion
			costDel := dp[i-1][j] + int(s[i-1]-'a'+1)
			// insertion
			costIns := dp[i][j-1] + int(t[j-1]-'a'+1)
			// substitution or match
			diff := int(s[i-1] - t[j-1])
			if diff < 0 {
				diff = -diff
			}
			costSub := dp[i-1][j-1] + diff
			// take min
			best := costDel
			if costIns < best {
				best = costIns
			}
			if costSub < best {
				best = costSub
			}
			dp[i][j] = best
		}
	}
	fmt.Println(dp[n][m])
}
