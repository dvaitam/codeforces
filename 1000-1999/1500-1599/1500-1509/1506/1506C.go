package main

import (
	"bufio"
	"fmt"
	"os"
)

func longestCommonSubstr(a, b string) int {
	n := len(a)
	m := len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	maxLen := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			}
		}
	}
	return maxLen
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b string
		fmt.Fscan(in, &a)
		fmt.Fscan(in, &b)
		l := longestCommonSubstr(a, b)
		res := len(a) + len(b) - 2*l
		fmt.Fprintln(out, res)
	}
}
