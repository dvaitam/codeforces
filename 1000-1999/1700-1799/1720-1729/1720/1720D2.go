package main

import (
	"bufio"
	"fmt"
	"os"
)

const window = 512

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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		dp := make([]int, n)
		ans := 0
		for j := 0; j < n; j++ {
			dp[j] = 1
			start := j - window
			if start < 0 {
				start = 0
			}
			for i := start; i < j; i++ {
				if a[i]^j < a[j]^i {
					if dp[i]+1 > dp[j] {
						dp[j] = dp[i] + 1
					}
				}
			}
			if dp[j] > ans {
				ans = dp[j]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
