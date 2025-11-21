package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const MAXN = 100

var comb [MAXN + 1][MAXN + 1]int64

func init() {
	for i := 0; i <= MAXN; i++ {
		comb[i][0] = 1
		for j := 1; j <= i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % MOD
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		s := make([]int, n+1)
		ok := true
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &s[i])
			if s[i] > 2 {
				ok = false
			}
		}
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		dp := make([][]int64, n+2)
		for i := range dp {
			dp[i] = make([]int64, n+2)
		}

		for len := 1; len <= n; len++ {
			for l := 1; l+len-1 <= n; l++ {
				r := l + len - 1
				var ways int64
				for root := l; root <= r; root++ {
					leftLen := root - l
					rightLen := r - root
					children := 0
					leftWays := int64(1)
					rightWays := int64(1)
					if leftLen > 0 {
						children++
						leftWays = dp[l][root-1]
					}
					if rightLen > 0 {
						children++
						rightWays = dp[root+1][r]
					}
					if leftLen > 0 && leftWays == 0 {
						continue
					}
					if rightLen > 0 && rightWays == 0 {
						continue
					}
					if s[root] != -1 && s[root] != children {
						continue
					}
					total := comb[len-1][leftLen]
					total = (total * leftWays) % MOD
					total = (total * rightWays) % MOD
					ways = (ways + total) % MOD
				}
				dp[l][r] = ways
			}
		}
		fmt.Fprintln(out, dp[1][n]%MOD)
	}
}
