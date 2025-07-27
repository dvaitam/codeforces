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

	var n, h, l, r int
	if _, err := fmt.Fscan(in, &n, &h, &l, &r); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	const negInf = -1 << 60
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, h)
		for j := range dp[i] {
			dp[i][j] = negInf
		}
	}
	dp[0][0] = 0

	for i := 1; i <= n; i++ {
		ai := a[i-1]
		for t := 0; t < h; t++ {
			if dp[i-1][t] == negInf {
				continue
			}
			// option 1: sleep after ai hours
			nt := (t + ai) % h
			val := dp[i-1][t]
			if l <= nt && nt <= r {
				val++
			}
			if val > dp[i][nt] {
				dp[i][nt] = val
			}
			// option 2: sleep after ai-1 hours
			nt = (t + ai - 1) % h
			val = dp[i-1][t]
			if l <= nt && nt <= r {
				val++
			}
			if val > dp[i][nt] {
				dp[i][nt] = val
			}
		}
	}

	ans := 0
	for t := 0; t < h; t++ {
		if dp[n][t] > ans {
			ans = dp[n][t]
		}
	}
	fmt.Fprintln(out, ans)
}
