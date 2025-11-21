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

	var n int
	fmt.Fscan(in, &n)

	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = 1 << 30
		x := i
		for x > 0 {
			d := x % 10
			x /= 10
			if d > 0 && dp[i-d]+1 < dp[i] {
				dp[i] = dp[i-d] + 1
			}
		}
	}

	fmt.Fprintln(out, dp[n])
}
