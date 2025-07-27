package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var mod int64
	if _, err := fmt.Fscan(reader, &n, &mod); err != nil {
		return
	}

	dp := make([]int64, n+2)
	diff := make([]int64, n+2)

	dp[1] = 1 % mod
	prefix := dp[1]
	// propagate contributions from dp[1]
	for j := 2; j <= n; j++ {
		l := j
		r := l + j
		if r > n+1 {
			r = n + 1
		}
		diff[l] = (diff[l] + dp[1]) % mod
		diff[r] = (diff[r] - dp[1]) % mod
	}

	cur := int64(0)
	for i := 2; i <= n; i++ {
		cur = (cur + diff[i]) % mod
		if cur < 0 {
			cur += mod
		}
		dp[i] = (prefix + cur) % mod
		prefix = (prefix + dp[i]) % mod
		for j := 2; i*j <= n; j++ {
			l := i * j
			r := l + j
			if r > n+1 {
				r = n + 1
			}
			diff[l] = (diff[l] + dp[i]) % mod
			diff[r] = (diff[r] - dp[i]) % mod
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, dp[n]%mod)
	writer.Flush()
}
