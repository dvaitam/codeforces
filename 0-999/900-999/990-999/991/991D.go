package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s1, s2 string
	fmt.Fscan(reader, &s1)
	fmt.Fscan(reader, &s2)
	n := len(s1)
	blocked := make([]int, n)
	for i := 0; i < n; i++ {
		if s1[i] == 'X' {
			blocked[i] |= 1
		}
		if s2[i] == 'X' {
			blocked[i] |= 2
		}
	}
	dp := make([][4]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 0; j < 4; j++ {
			dp[i][j] = -1000000
		}
	}
	dp[0][0] = 0
	for i := 0; i < n; i++ {
		for mask := 0; mask < 4; mask++ {
			if dp[i][mask] < 0 {
				continue
			}
			curBlock := blocked[i]
			// skip placing a piece at column i
			dp[i+1][0] = max(dp[i+1][0], dp[i][mask])
			if i+1 >= n {
				continue
			}
			nxtBlock := blocked[i+1]
			// orientation 1: missing top-left
			if (mask&2) == 0 && (curBlock&2) == 0 && (nxtBlock&1) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][3] = max(dp[i+1][3], dp[i][mask]+1)
			}
			// orientation 2: missing top-right
			if (mask&1) == 0 && (curBlock&1) == 0 && (mask&2) == 0 && (curBlock&2) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][2] = max(dp[i+1][2], dp[i][mask]+1)
			}
			// orientation 3: missing bottom-left
			if (mask&1) == 0 && (curBlock&1) == 0 && (nxtBlock&1) == 0 && (nxtBlock&2) == 0 {
				dp[i+1][3] = max(dp[i+1][3], dp[i][mask]+1)
			}
			// orientation 4: missing bottom-right
			if (mask&1) == 0 && (curBlock&1) == 0 && (nxtBlock&1) == 0 && (mask&2) == 0 && (curBlock&2) == 0 {
				dp[i+1][1] = max(dp[i+1][1], dp[i][mask]+1)
			}
		}
	}
	ans := dp[n][0]
	fmt.Println(ans)
}
