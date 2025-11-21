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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		dp := make([]int, n+2)
		for i := range dp {
			dp[i] = -1 << 60
		}
		dp[0] = 0

		for i := 0; i < n; i++ {
			next := make([]int, n+2)
			for j := range next {
				next[j] = -1 << 60
			}
			for j := 0; j <= i; j++ {
				if dp[j] <= -(1 << 55) {
					continue
				}
				// Monocarp rests; Stereocarp doesn't act.
				if next[j] < dp[j] {
					next[j] = dp[j]
				}
				// Monocarp trains.
				gain := a[i]
				if j > 0 {
					gain -= b[i]
				}
				if next[j+1] < dp[j]+gain {
					next[j+1] = dp[j] + gain
				}
			}
			dp = next
		}

		ans := 0
		for j := 0; j <= n; j++ {
			if dp[j] > ans {
				ans = dp[j]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
