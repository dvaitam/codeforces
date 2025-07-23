package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int32 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	// 1-index the string for easier handling
	bs := []byte(" " + s)

	// precompute longest common prefixes
	lcp := make([][]uint16, n+2)
	for i := range lcp {
		lcp[i] = make([]uint16, n+2)
	}
	for i := n; i >= 1; i-- {
		for j := n; j >= 1; j-- {
			if bs[i] == bs[j] {
				lcp[i][j] = lcp[i+1][j+1] + 1
			}
		}
	}

	dp := make([][]int32, n+2)
	pref := make([][]int32, n+2)
	for i := range dp {
		dp[i] = make([]int32, n+2)
		pref[i] = make([]int32, n+2)
	}

	// base case: empty prefix has one way
	dp[0][0] = 1
	pref[0][0] = 1

	for j := 1; j <= n; j++ {
		for i := 1; i <= j; i++ {
			if bs[i] == '0' {
				dp[i][j] = 0
				pref[i][j] = (pref[i-1][j] + dp[i][j]) % MOD
				continue
			}
			ans := pref[i-1][i-1]
			length := j - i + 1
			if k := i - length - 1; k >= 0 {
				val := pref[k][i-1]
				ans = (ans - val) % MOD
			}
			if k := i - length; k >= 1 {
				common := int(lcp[k][i])
				if common >= length || bs[k+common] > bs[i+common] {
					ans = (ans - dp[k][i-1]) % MOD
				}
			}
			if ans < 0 {
				ans += MOD
			}
			dp[i][j] = ans
			pref[i][j] = (pref[i-1][j] + dp[i][j]) % MOD
		}
	}

	res := pref[n][n]
	fmt.Println(res)
}
