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

	coins := []int{1, 3, 6, 10, 15}
	const maxBase = 29
	const inf = int(1e9)
	dp := make([]int, maxBase+1)
	for i := 1; i <= maxBase; i++ {
		dp[i] = inf
		for _, c := range coins {
			if i >= c && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}

	var t, n int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		fmt.Fscan(in, &n)
		if n < 15 {
			fmt.Fprintln(out, dp[n])
			continue
		}
		r := (n-15)%15 + 15
		ans := (n - r) / 15
		ans += dp[r]
		fmt.Fprintln(out, ans)
	}
}
