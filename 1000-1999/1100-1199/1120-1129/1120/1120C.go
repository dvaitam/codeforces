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
	var a, b int
	fmt.Fscan(in, &n, &a, &b)
	var s string
	fmt.Fscan(in, &s)

	dp := make([]int, n+1)
	const inf = int(1e9)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	nxt := make([][]int, n)
	for i := range nxt {
		nxt[i] = make([]int, 26)
		for j := 0; j < 26; j++ {
			nxt[i][j] = n
		}
	}
	last := make([]int, 26)
	for c := 0; c < 26; c++ {
		last[c] = n
	}
	for i := n - 1; i >= 0; i-- {
		for c := 0; c < 26; c++ {
			nxt[i][c] = last[c]
		}
		last[s[i]-'a'] = i
	}

	for i := 0; i < n; i++ {
		if dp[i]+a < dp[i+1] {
			dp[i+1] = dp[i] + a
		}
		j := i + 1
		best := n
		for ; j <= n; j++ {
			if best == n {
				if dp[i]+b < dp[j] {
					dp[j] = dp[i] + b
				}
				break
			}
			best = nxt[best][s[j-1]-'a']
			if best == n {
				break
			}
			if dp[i]+b < dp[j] {
				dp[j] = dp[i] + b
			}
		}
	}

	fmt.Fprintln(out, dp[n])
}
