package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	var t string
	fmt.Fscan(in, &t)
	m := len(t)
	size := m + n + 1
	dp := make([]int64, size)
	dp[0] = 1
	last := make([]int, k)
	for i := 1; i <= m; i++ {
		c := int(t[i-1] - 'a')
		dp[i] = (2*dp[i-1] - dp[last[c]] + mod) % mod
		last[c] = i
	}
	for i := m + 1; i <= m+n; i++ {
		best := 0
		minPos := last[0]
		for j := 1; j < k; j++ {
			if last[j] < minPos {
				best = j
				minPos = last[j]
			}
		}
		dp[i] = (2*dp[i-1] - dp[last[best]] + mod) % mod
		last[best] = i
	}
	fmt.Println(dp[m+n] % mod)
}
