package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(b []int) bool {
	n := len(b)
	dp := make([]bool, n+1)
	dp[0] = true
	for i := 0; i < n; i++ {
		if dp[i] && i+b[i]+1 <= n {
			dp[i+b[i]+1] = true
		}
		if i-b[i] >= 0 && dp[i-b[i]] {
			dp[i+1] = true
		}
	}
	return dp[n]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if solve(b) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
