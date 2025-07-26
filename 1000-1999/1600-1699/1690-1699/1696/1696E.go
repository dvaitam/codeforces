package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	maxA := a[0]
	dp := make([]int64, maxA+2)
	next := make([]int64, maxA+2)
	dp[0] = 1
	var ans int64

	for _, w := range a {
		for j := 0; j < w; j++ {
			val := dp[j] % MOD
			if val == 0 {
				continue
			}
			ans = (ans + val) % MOD
			next[j] = (next[j] + val) % MOD
			dp[j+1] = (dp[j+1] + val) % MOD
			dp[j] = 0
		}
		for j := w; j <= maxA; j++ {
			if dp[j] != 0 {
				next[j] = (next[j] + dp[j]) % MOD
				dp[j] = 0
			}
		}
		dp, next = next, make([]int64, maxA+2)
	}

	fmt.Fprintln(writer, ans%MOD)
}
