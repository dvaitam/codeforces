package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int64 = 1e18

func minResult(arr []int) int64 {
	m := len(arr)
	dp := make([]int64, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = INF
	}
	dp[m] = 0
	for i := m - 1; i >= 0; i-- {
		prod := int64(1)
		for j := i; j < m; j++ {
			prod *= int64(arr[j])
			if prod > INF {
				prod = INF
			}
			if val := prod + dp[j+1]; val < dp[i] {
				dp[i] = val
			}
		}
	}
	return dp[0]
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
		var s string
		fmt.Fscan(reader, &s)
		digits := make([]int, n)
		for i := 0; i < n; i++ {
			digits[i] = int(s[i] - '0')
		}
		ans := INF
		for merge := 0; merge < n-1; merge++ {
			arr := make([]int, 0, n-1)
			for i := 0; i < n; i++ {
				if i == merge {
					val := digits[i]*10 + digits[i+1]
					arr = append(arr, val)
					i++
					continue
				}
				arr = append(arr, digits[i])
			}
			val := minResult(arr)
			if val < ans {
				ans = val
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
