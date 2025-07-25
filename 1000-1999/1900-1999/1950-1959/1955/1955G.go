package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]int, n)
		for i := range grid {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &grid[i][j])
			}
		}
		dp := make([][]map[int]struct{}, n)
		for i := range dp {
			dp[i] = make([]map[int]struct{}, m)
		}
		dp[0][0] = map[int]struct{}{grid[0][0]: {}}
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if i == 0 && j == 0 {
					continue
				}
				set := make(map[int]struct{})
				if i > 0 {
					for g := range dp[i-1][j] {
						set[gcd(g, grid[i][j])] = struct{}{}
					}
				}
				if j > 0 {
					for g := range dp[i][j-1] {
						set[gcd(g, grid[i][j])] = struct{}{}
					}
				}
				dp[i][j] = set
			}
		}
		best := 0
		for g := range dp[n-1][m-1] {
			if g > best {
				best = g
			}
		}
		fmt.Fprintln(out, best)
	}
}
