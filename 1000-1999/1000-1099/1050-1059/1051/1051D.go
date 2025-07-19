package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	// dp0[i][j] = number of ways state 0, dp1[i][j] for state 1
	dp0 := make([][]int, n+1)
	dp1 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp0[i] = make([]int, k+1)
		dp1[i] = make([]int, k+1)
	}
	// Base for i=1
	if k >= 2 {
		dp0[1][2] = 1
	}
	if k >= 1 {
		dp1[1][1] = 1
	}
	for i := 2; i <= n; i++ {
		for j := 1; j <= k; j++ {
			// state 0
			v0 := dp0[i-1][j]
			// from state1 with one less, two ways
			v0 += 2 * dp1[i-1][j-1]
			if j >= 2 {
				v0 += dp0[i-1][j-2]
			}
			dp0[i][j] = v0 % MOD
			// state 1
			v1 := dp1[i-1][j]
			v1 += dp1[i-1][j-1]
			v1 += 2 * dp0[i-1][j]
			dp1[i][j] = v1 % MOD
		}
	}
	res := (dp0[n][k] + dp1[n][k]) * 2 % MOD
	fmt.Println(res)
}
