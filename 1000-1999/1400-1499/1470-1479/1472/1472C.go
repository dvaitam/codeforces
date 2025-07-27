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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		dp := make([]int64, n)
		var ans int64
		for i := n - 1; i >= 0; i-- {
			dp[i] = a[i]
			next := i + int(a[i])
			if next < n {
				dp[i] += dp[next]
			}
			if dp[i] > ans {
				ans = dp[i]
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
