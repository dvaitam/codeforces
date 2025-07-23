package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	// dp[l][r] represents minimal moves to remove subarray [l,r]
	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}
	for i := 0; i < n; i++ {
		dp[i][i] = 1
	}

	for length := 2; length <= n; length++ {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			best := dp[l+1][r] + 1 // remove a[l] separately
			if a[l] == a[l+1] {
				if v := 1 + dp[l+2][r]; v < best {
					best = v
				}
			}
			for k := l + 2; k <= r; k++ {
				if a[l] == a[k] {
					if v := dp[l+1][k-1] + dp[k+1][r]; v < best {
						best = v
					}
				}
			}
			dp[l][r] = best
		}
	}

	if n == 0 {
		fmt.Println(0)
	} else {
		fmt.Println(dp[0][n-1])
	}
}
