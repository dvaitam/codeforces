package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func add(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
	}
	// dp[layer][state][value]
	var dp [2][2][201]int
	prev := 1
	// initial
	if A[0] != -1 {
		dp[prev][1][A[0]] = 1
	} else {
		for v := 1; v <= 200; v++ {
			dp[prev][1][v] = 1
		}
	}

	cur := 0
	for i := 1; i < n; i++ {
		cur, prev = prev, cur
		// clear current layer
		for s := 0; s < 2; s++ {
			for v := 1; v <= 200; v++ {
				dp[cur][s][v] = 0
			}
		}
		if A[i] != -1 {
			ai := A[i]
			for j := 1; j <= 200; j++ {
				if ai < j {
					dp[cur][0][ai] = add(dp[cur][0][ai], dp[prev][0][j])
				}
				if ai == j {
					dp[cur][0][ai] = add(dp[cur][0][ai], dp[prev][0][j])
					dp[cur][0][ai] = add(dp[cur][0][ai], dp[prev][1][j])
				}
				if ai > j {
					dp[cur][1][ai] = add(dp[cur][1][ai], dp[prev][0][j])
					dp[cur][1][ai] = add(dp[cur][1][ai], dp[prev][1][j])
				}
				// clear prev layer
				dp[prev][0][j] = 0
				dp[prev][1][j] = 0
			}
		} else {
			sum0 := 0
			for j := 1; j <= 200; j++ {
				sum0 = add(sum0, dp[prev][0][j])
			}
			sum1 := 0
			for j := 1; j <= 200; j++ {
				// dp[cur][0][j]
				dp[cur][0][j] = add(dp[cur][0][j], sum0)
				dp[cur][0][j] = add(dp[cur][0][j], dp[prev][1][j])
				// update sum0 -= dp[prev][0][j]
				sum0 = sum0 - dp[prev][0][j]
				if sum0 < 0 {
					sum0 += mod
				}
				// dp[cur][1][j]
				dp[cur][1][j] = add(dp[cur][1][j], sum1)
				// update sum1 += dp[prev][1][j] + dp[prev][0][j]
				sum1 = add(sum1, dp[prev][1][j])
				sum1 = add(sum1, dp[prev][0][j])
				// clear prev layer
				dp[prev][0][j] = 0
				dp[prev][1][j] = 0
			}
		}
	}
	// result
	res := 0
	for v := 1; v <= 200; v++ {
		res = add(res, dp[cur][0][v])
	}
	fmt.Fprintln(writer, res)
}
