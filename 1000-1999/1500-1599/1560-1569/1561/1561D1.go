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

	var n, mod int
	if _, err := fmt.Fscan(reader, &n, &mod); err != nil {
		return
	}

	dp := make([]int, n+2)
	pref := make([]int, n+2)
	diff := make([]int, n+2)

	dp[1] = 1
	pref[1] = 1

	// propagate transitions from cell 1 using division
	for j := 2; j <= n; j++ {
		start := j
		end := j + j
		if end > n+1 {
			end = n + 1
		}
		diff[start] = (diff[start] + dp[1]) % mod
		if end <= n {
			diff[end] = (diff[end] - dp[1]) % mod
		}
	}

	add := 0
	for i := 2; i <= n; i++ {
		add = (add + diff[i]) % mod
		if add < 0 {
			add += mod
		}
		dp[i] = (pref[i-1] + add) % mod
		pref[i] = (pref[i-1] + dp[i]) % mod

		for j := 2; i*j <= n; j++ {
			start := i * j
			end := start + j
			if end > n+1 {
				end = n + 1
			}
			diff[start] = (diff[start] + dp[i]) % mod
			if end <= n {
				diff[end] = (diff[end] - dp[i]) % mod
			}
		}
	}

	ans := dp[n] % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(writer, ans)
}
