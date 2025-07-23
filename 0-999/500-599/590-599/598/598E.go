package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 31
const maxK = 51

var (
	dp  [maxN][maxN][maxK]int
	vis [maxN][maxN][maxK]bool
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, m, k int) int {
	if k == 0 || k == n*m {
		return 0
	}
	if k > n*m { // impossible but guard
		return 1<<31 - 1
	}
	if vis[n][m][k] {
		return dp[n][m][k]
	}
	vis[n][m][k] = true
	ans := 1<<31 - 1
	// horizontal cuts
	for i := 1; i <= n/2; i++ {
		maxK1 := min(k, i*m)
		for k1 := 0; k1 <= maxK1; k1++ {
			if k1 > i*m || k-k1 > (n-i)*m {
				continue
			}
			cost := m*m + solve(i, m, k1) + solve(n-i, m, k-k1)
			if cost < ans {
				ans = cost
			}
		}
	}
	// vertical cuts
	for j := 1; j <= m/2; j++ {
		maxK1 := min(k, j*n)
		for k1 := 0; k1 <= maxK1; k1++ {
			if k1 > j*n || k-k1 > n*(m-j) {
				continue
			}
			cost := n*n + solve(n, j, k1) + solve(n, m-j, k-k1)
			if cost < ans {
				ans = cost
			}
		}
	}
	dp[n][m][k] = ans
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		res := solve(n, m, k)
		fmt.Fprintln(writer, res)
	}
}
