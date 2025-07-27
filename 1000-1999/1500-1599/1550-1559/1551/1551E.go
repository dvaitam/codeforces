package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	const neg = -1000000000
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = neg
	}
	dp[0] = 0
	for i := 0; i < n; i++ {
		next := make([]int, n+1)
		for j := range next {
			next[j] = neg
		}
		for t := 0; t <= i; t++ {
			if dp[t] == neg {
				continue
			}
			if dp[t] > next[t] {
				next[t] = dp[t]
			}
			val := dp[t]
			if a[i] == t+1 {
				val++
			}
			if val > next[t+1] {
				next[t+1] = val
			}
		}
		dp = next
	}
	ans := -1
	for t := 0; t <= n; t++ {
		if dp[t] >= k {
			cand := n - t
			if ans == -1 || cand < ans {
				ans = cand
			}
		}
	}
	fmt.Fprintln(writer, ans)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		solve(reader, writer)
	}
}
