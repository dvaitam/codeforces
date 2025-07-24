package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1e9 + 7

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &adj[i][j])
		}
	}

	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	for length := 1; length <= n; length++ {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			if l == r {
				dp[l][r] = 1
				continue
			}
			var res int
			for k := l; k < r; k++ {
				if adj[k][r] == 0 {
					continue
				}
				res = (res + dp[l][k]*dp[k+1][r]) % mod
			}
			dp[l][r] = res
		}
	}

	fmt.Fprintln(out, dp[0][n-1])
}
