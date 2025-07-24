package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	const INF = int(1e9)
	costs := make([][3]int, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		for j := 0; j < 3; j++ {
			costs[i][j] = INF
		}
		for j := 0; j < m; j++ {
			dist := j
			if m-j < dist {
				dist = m - j
			}
			c := s[j]
			if c >= '0' && c <= '9' {
				if dist < costs[i][0] {
					costs[i][0] = dist
				}
			} else if c >= 'a' && c <= 'z' {
				if dist < costs[i][1] {
					costs[i][1] = dist
				}
			} else if c == '#' || c == '*' || c == '&' {
				if dist < costs[i][2] {
					costs[i][2] = dist
				}
			}
		}
	}

	dp := make([][8]int, n+1)
	for i := 0; i <= n; i++ {
		for j := 0; j < 8; j++ {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0
	for i := 0; i < n; i++ {
		for mask := 0; mask < 8; mask++ {
			if dp[i][mask] == INF {
				continue
			}
			for t := 0; t < 3; t++ {
				if costs[i][t] == INF {
					continue
				}
				newMask := mask | (1 << t)
				val := dp[i][mask] + costs[i][t]
				if val < dp[i+1][newMask] {
					dp[i+1][newMask] = val
				}
			}
		}
	}
	fmt.Fprintln(writer, dp[n][7])
}
