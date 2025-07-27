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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		dp := make([]int, n+1)
		ans := 1
		for i := 1; i <= n; i++ {
			if dp[i] == 0 {
				dp[i] = 1
			}
			for j := i * 2; j <= n; j += i {
				if a[j] > a[i] && dp[i]+1 > dp[j] {
					dp[j] = dp[i] + 1
				}
			}
			if dp[i] > ans {
				ans = dp[i]
			}
		}
		for i := 1; i <= n; i++ {
			if dp[i] > ans {
				ans = dp[i]
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
