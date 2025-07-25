package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m, n int
	if _, err := fmt.Fscan(reader, &m, &n); err != nil {
		return
	}

	s := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &s[i])
	}
	l := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &l[i])
	}

	ways := make([][]int64, m)
	for i := 0; i < m; i++ {
		ways[i] = make([]int64, m)
		si := int64(s[i])
		li := int64(l[i])
		for j := 0; j < m; j++ {
			ways[i][j] = si*(int64(s[j])+int64(l[j])) + li*int64(s[j])
		}
	}

	dp := make([]int64, m)
	dp[0] = 1
	for day := 0; day < n; day++ {
		ndp := make([]int64, m)
		for i := 0; i < m; i++ {
			if dp[i] == 0 {
				continue
			}
			for j := 0; j < m; j++ {
				ndp[j] = (ndp[j] + dp[i]*ways[i][j]) % MOD
			}
		}
		dp = ndp
	}
	ans := int64(0)
	for i := 0; i < m; i++ {
		ans = (ans + dp[i]) % MOD
	}
	fmt.Fprintln(writer, ans)
}
