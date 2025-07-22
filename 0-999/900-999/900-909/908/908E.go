package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	var m, n int
	if _, err := fmt.Fscan(in, &m, &n); err != nil {
		return
	}
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &lines[i])
	}

	patternCount := make(map[uint64]int)
	for j := 0; j < m; j++ {
		var pat uint64 = 0
		for i := 0; i < n; i++ {
			if lines[i][j] == '1' {
				pat |= 1 << uint(i)
			}
		}
		patternCount[pat]++
	}

	maxSize := 0
	for _, c := range patternCount {
		if c > maxSize {
			maxSize = c
		}
	}

	bell := bellNumbers(maxSize)

	ans := 1
	for _, c := range patternCount {
		ans = int(int64(ans) * int64(bell[c]) % int64(MOD))
	}
	fmt.Println(ans)
}

func bellNumbers(n int) []int {
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, n+1)
	}
	dp[0][0] = 1
	for i := 1; i <= n; i++ {
		for k := 1; k <= i; k++ {
			dp[i][k] = (dp[i-1][k-1] + int(int64(k)*int64(dp[i-1][k])%int64(MOD))) % MOD
		}
	}
	bell := make([]int, n+1)
	for i := 0; i <= n; i++ {
		sum := 0
		for k := 0; k <= i; k++ {
			sum += dp[i][k]
			if sum >= MOD {
				sum -= MOD
			}
		}
		bell[i] = sum
	}
	return bell
}
