package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf = int(1e9)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveDP(n int) int {
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}
	dp[0][0] = 0
	combos := [][2]int{
		{1, 0},
		{2, 0},
		{3, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{0, 2},
		{1, 2},
	}
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			if dp[i][j] == inf {
				continue
			}
			for _, c := range combos {
				ni := i + c[0]
				nj := j + c[1]
				if ni > n || nj > n {
					continue
				}
				if dp[i][j]+1 < dp[ni][nj] {
					dp[ni][nj] = dp[i][j] + 1
				}
			}
		}
	}
	return dp[n][n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	planks25 := (n + 1) / 2
	planks18and21 := solveDP(n)
	fmt.Fprintln(out, planks25+planks18and21)
}
