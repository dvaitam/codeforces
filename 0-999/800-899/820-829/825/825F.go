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

	var s string
	fmt.Fscan(reader, &s)
	n := len(s)

	digits := make([]int, n+1)
	for i := 1; i <= n; i++ {
		sz := 1
		if i >= 10 {
			sz = 2
		}
		if i >= 100 {
			sz = 3
		}
		if i >= 1000 {
			sz = 4
		}
		digits[i] = sz
	}

	const inf int = 1 << 30
	dp := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = inf
	}

	pi := make([]int, n)
	for start := 0; start < n; start++ {
		m := n - start
		pi = pi[:m]
		pi[0] = 0
		for i := 1; i < m; i++ {
			k := pi[i-1]
			for k > 0 && s[start+i] != s[start+k] {
				k = pi[k-1]
			}
			if s[start+i] == s[start+k] {
				k++
			}
			pi[i] = k
		}
		for l := 1; l <= m; l++ {
			period := l - pi[l-1]
			rep := 1
			if l%period == 0 {
				rep = l / period
			}
			cost := dp[start] + digits[rep] + period
			if cost < dp[start+l] {
				dp[start+l] = cost
			}
		}
	}

	fmt.Fprintln(writer, dp[n])
}
