package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const negInf int64 = -1 << 60

func maxK(n int) int {
	return int((math.Sqrt(float64(8*n+1)) - 1) / 2)
}

func solve(n int, a []int64) int {
	maxLen := maxK(n)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}

	dp := make([][]int64, maxLen+1)
	suf := make([][]int64, maxLen+1)

	dp[1] = make([]int64, n)
	suf[1] = make([]int64, n+1)
	for i := 0; i <= n; i++ {
		suf[1][i] = negInf
	}
	for i := n - 1; i >= 0; i-- {
		dp[1][i] = a[i]
		if dp[1][i] > suf[1][i+1] {
			suf[1][i] = dp[1][i]
		} else {
			suf[1][i] = suf[1][i+1]
		}
	}

	for l := 2; l <= maxLen; l++ {
		dp[l] = make([]int64, n)
		suf[l] = make([]int64, n+1)
		for i := 0; i <= n; i++ {
			suf[l][i] = negInf
		}
		for i := n - 1; i >= 0; i-- {
			if i+l <= n {
				sum := prefix[i+l] - prefix[i]
				if suf[l-1][i+l] > sum {
					dp[l][i] = sum
				} else {
					dp[l][i] = negInf
				}
			} else {
				dp[l][i] = negInf
			}
			if dp[l][i] > suf[l][i+1] {
				suf[l][i] = dp[l][i]
			} else {
				suf[l][i] = suf[l][i+1]
			}
		}
	}

	for k := maxLen; k >= 1; k-- {
		if suf[k][0] > negInf {
			return k
		}
	}
	return 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		fmt.Fprintln(writer, solve(n, a))
	}
}
